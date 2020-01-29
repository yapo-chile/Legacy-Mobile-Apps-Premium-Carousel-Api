package handlers

import (
	"net/http"

	"github.com/Yapo/goutils"
)

// GetHealthcheckInteractor allows the interaction between the handler and the use case.
// The method GetHealthcheck is implemented in the use case and allows the MS to query
// it's own state
type GetHealthcheckInteractor interface {
	GetHealthcheck() (string, error)
}

// GetHealthcheckHandler implements the handler interface and responds to
// /getHealthcheck requests using an interactor.
type GetHealthcheckHandler struct {
	GetHealthcheckInteractor GetHealthcheckInteractor
}

type getHealthcheckHandlerInput struct{}

type getHealthcheckRequestOutput struct {
	Status string `json:"ClientStatus"`
}

// Input returns a fresh, empty instance of getHealthcheckHandler
func (h *GetHealthcheckHandler) Input(ir InputRequest) HandlerInput {
	return &getHealthcheckHandlerInput{}
}

type getHealthcheckRequestError goutils.GenericError

// Execute queries itself to get the service health status.
func (h *GetHealthcheckHandler) Execute(ig InputGetter) *goutils.Response {
	resp, err := h.GetHealthcheckInteractor.GetHealthcheck()
	if err != nil {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: getHealthcheckRequestError{
				err.Error(),
			},
		}
	}

	return &goutils.Response{
		Code: http.StatusOK,
		Body: getHealthcheckRequestOutput{
			Status: resp,
		},
	}
}
