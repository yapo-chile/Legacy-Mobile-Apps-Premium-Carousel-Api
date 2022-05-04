package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/domain"
)

func TestAddUserProductHandlerInput(t *testing.T) {
	var h AddUserProductHandler
	mMockInputRequest := &MockInputRequest{}
	mTargetRequest := &MockTargetRequest{}
	mMockInputRequest.On("Set",
		mock.AnythingOfType("*handlers.addUserProductHandlerInput")).Return(mTargetRequest)
	mTargetRequest.On("FromJSONBody").Return(mTargetRequest)
	input := h.Input(mMockInputRequest)
	var expected *addUserProductHandlerInput
	assert.IsType(t, expected, input)
	mMockInputRequest.AssertExpectations(t)
	mTargetRequest.AssertExpectations(t)
}

type mockAddUserProductInteractor struct {
	mock.Mock
}

func (m *mockAddUserProductInteractor) AddUserProduct(userID int,
	email string, purchaseNumber, purchasePrice int,
	purchaseType domain.PurchaseType, productType domain.ProductType,
	expiredAt time.Time, config domain.ProductParams) error {
	args := m.Called(userID, email, purchaseNumber, purchasePrice,
		purchaseType, productType, expiredAt, config)
	return args.Error(0)
}

func TestAddUserProductHandlerErrorBadInput(t *testing.T) {
	mInteractor := &mockAddUserProductInteractor{}
	h := AddUserProductHandler{
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

func TestAddUserProductHandlerOK(t *testing.T) {
	mInteractor := &mockAddUserProductInteractor{}
	mInteractor.On("AddUserProduct",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.PurchaseType"),
		domain.PremiumCarousel,
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("domain.ProductParams"),
	).Return(nil)
	h := AddUserProductHandler{
		Interactor: mInteractor,
	}
	input := addUserProductHandlerInput{
		UserID:     123,
		Email:      "test@test.cl",
		Categories: "2000,1000,3000",
		ExpiredAt:  time.Now().Add(time.Hour * 24 * 365),
	}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: addUserProductRequestOutput{
			response: "OK",
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestAddUserProductHandlerBadPurchaseType(t *testing.T) {
	mInteractor := &mockAddUserProductInteractor{}
	h := AddUserProductHandler{
		Interactor: mInteractor,
	}
	input := addUserProductHandlerInput{
		UserID:       123,
		PurchaseType: "ASDKAD",
		Email:        "test@test.cl",
		Categories:   "2000,1000,3000",
		ExpiredAt:    time.Now().Add(time.Hour * 24 * 365),
	}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)

	assert.Equal(t, http.StatusBadRequest, r.Code)
	mInteractor.AssertExpectations(t)
}

func TestAddUserProductHandlerError(t *testing.T) {
	err := fmt.Errorf("err")
	mInteractor := &mockAddUserProductInteractor{}
	mInteractor.On("AddUserProduct",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.PurchaseType"),
		domain.PremiumCarousel,
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("domain.ProductParams"),
	).Return(err)
	h := AddUserProductHandler{
		Interactor: mInteractor,
	}
	input := addUserProductHandlerInput{
		UserID:    123,
		Email:     "test@test.cl",
		ExpiredAt: time.Now().Add(time.Hour * 24 * 365),
		Exclude:   "1234",
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

func TestAddUserProductHandlerBadExpiredAtTime(t *testing.T) {
	mInteractor := &mockAddUserProductInteractor{}
	h := AddUserProductHandler{
		Interactor: mInteractor,
	}
	input := addUserProductHandlerInput{
		UserID:    123,
		Email:     "test@test.cl",
		ExpiredAt: time.Now().Add(-1 * time.Hour * 24 * 365),
	}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: goutils.GenericError{
			ErrorMessage: fmt.Sprintf(`bad expiration date: %+v`,
				input.ExpiredAt),
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}
