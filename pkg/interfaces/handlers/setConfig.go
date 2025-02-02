package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Yapo/goutils"

	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/domain"
	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/usecases"
)

// SetConfigHandler implements the handler interface and responds to /ads with
// related user ads
type SetConfigHandler struct {
	Interactor usecases.SetConfigInteractor
}

// SetConfigLogger logger for SetConfig Handler
type SetConfigLogger interface{}

// setConfigHandlerInput is the handler expected input
type setConfigHandlerInput struct {
	UserProductID      int       `path:"ID"`
	Categories         string    `json:"categories"`
	Exclude            string    `json:"exclude"`
	Keywords           string    `json:"keywords"`
	Comment            string    `json:"comment"`
	Limit              int       `json:"limit"`
	PriceRange         int       `json:"price_range"`
	ExpiredAt          time.Time `json:"expiration"`
	FillGapsWithRandom bool      `json:"fill_random"`
}

// getUserRequestOutput is the handler output
type setConfigRequestOutput struct {
	response string
}

// Input returns a fresh, empty instance of setConfigHandlerInput
func (*SetConfigHandler) Input(ir InputRequest) HandlerInput {
	input := setConfigHandlerInput{}
	ir.Set(&input).FromJSONBody().FromPath()
	return &input
}

// Execute sets configuration for userProduct
func (h *SetConfigHandler) Execute(ig InputGetter) *goutils.Response {
	input, response := ig()
	if response != nil {
		return response
	}
	in := input.(*setConfigHandlerInput)
	if in.UserProductID < 1 {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: goutils.GenericError{
				ErrorMessage: fmt.Sprintf(`error with ProductID: %+v`,
					in.UserProductID),
			},
		}
	}
	config := domain.ProductParams{
		Categories:         h.getCategories(in.Categories),
		Exclude:            h.getCommaSeparedArr(in.Exclude),
		Keywords:           h.getCommaSeparedArr(in.Keywords),
		Comment:            in.Comment,
		Limit:              in.Limit,
		PriceRange:         in.PriceRange,
		FillGapsWithRandom: in.FillGapsWithRandom,
	}
	if err := h.Interactor.SetConfig(in.UserProductID,
		config, in.ExpiredAt); err != nil {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: goutils.GenericError{
				ErrorMessage: fmt.Sprintf(`%+v`, err),
			},
		}
	}
	body := setConfigRequestOutput{
		response: "OK",
	}
	return &goutils.Response{
		Code: http.StatusOK,
		Body: body,
	}
}

func (h *SetConfigHandler) getCategories(raw string) (categories []int) {
	if raw == "" {
		return []int{}
	}
	categoriesArr := strings.Split(raw, ",")
	for _, c := range categoriesArr {
		cat, _ := strconv.Atoi(c)
		categories = append(categories, cat)
	}
	return categories
}

func (h *SetConfigHandler) getCommaSeparedArr(raw string) []string {
	if raw == "" {
		return []string{}
	}
	return strings.Split(raw, ",")
}
