package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

func TestSetConfigHandlerInput(t *testing.T) {
	var h SetConfigHandler
	mMockInputRequest := &MockInputRequest{}
	mTargetRequest := &MockTargetRequest{}
	mMockInputRequest.On("Set",
		mock.AnythingOfType("*handlers.setConfigHandlerInput")).Return(mTargetRequest)
	mTargetRequest.On("FromJSONBody").Return(mTargetRequest)
	mTargetRequest.On("FromPath").Return(mTargetRequest)
	input := h.Input(mMockInputRequest)
	var expected *setConfigHandlerInput
	assert.IsType(t, expected, input)
	mMockInputRequest.AssertExpectations(t)
	mTargetRequest.AssertExpectations(t)
}

type mockSetConfigInteractor struct {
	mock.Mock
}

func (m *mockSetConfigInteractor) SetConfig(userProductID int,
	config domain.ProductParams, expiredAt time.Time) error {
	args := m.Called(userProductID, config, expiredAt)
	return args.Error(0)
}

func TestSetConfigHandlerErrorBadInput(t *testing.T) {
	mInteractor := &mockSetConfigInteractor{}
	h := SetConfigHandler{
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

func TestSetConfigHandlerOK(t *testing.T) {
	mInteractor := &mockSetConfigInteractor{}
	mInteractor.On("SetConfig",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.ProductParams"),
		mock.AnythingOfType("time.Time"),
	).Return(nil)
	h := SetConfigHandler{
		Interactor: mInteractor,
	}
	input := setConfigHandlerInput{
		UserProductID: 123,
		Categories:    "2000,1000,3000",
		ExpiredAt:     time.Now().Add(time.Hour * 24 * 365),
		Exclude:       "12345",
	}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: setConfigRequestOutput{
			response: "OK",
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestSetConfigHandlerError(t *testing.T) {
	err := fmt.Errorf("err")
	mInteractor := &mockSetConfigInteractor{}
	mInteractor.On("SetConfig",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.ProductParams"),
		mock.AnythingOfType("time.Time"),
	).Return(err)
	h := SetConfigHandler{
		Interactor: mInteractor,
	}
	input := setConfigHandlerInput{
		UserProductID: 123,
		ExpiredAt:     time.Now().Add(time.Hour * 24 * 365),
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

func TestSetConfigHandlerBadUserProductID(t *testing.T) {
	mInteractor := &mockSetConfigInteractor{}
	h := SetConfigHandler{
		Interactor: mInteractor,
	}
	input := setConfigHandlerInput{
		UserProductID: 0,
	}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: goutils.GenericError{
			ErrorMessage: fmt.Sprintf(`error with ProductID: %+v`,
				input.UserProductID),
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}
