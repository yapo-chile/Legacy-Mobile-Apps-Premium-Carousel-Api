package handlers

import (
	"fmt"
	"net/http"

	"github.com/Yapo/goutils"

	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/usecases"
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
	UserProductID int                    `path:"ID"`
	Body          map[string]interface{} `body:"body"`
}

// getUserRequestOutput is the handler output
type setPartialConfigRequestOutput struct {
	response string
}

// Input returns a fresh, empty instance of setPartialConfigHandlerInput
func (*SetPartialConfigHandler) Input(ir InputRequest) HandlerInput {
	input := setPartialConfigHandlerInput{}
	ir.Set(&input).FromPath()
	ir.Set(&input.Body).FromJSONBody()
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
	err := h.Interactor.SetPartialConfig(in.UserProductID, in.Body)
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
