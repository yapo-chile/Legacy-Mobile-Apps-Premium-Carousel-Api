package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

func TestGetUserProductsHandlerInput(t *testing.T) {
	var h GetUserProductsHandler
	mMockInputRequest := &MockInputRequest{}
	mTargetRequest := &MockTargetRequest{}
	mMockInputRequest.On("Set",
		mock.AnythingOfType("*handlers.getUserProductsHandlerInput")).Return(mTargetRequest)
	mTargetRequest.On("FromQuery").Return(mTargetRequest)
	input := h.Input(mMockInputRequest)
	var expected *getUserProductsHandlerInput
	assert.IsType(t, expected, input)
	mMockInputRequest.AssertExpectations(t)
	mTargetRequest.AssertExpectations(t)
}

type mockGetUserProductsInteractor struct {
	mock.Mock
}

func (m *mockGetUserProductsInteractor) GetUserProducts(email string,
	page int) ([]usecases.Product, int, int, error) {
	args := m.Called(email, page)
	return args.Get(0).([]usecases.Product), args.Int(1), args.Int(2), args.Error(3)
}

func TestGetUserProductsHandlerErrorBadInput(t *testing.T) {
	mInteractor := &mockGetUserProductsInteractor{}
	h := GetUserProductsHandler{
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

func TestGetUserProductsHandlerOK(t *testing.T) {
	mInteractor := &mockGetUserProductsInteractor{}
	mInteractor.On("GetUserProducts", mock.AnythingOfType("string"),
		mock.AnythingOfType("int")).
		Return([]usecases.Product{{ID: 123}}, 1, 1, nil)
	h := GetUserProductsHandler{
		Interactor: mInteractor,
	}
	input := getUserProductsHandlerInput{
		Email: "test@test.cl",
		Page:  1,
	}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: getUserProductsRequestOutput{
			Products: []productsOutput{{ID: 123}},
			Metadata: metadata{CurrentPage: 1, TotalPages: 1},
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestGetUserProductsHandlerError(t *testing.T) {
	mInteractor := &mockGetUserProductsInteractor{}
	err := fmt.Errorf("err")
	mInteractor.On("GetUserProducts", mock.AnythingOfType("string"),
		mock.AnythingOfType("int")).
		Return([]usecases.Product{}, 0, 0, err)
	h := GetUserProductsHandler{
		Interactor: mInteractor,
	}
	input := getUserProductsHandlerInput{
		Email: "test@test.cl",
		Page:  1,
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
