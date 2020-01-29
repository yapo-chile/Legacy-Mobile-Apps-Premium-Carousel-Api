package repository

import (
	"encoding/json"
	"fmt"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

// HTTPAdRepo implements the repository interface and gets ads from search-ms
type HTTPAdRepo struct {
	Handler HTTPHandler
	Path    string
}

// NewHTTPAdRepo returns a fresh instance of NewHTTPAdRepo
func NewHTTPAdRepo(handler HTTPHandler, path string) *HTTPAdRepo {
	return &HTTPAdRepo{
		Handler: handler,
		Path:    path,
	}
}

// Get gets ads from search microservice
func (repo *HTTPAdRepo) Get(ads usecases.SearchInput) (usecases.SearchResponse, error) {
	request := repo.Handler.NewRequest().
		SetMethod("POST").
		SetPath(repo.Path).
		SetBody(ads)
	adJSON, err := repo.Handler.Send(request)
	if err == nil && adJSON != nil {
		ad := fmt.Sprintf("%s", adJSON)
		var adResult usecases.SearchResponse
		err = json.Unmarshal([]byte(ad), &adResult)
		if err != nil {
			return usecases.SearchResponse{}, fmt.Errorf("There was an error retrieving ads info from search-ms: %+v. \nBody request: %+v", err, ads)
		}
		if len(adResult) == 0 {
			return usecases.SearchResponse{}, fmt.Errorf("The specified ads %+v don't return results from search-ms", ads)
		}
		return adResult, nil
	}
	return usecases.SearchResponse{}, fmt.Errorf("There was an error retrieving ads info: %+v", err)
}
