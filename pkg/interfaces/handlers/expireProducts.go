package handlers

import (
	"fmt"
	"net/http"

	"github.com/Yapo/goutils"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

// ExpireProductsHandler implements expire product handler to respond to /expire-product
type ExpireProductsHandler struct {
	Interactor usecases.ExpireProductsInteractor
}

// ExpireProductsLogger logger for ExpireProducts Handler
type ExpireProductsLogger interface{}

// expireProductsHandlerInput is the handler expected input
type expireProductsHandlerInput struct{}

// getUserRequestOutput is the handler output
type expireProductsRequestOutput struct {
	Response string `json:"status"`
}

// Input returns a fresh, empty instance of expireProductsHandlerInput
func (*ExpireProductsHandler) Input(ir InputRequest) HandlerInput {
	input := expireProductsHandlerInput{}
	return &input
}

// Execute edits a new user product using controlpanel
func (h *ExpireProductsHandler) Execute(ig InputGetter) *goutils.Response {
	if err := h.Interactor.ExpireProducts(); err != nil {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: goutils.GenericError{
				ErrorMessage: fmt.Sprintf(`%+v`, err),
			},
		}
	}
	return &goutils.Response{
		Code: http.StatusOK,
		Body: expireProductsRequestOutput{
			Response: "OK",
		},
	}
}
