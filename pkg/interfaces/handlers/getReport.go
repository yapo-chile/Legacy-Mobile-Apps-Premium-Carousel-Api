package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Yapo/goutils"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

// GetReportHandler implements the handler interface and responds to /report with
// sales report data
type GetReportHandler struct {
	Interactor usecases.GetReportInteractor
}

// GetReportLogger logger for GetReport Handler
type GetReportLogger interface{}

// getReportHandlerInput is the handler expected input
type getReportHandlerInput struct {
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
}

// getReportRequestOutput is the handler output
type getReportRequestOutput struct {
	Products []productsOutput `json:"assigns"`
}

// Input returns a fresh, empty instance of getReportHandlerInput
func (*GetReportHandler) Input(ir InputRequest) HandlerInput {
	input := getReportHandlerInput{}
	ir.Set(&input).FromQuery()
	return &input
}

// Execute gets sales report for controlpanel
func (h *GetReportHandler) Execute(ig InputGetter) *goutils.Response {
	input, response := ig()
	if response != nil {
		return response
	}
	in := input.(*getReportHandlerInput)
	startDate, err := time.Parse(time.RFC3339, in.StartDate)
	if err != nil {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: goutils.GenericError{
				ErrorMessage: fmt.Sprintf(`bad start_date format: %+v`, err),
			},
		}
	}
	endDate, err := time.Parse(time.RFC3339, in.EndDate)
	if err != nil {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: goutils.GenericError{
				ErrorMessage: fmt.Sprintf(`bad end_date format: %+v`, err),
			},
		}
	}
	if startDate.After(endDate) {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: goutils.GenericError{
				ErrorMessage: fmt.Sprintf(`invalid date interval`),
			},
		}
	}
	products, err := h.Interactor.GetReport(startDate, endDate)
	if err != nil {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: goutils.GenericError{
				ErrorMessage: fmt.Sprintf(`%+v`, err),
			},
		}
	}
	productsOut := []productsOutput{}
	for _, v := range products {
		p := productsOutput{
			ID:             v.ID,
			Email:          v.Email,
			UserID:         strconv.Itoa(v.UserID),
			Status:         string(v.Status),
			Type:           string(v.Type),
			PurchaseNumber: v.Purchase.Number,
			PurchasePrice:  v.Purchase.Price,
			PurchaseStatus: string(v.Purchase.Status),
			ExpiredAt:      v.ExpiredAt,
			CreatedAt:      v.CreatedAt,
			Comment:        v.Config.Comment,
			Keywords:       strings.Join(v.Config.Keywords, ","),
			PriceRange:     v.Config.PriceRange,
			Categories: strings.Trim(strings.Join(
				strings.Fields(fmt.Sprint(v.Config.Categories)), ","), "[]"),
			Limit:              v.Config.Limit,
			FillGapsWithRandom: v.Config.FillGapsWithRandom,
		}
		productsOut = append(productsOut, p)
	}
	body := getReportRequestOutput{
		Products: productsOut,
	}
	return &goutils.Response{
		Code: http.StatusOK,
		Body: body,
	}
}
