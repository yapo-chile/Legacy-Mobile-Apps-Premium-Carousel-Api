package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/infrastructure"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/interfaces/handlers"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/interfaces/repository"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
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
			"premium-carousel-api_service_events_total",
			"events tracker counter for premium-carousel-api service",
		),
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(2) //nolint: gomnd
	}

	shutdownSequence.Push(prometheus)
	logger.Info("Initializing resources")

	regions, errorRegions := infrastructure.NewEtcd(
		conf.EtcdConf.Host,
		conf.EtcdConf.Region,
		conf.EtcdConf.Prefix,
		logger,
	)

	if errorRegions != nil {
		panic("Unable to load regions remote config from etcd")
	}

	HTTPHandler := infrastructure.NewHTTPHandler(logger)

	elasticsearch := infrastructure.NewElasticsearch("http://10.15.1.78", "19200", logger)

	adRepo := repository.MakeAdRepository(
		elasticsearch,
		regions,
		conf.AdConf.Host+conf.AdConf.Path,
		conf.AdConf.MaxAdsToDisplay,
	)

	userProfileRepo := repository.MakeUserProfileRepository(
		HTTPHandler,
		conf.ProfileConf.Host+conf.ProfileConf.UserDataPath+conf.ProfileConf.UserDataTokens,
	)

	configRepo := repository.MakeConfigRepository(nil)

	getUserAdsInteractor := usecases.MakeGetUserAdsInteractor(adRepo, configRepo)
	getUserDataInteractor := usecases.MakeGetUserDataInteractor(userProfileRepo)

	// UserAdsHandler
	getUserAdsHandler := handlers.GetUserAdsHandler{
		Interactor:            getUserAdsInteractor,
		GetUserDataInteractor: getUserDataInteractor,
		Logger:                nil,
		UnitOfAccountSymbol:   conf.AdConf.UnitOfAccountSymbol,
		CurrencySymbol:        conf.AdConf.CurrencySymbol,
	}

	// HealthHandler
	var healthHandler handlers.HealthHandler

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
				Groups: []infrastructure.Route{
					{
						Name:    "Check service health",
						Method:  "GET",
						Pattern: "/healthcheck",
						Handler: &healthHandler,
					},
					{
						Name:    "Get user ads",
						Method:  "GET",
						Pattern: "/ads/{token:.*}",
						Handler: &getUserAdsHandler,
					},
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
