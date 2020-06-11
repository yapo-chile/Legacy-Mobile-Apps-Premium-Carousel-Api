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

func TestGetReportHandlerInput(t *testing.T) {
	var h GetReportHandler
	mMockInputRequest := &MockInputRequest{}
	mTargetRequest := &MockTargetRequest{}
	mMockInputRequest.On("Set",
		mock.AnythingOfType("*handlers.getReportHandlerInput")).Return(mTargetRequest)
	mTargetRequest.On("FromQuery").Return(mTargetRequest)
	input := h.Input(mMockInputRequest)
	var expected *getReportHandlerInput
	assert.IsType(t, expected, input)
	mMockInputRequest.AssertExpectations(t)
	mTargetRequest.AssertExpectations(t)
}

type mockGetReportInteractor struct {
	mock.Mock
}

func (m *mockGetReportInteractor) GetReport(start,
	end time.Time) ([]domain.Product, error) {
	args := m.Called(start, end)
	return args.Get(0).([]domain.Product), args.Error(1)
}

func TestGetReportHandlerErrorBadInput(t *testing.T) {
	mInteractor := &mockGetReportInteractor{}
	h := GetReportHandler{
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

func TestGetReportHandlerOK(t *testing.T) {
	mInteractor := &mockGetReportInteractor{}
	testTime := time.Now()
	mInteractor.On("GetReport", mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("time.Time")).
		Return([]domain.Product{{ID: 123, Purchase: domain.Purchase{ID: 1,
			Type: domain.AdminPurchase}}}, nil)
	h := GetReportHandler{
		Interactor: mInteractor,
	}
	input := getReportHandlerInput{
		StartDate: testTime.Format(time.RFC3339),
		EndDate:   testTime.Add(time.Hour).Format(time.RFC3339),
	}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: getReportRequestOutput{
			Products: []productsOutput{{ID: 123, UserID: "0", PurchaseID: 1,
				PurchaseType: "ADMIN"}},
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestGetReportHandlerError(t *testing.T) {
	mInteractor := &mockGetReportInteractor{}
	testTime := time.Now()
	err := fmt.Errorf("err")
	mInteractor.On("GetReport", mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("time.Time")).
		Return([]domain.Product{{ID: 123}}, err)
	h := GetReportHandler{
		Interactor: mInteractor,
	}
	input := getReportHandlerInput{
		StartDate: testTime.Format(time.RFC3339),
		EndDate:   testTime.Add(time.Hour).Format(time.RFC3339),
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

func TestGetReportHandlerErrorNotValidInterval(t *testing.T) {
	mInteractor := &mockGetReportInteractor{}
	testTime := time.Now()
	h := GetReportHandler{
		Interactor: mInteractor,
	}
	input := getReportHandlerInput{
		StartDate: testTime.Add(time.Hour).Format(time.RFC3339),
		EndDate:   testTime.Format(time.RFC3339),
	}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: goutils.GenericError{
			ErrorMessage: fmt.Sprintf(`invalid date interval`),
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
}

func TestGetReportHandlerErrorBadStartDate(t *testing.T) {
	mInteractor := &mockGetReportInteractor{}
	testTime := time.Now()
	h := GetReportHandler{
		Interactor: mInteractor,
	}
	input := getReportHandlerInput{
		StartDate: "asdf",
		EndDate:   testTime.Format(time.RFC3339),
	}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusBadRequest,
	}
	assert.Equal(t, expected.Code, r.Code)
	mInteractor.AssertExpectations(t)
}

func TestGetReportHandlerErrorBadEndDate(t *testing.T) {
	mInteractor := &mockGetReportInteractor{}
	testTime := time.Now()
	h := GetReportHandler{
		Interactor: mInteractor,
	}
	input := getReportHandlerInput{
		StartDate: testTime.Format(time.RFC3339),
		EndDate:   "asdf",
	}
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusBadRequest,
	}
	assert.Equal(t, expected.Code, r.Code)
	mInteractor.AssertExpectations(t)
}
