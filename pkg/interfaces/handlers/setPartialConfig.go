package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Yapo/goutils"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

// SetPartialConfigHandler implements the handler interface and responds to /ads with
// related user ads
type SetPartialConfigHandler struct {
	Interactor usecases.SetPartialConfigInteractor
}

// SetPartialConfigLogger logger for SetPartialConfig Handler
type SetPartialConfigLogger interface{}

// setPartialConfigHandlerInput is the handler expected input
type setPartialConfigHandlerInput struct {
	UserProductID int    `path:"ID"`
	Body          []byte `raw:"body"`
}

// getUserRequestOutput is the handler output
type setPartialConfigRequestOutput struct {
	response string
}

// Input returns a fresh, empty instance of setPartialConfigHandlerInput
func (*SetPartialConfigHandler) Input(ir InputRequest) HandlerInput {
	input := setPartialConfigHandlerInput{}
	ir.Set(&input).FromRawBody().FromPath()
	return &input
}

// Execute sets partial configuration for supported params
func (h *SetPartialConfigHandler) Execute(ig InputGetter) *goutils.Response {
	input, response := ig()
	if response != nil {
		return response
	}
	in := input.(*setPartialConfigHandlerInput)
	if in.UserProductID < 1 {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: goutils.GenericError{
				ErrorMessage: fmt.Sprintf(`Wrong ProductID: %d`, in.UserProductID),
			},
		}
	}
	configMap := make(map[string]interface{})
	err := json.Unmarshal(in.Body, &configMap)
	if err != nil {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: goutils.GenericError{
				ErrorMessage: fmt.Sprintf(`error decoding input: %+v`, err),
			},
		}
	}
	err = h.Interactor.SetPartialConfig(in.UserProductID, configMap)
	if err != nil {
		return &goutils.Response{
			Code: http.StatusBadRequest,
			Body: goutils.GenericError{
				ErrorMessage: fmt.Sprintf(`%+v`, err),
			},
		}
	}
	body := setPartialConfigRequestOutput{
		response: "OK",
	}
	return &goutils.Response{
		Code: http.StatusOK,
		Body: body,
	}
}
