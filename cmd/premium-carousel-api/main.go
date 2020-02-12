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
			"premium_carousel_api_service_events_total",
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
		conf.EtcdConf.RegionPath,
		conf.EtcdConf.Prefix,
		logger,
	)

	if errorRegions != nil {
		panic("Unable to load regions remote config from etcd")
	}

	elasticsearch := infrastructure.NewElasticsearch(conf.AdConf.Host, conf.AdConf.Port, logger)

	adRepo := repository.MakeAdRepository(
		elasticsearch,
		regions,
		conf.AdConf.Index,
		conf.AdConf.ImageServerURL,
		conf.AdConf.MaxAdsToDisplay,
	)

	configRepo := repository.MakeConfigRepository(nil) // TODO

	getUserAdsInteractor := usecases.MakeGetUserAdsInteractor(adRepo, configRepo)
	getAdInteractor := usecases.MakeGetAdInteractor(adRepo)

	// UserAdsHandler
	getUserAdsHandler := handlers.GetUserAdsHandler{
		Interactor:          getUserAdsInteractor,
		GetAdInteractor:     getAdInteractor,
		Logger:              nil, // TODO
		UnitOfAccountSymbol: conf.AdConf.UnitOfAccountSymbol,
		CurrencySymbol:      conf.AdConf.CurrencySymbol,
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
						Pattern: "/ads/{listID:[0-9]+}",
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
