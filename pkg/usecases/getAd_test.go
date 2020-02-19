package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

func TestGetAdOk(t *testing.T) {
	mAdRepo := &mockAdRepo{}
	interactor := MakeGetAdInteractor(mAdRepo)
	tAd := domain.Ad{ID: "1", Subject: "Mi auto", UserID: "123"}
	mAdRepo.On("GetAd", mock.AnythingOfType("string")).Return(tAd, nil)
	ads, err := interactor.GetAd("1")
	expected := tAd
	assert.NoError(t, err)
	assert.Equal(t, expected, ads)
	mAdRepo.AssertExpectations(t)
}
