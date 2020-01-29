package usecases

// GetHealthcheckInteractor contains the interfaces that allows it to interact with:
// GomsRepository: the service it needs to consume in order to know it's state
// Logger: to log useful information
type GetHealthcheckInteractor struct {
	GomsRepository GomsRepository
	Logger         HealthcheckPrometheusLogger
}

// HealthcheckPrometheusLogger defines all the events a GetHealthcheckInteractor may
// need/like to report as they happen
type HealthcheckPrometheusLogger interface {
	LogURI(string)
	LogRequestErr(error)
	LogHealthcheckOK(string)
}

// GetHealthcheck allows the service to ask for it's own service state via http
func (i *GetHealthcheckInteractor) GetHealthcheck() (string, error) {
	i.Logger.LogURI("Getting local healthcheck information")
	resp, err := i.GomsRepository.GetHealthcheck()

	if err != nil {
		i.Logger.LogRequestErr(err)
		return "", err
	}

	i.Logger.LogHealthcheckOK("Goms healthcheck answered successfully")

	return resp, nil
}
