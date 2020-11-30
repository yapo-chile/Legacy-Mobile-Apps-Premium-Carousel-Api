package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExpireProductsHandlerInput(t *testing.T) {
	var h ExpireProductsHandler
	mMockInputRequest := &MockInputRequest{}
	mTargetRequest := &MockTargetRequest{}
	input := h.Input(mMockInputRequest)
	var expected *expireProductsHandlerInput
	assert.IsType(t, expected, input)
	mMockInputRequest.AssertExpectations(t)
	mTargetRequest.AssertExpectations(t)
}

type mockExpireProductsInteractor struct {
	mock.Mock
}

func (m *mockExpireProductsInteractor) ExpireProducts() error {
	args := m.Called()
	return args.Error(0)
}

func TestExpireProductsHandlerOK(t *testing.T) {
	mInteractor := &mockExpireProductsInteractor{}
	mInteractor.On("ExpireProducts").Return(nil)
	h := ExpireProductsHandler{
		Interactor: mInteractor,
	}
	input := expireProductsHandlerInput{}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: expireProductsRequestOutput{
			Response: "OK",
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestExpireProductsHandlerInteractorError(t *testing.T) {
	mInteractor := &mockExpireProductsInteractor{}
	mInteractor.On("ExpireProducts").Return(fmt.Errorf("err"))
	h := ExpireProductsHandler{
		Interactor: mInteractor,
	}
	input := expireProductsHandlerInput{}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	assert.Equal(t, http.StatusBadRequest, r.Code)
	mInteractor.AssertExpectations(t)
}
