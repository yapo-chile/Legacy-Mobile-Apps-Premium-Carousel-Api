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
	CustomQuery        string    `json:"keywords"`
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
	fmt.Printf(`set config input: %+v`, input)
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
	if in.ExpiredAt.Before(time.Now()) {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: goutils.GenericError{
				ErrorMessage: fmt.Sprintf(`bad expiration date: %+v`,
					in.ExpiredAt),
			},
		}
	}
	config := usecases.CpConfig{
		Categories:         h.getCategories(in.Categories),
		Exclude:            h.getExclude(in.Exclude),
		CustomQuery:        in.CustomQuery,
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

func (h *SetConfigHandler) getExclude(raw string) []string {
	if raw == "" {
		return []string{}
	}
	return strings.Split(raw, ",")
}
