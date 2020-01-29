package main

import (
	"encoding/json"
	"fmt"
	"os"

	// CLONE REMOVE START
	"regexp"
	"time"

	// CLONE REMOVE END
	"github.mpi-internal.com/Yapo/goms/pkg/infrastructure"
	"github.mpi-internal.com/Yapo/goms/pkg/interfaces/handlers"

	// CLONE REMOVE START
	"github.mpi-internal.com/Yapo/goms/pkg/interfaces/loggers"
	"github.mpi-internal.com/Yapo/goms/pkg/interfaces/repository"
	"github.mpi-internal.com/Yapo/goms/pkg/usecases"
	// CLONE REMOVE END
)

func main() { //nolint: funlen
	var shutdownSequence = infrastructure.NewShutdownSequence()
	var conf infrastructure.Config

	fmt.Printf("Etag:%d\n", conf.CacheConf.InitEtag())
	shutdownSequence.Listen()
	infrastructure.LoadFromEnv(&conf)

	if jconf, err := json.MarshalIndent(conf, "", "    "); err == nil {
		fmt.Printf("Config: \n%s\n", jconf)
	} else {
		fmt.Printf("Config: \n%+v\n", conf)
	}

	fmt.Printf("Setting up Prometheus\n")

	prometheus := infrastructure.MakePrometheusExporter(
		conf.PrometheusConf.Port,
		conf.PrometheusConf.Enabled,
	)

	fmt.Printf("Setting up logger\n")

	logger, err := infrastructure.MakeYapoLogger(&conf.LoggerConf,
		prometheus.NewEventsCollector(
			"goms_service_events_total",
			"events tracker counter for goms service",
		),
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(2) //nolint: gomnd
	}

	shutdownSequence.Push(prometheus)
	logger.Info("Initializing resources")

	// HealthHandler
	var healthHandler handlers.HealthHandler

	// CLONE REMOVE START
	// To handle http connections you can use an httpHandler
	HTTPHandler := infrastructure.NewHTTPHandler(logger)

	// FibonacciHandler
	fibonacciLogger := loggers.MakeFibonacciLogger(logger)
	fibonacciRepository := repository.NewMapFibonacciRepository()
	fibonacciInteractor := usecases.FibonacciInteractor{
		Logger:     fibonacciLogger,
		Repository: fibonacciRepository,
	}
	fibonacciHandler := handlers.FibonacciHandler{
		Interactor: &fibonacciInteractor,
	}

	getUserDataHandlerLogger := loggers.MakeGetUserDataHandlerLogger(logger)

	emailFormat := "^[\\w_+-]+(\\.[\\w_+-]+)*\\.?@([\\w_+-]+\\.)+[\\w]{2,4}$"
	emailValidate := regexp.MustCompile(emailFormat)
	userProfileRepo := repository.NewUserProfileRepository(
		HTTPHandler,
		conf.ProfileConf.Host+conf.ProfileConf.UserDataPath+conf.ProfileConf.UserDataTokens,
	)
	getUserDataInteractor := usecases.GetUserDataInteractor{
		UserProfileRepository: userProfileRepo,
	}

	getUserDataHandler := handlers.GetUserDataHandler{
		Interactor:    &getUserDataInteractor,
		Logger:        getUserDataHandlerLogger,
		EmailValidate: emailValidate,
	}

	// httpCircuitBreakerHandler which retries requests with it's client
	// until it returns a valid answer and then continues normal execution
	// OPTION: classic HTTP

	getHealthLogger := loggers.MakeGomsRepoLogger(logger)
	getHealthInteractor := usecases.GetHealthcheckInteractor{
		GomsRepository: repository.NewGomsRepository(
			HTTPHandler,
			conf.GomsClientConf.TimeOut,
			conf.GomsClientConf.GetHealthcheckPath),
		Logger: getHealthLogger,
	}
	getHealthHandler := handlers.GetHealthcheckHandler{
		GetHealthcheckInteractor: &getHealthInteractor,
	}

	// OPTION: HTTP with Circuit Breaker
	circuitBreaker := infrastructure.NewCircuitBreaker(
		conf.CircuitBreakerConf.Name,
		conf.CircuitBreakerConf.ConsecutiveFailure,
		conf.CircuitBreakerConf.FailureRatio,
		conf.CircuitBreakerConf.Timeout,
		conf.CircuitBreakerConf.Interval,
		logger,
	)
	HTTPCBHandler := infrastructure.NewHTTPCircuitBreakerHandler(circuitBreaker, logger, HTTPHandler)
	getHealthCBHandler := handlers.GetHealthcheckHandler{
		GetHealthcheckInteractor: &usecases.GetHealthcheckInteractor{
			GomsRepository: repository.NewGomsRepository(
				HTTPCBHandler,
				conf.GomsClientConf.TimeOut,
				conf.GomsClientConf.GetHealthcheckPath),
			Logger: getHealthLogger,
		},
	}

	// CLONE REMOVE END

	// CLONE-RCONF REMOVE START
	// Initialize remote conf example
	lastUpdate, errRconf := infrastructure.NewRconf(
		conf.EtcdConf.Host,
		conf.EtcdConf.LastUpdate,
		conf.EtcdConf.Prefix,
		logger,
	)

	if errRconf != nil {
		logger.Error("Error loading remote conf")
	} else {
		logger.Info("Remote Conf Updated at %s", lastUpdate.Content.Node.Value)
	}
	// CLONE-RCONF REMOVE END

	useBrowserCache := handlers.Cache{
		MaxAge:  conf.CacheConf.MaxAge,
		Etag:    conf.CacheConf.Etag,
		Enabled: conf.CacheConf.Enabled,
	}
	// Setting up router
	maker := infrastructure.RouterMaker{
		Logger:        logger,
		Cors:          conf.CorsConf,
		Cache:         useBrowserCache,
		WrapperFuncs:  []infrastructure.WrapperFunc{prometheus.TrackHandlerFunc},
		WithProfiling: conf.ServiceConf.Profiling,
		Routes: infrastructure.Routes{
			{
				// This is the base path, all routes will start with this prefix
				Prefix: "",
				Groups: []infrastructure.Route{
					{
						Name:    "Check service health",
						Method:  "GET",
						Pattern: "/healthcheck",
						Handler: &healthHandler,
					},
					// CLONE REMOVE START
					{
						Name:      "Retrieve the Nth Fibonacci with Clean Architecture",
						Method:    "GET",
						Pattern:   "/fibonacci",
						Handler:   &fibonacciHandler,
						UseCache:  true,
						TimeCache: 60 * time.Minute, //nolint: gomnd
					},
					{
						Name:    "Retrieve healthcheck by doing a client request to itself",
						Method:  "GET",
						Pattern: "/http/healthcheck",
						Handler: &getHealthHandler,
					},
					{
						Name:    "Retrieve healthcheck by doing a client request to itself using Circuit Breaker",
						Method:  "GET",
						Pattern: "/httpcb/healthcheck",
						Handler: &getHealthCBHandler,
					},
					{
						Name:    "Retrieve the user basic data",
						Method:  "GET",
						Pattern: "/user/basic-data",
						Handler: &getUserDataHandler,
					},
					// CLONE REMOVE END
				},
			},
		},
	}

	router := maker.NewRouter()

	server := infrastructure.NewHTTPServer(
		fmt.Sprintf("%s:%d", conf.Runtime.Host, conf.Runtime.Port),
		router,
		logger,
	)
	shutdownSequence.Push(server)
	logger.Info("Starting request serving")

	go server.ListenAndServe()
	shutdownSequence.Wait()
	logger.Info("Server exited normally")
}
