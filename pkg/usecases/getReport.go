package usecases

import (
	"fmt"
	"time"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

// GetReportInteractor wraps GetReport operations
type GetReportInteractor interface {
	GetReport(startDate,
		endDate time.Time) (products []domain.Product, err error)
}

// getReportInteractor defines the interactor for GetReport usecase
type getReportInteractor struct {
	productRepo ProductRepository
	logger      GetReportLogger
}

// GetReportLogger logs GetReport events
type GetReportLogger interface {
	LogErrorGettingReport(err error)
}

// MakeGetReportInteractor creates a new instance of GetReportInteractor
func MakeGetReportInteractor(productRepo ProductRepository,
	logger GetReportLogger) GetReportInteractor {
	return &getReportInteractor{productRepo: productRepo, logger: logger}
}

// GetReport gets sales report using start & end date
func (interactor *getReportInteractor) GetReport(startDate,
	endDate time.Time) (products []domain.Product, err error) {
	products, err = interactor.productRepo.GetReport(startDate, endDate)
	if err != nil {
		interactor.logger.LogErrorGettingReport(err)
		return []domain.Product{}, fmt.Errorf("error loading report: %+v", err)
	}
	return
}
