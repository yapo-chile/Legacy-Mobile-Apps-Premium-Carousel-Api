package usecases

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/domain"
)

type mockGetReportLogger struct {
	mock.Mock
}

func (m *mockGetReportLogger) LogErrorGettingReport(err error) {
	m.Called(err)
}

func TestGetReportOk(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mLogger := &mockGetReportLogger{}
	testTime := time.Now()
	interactor := MakeGetReportInteractor(mProductRepo, mLogger)
	products := []domain.Product{}
	mProductRepo.On("GetReport",
		testTime,
		testTime,
	).Return(products, nil)
	res, err := interactor.GetReport(testTime, testTime)
	assert.NoError(t, err)
	assert.Equal(t, products, res)
	mProductRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetReportError(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mLogger := &mockGetReportLogger{}
	testTime := time.Now()
	interactor := MakeGetReportInteractor(mProductRepo, mLogger)
	mLogger.On("LogErrorGettingReport",
		mock.Anything, mock.Anything)
	mProductRepo.On("GetReport",
		testTime,
		testTime,
	).Return([]domain.Product{}, fmt.Errorf("err"))
	_, err := interactor.GetReport(testTime, testTime)
	assert.Error(t, err)
	mProductRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}
