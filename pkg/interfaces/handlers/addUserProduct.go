package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Yapo/goutils"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
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
	PurchaseOrder      int       `json:"purchase_order"`
	PurchasePrice      int       `json:"purchase_price"`
	PurchaseType       string    `json:"purchase_type"`
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
	config := domain.ProductParams{
		Categories:         h.getCategories(in.Categories),
		Exclude:            h.getCommaSeparedArr(in.Exclude),
		Keywords:           h.getCommaSeparedArr(in.Keywords),
		Limit:              in.Limit,
		PriceRange:         in.PriceRange,
		FillGapsWithRandom: in.FillGapsWithRandom,
		Comment:            in.Comment,
	}

	purchaseType, err := h.getPurchaseType(in.PurchaseType)
	if err != nil {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: goutils.GenericError{
				ErrorMessage: fmt.Sprintf(`%+v`, err),
			},
		}
	}
	err = h.Interactor.AddUserProduct(strconv.Itoa(in.UserID), in.Email,
		in.PurchaseOrder, in.PurchasePrice, purchaseType,
		domain.PremiumCarousel, in.ExpiredAt, config)
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

func (h *AddUserProductHandler) getCategories(raw string) (categories []int) {
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

func (h *AddUserProductHandler) getCommaSeparedArr(raw string) []string {
	if raw == "" {
		return []string{}
	}
	return strings.Split(raw, ",")
}

func (h *AddUserProductHandler) getPurchaseType(raw string) (domain.PurchaseType, error) {
	switch raw {
	case "": // retrocompatibility site version 23.03.00
		fallthrough
	case "admin":
		return domain.AdminPurchase, nil
	default:
		return "", fmt.Errorf("PurchaseType %s not supported", raw)
	}
}
