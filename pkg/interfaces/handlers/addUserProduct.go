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

// AddUserProductHandler implements the handler interface and responds to /ads with
// related user ads
type AddUserProductHandler struct {
	Interactor usecases.AddUserProductInteractor
}

// AddUserProductLogger logger for AddUserProduct Handler
type AddUserProductLogger interface{}

// addUserProductHandlerInput is the handler expected input
type addUserProductHandlerInput struct {
	UserID             int       `json:"user_id"`
	Email              string    `json:"email"`
	Categories         string    `json:"categories"`
	Exclude            string    `json:"exclude"`
	CustomQuery        string    `json:"keywords"`
	Comment            string    `json:"comment"`
	Limit              int       `json:"limit"`
	PriceRange         int       `json:"price_range"`
	ExpiredAt          time.Time `json:"expiration"`
	FillGapsWithRandom bool      `json:"fill_random"`
}

// getUserRequestOutput is the handler output
type addUserProductRequestOutput struct {
	response string
}

// Input returns a fresh, empty instance of addUserProductHandlerInput
func (*AddUserProductHandler) Input(ir InputRequest) HandlerInput {
	input := addUserProductHandlerInput{}
	ir.Set(&input).FromJSONBody()
	return &input
}

// Execute adds a new user product using controlpanel
func (h *AddUserProductHandler) Execute(ig InputGetter) *goutils.Response {
	input, response := ig()
	if response != nil {
		return response
	}
	in := input.(*addUserProductHandlerInput)
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
		Exclude:            strings.Split(in.Exclude, ","),
		CustomQuery:        in.CustomQuery,
		Limit:              in.Limit,
		PriceRange:         in.PriceRange,
		FillGapsWithRandom: in.FillGapsWithRandom,
	}

	err := h.Interactor.AddUserProduct(strconv.Itoa(in.UserID), in.Email, in.Comment,
		usecases.PremiumCarousel, in.ExpiredAt, config)
	if err != nil {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: goutils.GenericError{
				ErrorMessage: fmt.Sprintf(`%+v`, err),
			},
		}
	}
	body := addUserProductRequestOutput{
		response: "OK",
	}

	return &goutils.Response{
		Code: http.StatusOK,
		Body: body,
	}
}

func (h *AddUserProductHandler) getCategories(raw string) []int {
	categories := []int{}
	categoriesArr := strings.Split(raw, ",")
	for _, c := range categoriesArr {
		cat, _ := strconv.Atoi(c)
		categories = append(categories, cat)
	}
	return categories
}
