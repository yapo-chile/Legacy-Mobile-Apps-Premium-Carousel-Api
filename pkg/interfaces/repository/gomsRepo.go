package repository

import (
	"encoding/json"
	"fmt"

	"github.mpi-internal.com/Yapo/goms/pkg/usecases"
)

// GomsRepository allows the interaction with goms service using http request
type GomsRepository struct {
	Handler HTTPHandler
	Path    string
	TimeOut int
}

// GomsResponse represents the response of the goms microservice
type GomsResponse struct {
	Status string `json:"status"`
}

// NewGomsRepository instanciates a HTTPRepository repository
func NewGomsRepository(handler HTTPHandler, timeOut int, path string) usecases.GomsRepository {
	return &GomsRepository{
		Path:    path,
		Handler: handler,
		TimeOut: timeOut,
	}
}

// GetHealthcheck obtains the status of the goms application
func (repo *GomsRepository) GetHealthcheck() (string, error) {
	request := repo.Handler.NewRequest().
		SetMethod("GET").
		SetPath(repo.Path).
		SetTimeOut(repo.TimeOut)

	response, err := repo.Handler.Send(request)
	if err != nil {
		return "", fmt.Errorf("there was an error obtaining healthcheck from Goms: %+v", err)
	}

	var gomsresp GomsResponse

	if response != "" {
		if err = json.Unmarshal([]byte(response.(string)), &gomsresp); err != nil {
			return "", fmt.Errorf("there was an error parsing goms response: %+v", err)
		}

		return gomsresp.Status, nil
	}

	return "", fmt.Errorf("goms response received is empty")
}
