package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetPartialConfigHandlerInput(t *testing.T) {
	var h SetPartialConfigHandler
	mMockInputRequest := &MockInputRequest{}
	mTargetRequest := &MockTargetRequest{}
	mMockInputRequest.On("Set",
		mock.AnythingOfType("*handlers.setPartialConfigHandlerInput")).Return(mTargetRequest)
	mTargetRequest.On("FromRawBody").Return(mTargetRequest)
	mTargetRequest.On("FromPath").Return(mTargetRequest)
	input := h.Input(mMockInputRequest)
	var expected *setPartialConfigHandlerInput
	assert.IsType(t, expected, input)
	mMockInputRequest.AssertExpectations(t)
	mTargetRequest.AssertExpectations(t)
}

type mockSetPartialConfigInteractor struct {
	mock.Mock
}

func (m *mockSetPartialConfigInteractor) SetPartialConfig(userProductID int,
	configMap map[string]interface{}) error {
	args := m.Called(userProductID, configMap)
	return args.Error(0)
}

func TestSetPartialConfigHandlerErrorBadInput(t *testing.T) {
	mInteractor := &mockSetPartialConfigInteractor{}
	h := SetPartialConfigHandler{
		Interactor: mInteractor,
	}
	var input getUserAdsHandlerInput
	getter := MakeMockInputGetter(&input, &goutils.Response{
		Code: http.StatusNoContent,
	})
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusNoContent,
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestSetPartialConfigHandlerOK(t *testing.T) {
	mInteractor := &mockSetPartialConfigInteractor{}
	mInteractor.On("SetPartialConfig",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("map[string]interface {}"),
	).Return(nil)
	h := SetPartialConfigHandler{
		Interactor: mInteractor,
	}
	input := setPartialConfigHandlerInput{
		UserProductID: 123,
		Body:          []byte(`{"status":"ACTIVE"}`),
	}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: setPartialConfigRequestOutput{
			response: "OK",
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestSetPartialConfigHandlerError(t *testing.T) {
	err := fmt.Errorf("err")
	mInteractor := &mockSetPartialConfigInteractor{}
	mInteractor.On("SetPartialConfig",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("map[string]interface {}"),
	).Return(err)
	h := SetPartialConfigHandler{
		Interactor: mInteractor,
	}
	input := setPartialConfigHandlerInput{
		UserProductID: 123,
		Body:          []byte(`{"status":"ACTIVE"}`),
	}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: goutils.GenericError{
			ErrorMessage: fmt.Sprintf(`%+v`, err),
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestSetPartialConfigHandlerBadUserProductID(t *testing.T) {
	mInteractor := &mockSetPartialConfigInteractor{}
	h := SetPartialConfigHandler{
		Interactor: mInteractor,
	}
	input := setPartialConfigHandlerInput{
		UserProductID: 0,
	}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: goutils.GenericError{
			ErrorMessage: fmt.Sprintf(`Wrong ProductID: %d`,
				input.UserProductID),
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestSetPartialConfigHandlerErrorDecodingInput(t *testing.T) {
	mInteractor := &mockSetPartialConfigInteractor{}
	h := SetPartialConfigHandler{
		Interactor: mInteractor,
	}
	input := setPartialConfigHandlerInput{
		UserProductID: 1,
		Body:          []byte("{{{{{{"),
	}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: goutils.GenericError{
			ErrorMessage: fmt.Sprintf(`error decoding input:` +
				` invalid character '{' looking for beginning of object key string`),
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}
