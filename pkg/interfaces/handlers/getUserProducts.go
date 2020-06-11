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

// GetUserProductsHandler implements the handler interface and responds to /ads with
// related user ads
type GetUserProductsHandler struct {
	Interactor usecases.GetUserProductsInteractor
}

// GetUserProductsLogger logger for GetUserProducts Handler
type GetUserProductsLogger interface{}

// getUserProductsHandlerInput is the handler expected input
type getUserProductsHandlerInput struct {
	Email string `query:"email"`
	Page  int    `query:"page"`
}

// getUserRequestOutput is the handler output
type getUserProductsRequestOutput struct {
	Products []productsOutput `json:"assigns"`
	Metadata metadata         `json:"metadata"`
}

type productsOutput struct {
	ID                 int       `json:"id"`
	UserID             string    `json:"user_id"`
	Email              string    `json:"email"`
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	PurchaseID         int       `json:"purchase_id"`
	PurchaseNumber     int       `json:"purchase_number"`
	PurchasePrice      int       `json:"purchase_price"`
	PurchaseStatus     string    `json:"purchase_status"`
	PurchaseType       string    `json:"purchase_type"`
	ExpiredAt          time.Time `json:"expiration"`
	CreatedAt          time.Time `json:"creation"`
	Comment            string    `json:"comment"`
	Keywords           string    `json:"keywords"`
	Categories         string    `json:"categories"`
	PriceRange         int       `json:"price_range"`
	Limit              int       `json:"limit"`
	FillGapsWithRandom bool      `json:"fill_random"`
}

type metadata struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
}

// Input returns a fresh, empty instance of getUserProductsHandlerInput
func (*GetUserProductsHandler) Input(ir InputRequest) HandlerInput {
	input := getUserProductsHandlerInput{}
	ir.Set(&input).FromQuery()
	return &input
}

// Execute get a list of user products using pagination for controlpanel
func (h *GetUserProductsHandler) Execute(ig InputGetter) *goutils.Response {
	input, response := ig()
	if response != nil {
		return response
	}
	in := input.(*getUserProductsHandlerInput)
	products, currentPage, totalPages, err := h.Interactor.GetUserProducts(in.Email, in.Page)
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
			PurchaseID:     v.Purchase.ID,
			PurchaseNumber: v.Purchase.Number,
			PurchasePrice:  v.Purchase.Price,
			PurchaseStatus: string(v.Purchase.Status),
			PurchaseType:   string(v.Purchase.Type),
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
	body := getUserProductsRequestOutput{
		Products: productsOut,
		Metadata: metadata{
			CurrentPage: currentPage,
			TotalPages:  totalPages,
		},
	}
	return &goutils.Response{
		Code: http.StatusOK,
		Body: body,
	}
}
