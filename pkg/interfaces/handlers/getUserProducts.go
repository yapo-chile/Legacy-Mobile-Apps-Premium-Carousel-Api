package handlers

import (
	"fmt"
	"net/http"
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
	Products []productsOutput
	Metadata metadata
}

type productsOutput struct {
	ID        int
	UserID    string
	Email     string
	Type      string
	Status    string
	ExpiredAt time.Time
	CreatedAt time.Time
	Comment   string
	Config    usecases.CpConfig
}

type metadata struct {
	CurrentPage int
	TotalPages  int
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
			Body: fmt.Sprintf(`{"error": "%+v"}`, err),
		}
	}
	productsOut := []productsOutput{}
	for _, v := range products {
		p := productsOutput{
			ID:        v.ID,
			Email:     v.Email,
			UserID:    v.UserID,
			Status:    string(v.Status),
			Type:      string(v.Type),
			ExpiredAt: v.ExpiredAt,
			CreatedAt: v.CreatedAt,
			Comment:   v.Comment,
			Config:    v.Config,
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
