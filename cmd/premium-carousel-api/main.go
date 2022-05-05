package main

import (
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	mpgsql "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"

	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/infrastructure"
	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/interfaces/handlers"
	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/interfaces/loggers"
	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/interfaces/repository"
	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/usecases"
)

func main() { //nolint: funlen
	var shutdownSequence = infrastructure.NewShutdownSequence()
	var conf infrastructure.Config

	fmt.Printf("Etag:%d\n", conf.BrowserCacheConf.InitEtag())
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
		panic(fmt.Errorf("error starting prometheus & loggers: %v", err))

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

	redisHandler := infrastructure.NewRedisHandler(
		conf.CacheConf.Host,
		conf.CacheConf.Password,
		conf.CacheConf.DB,
		logger,
	)

	dbHandler, err := infrastructure.MakePgsqlHandler(conf.DatabaseConf, logger)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect with postgres database: %+v", err))
	}
	shutdownSequence.Push(dbHandler)

	setupMigrations(conf, dbHandler, logger)

	elasticsearch := infrastructure.NewElasticsearch(
		conf.AdConf.Host,
		conf.AdConf.Port,
		conf.AdConf.Username,
		conf.AdConf.Password,
		logger,
	)
	var backendEventsProducer repository.KafkaProducer
	var backendEventsRepository usecases.BackendEventsRepository
	if conf.BackendEventsConf.Enabled {
		backendEventsProducer, err = infrastructure.NewKafkaProducer(
			conf.KafkaProducerConf.Host,
			conf.KafkaProducerConf.Port,
			conf.KafkaProducerConf.Acks,
			conf.KafkaProducerConf.CompressionType,
			conf.KafkaProducerConf.Retries,
			conf.KafkaProducerConf.LingerMS,
			conf.KafkaProducerConf.RequestTimeoutMS,
			conf.KafkaProducerConf.EnableIdempotence,
		)
		if err != nil {
			panic(fmt.Errorf("Error starting kafka producer: %+v", err))
		}
		shutdownSequence.Push(backendEventsProducer)
		backendEventsRepository = repository.MakeBackendEventsProducer(
			backendEventsProducer,
			conf.BackendEventsConf.PremiumProductsTopic,
		)
	}

	adRepo := repository.MakeAdRepository(
		elasticsearch,
		regions,
		conf.AdConf.Index,
		conf.AdConf.ImageServerURL,
		conf.AdConf.MaxAdsToDisplay,
	)

	cacheRepo := repository.NewCacheRepository(
		redisHandler,
		conf.CacheConf.Prefix,
		conf.CacheConf.DefaultTTL,
	)

	purchaseRepo := repository.MakePurchaseRepository(
		dbHandler,
	)

	productRepo := repository.MakeProductRepository(
		dbHandler,
		conf.ControlPanelConf.ResultsPerPage,
		loggers.MakeProductRepositoryLogger(logger),
	)

	getUserAdsInteractor := usecases.MakeGetUserAdsInteractor(
		adRepo,
		productRepo,
		cacheRepo,
		loggers.MakeGetUserAdsLogger(logger),
		conf.CacheConf.DefaultTTL,
		conf.AdConf.MinAdsToDisplay,
	)

	getAdInteractor := usecases.MakeGetAdInteractor(
		adRepo,
		cacheRepo,
		loggers.MakeGetAdLogger(logger),
		conf.CacheConf.DefaultTTL,
	)

	addUserProductInteractor := usecases.MakeAddUserProductInteractor(
		productRepo,
		purchaseRepo,
		cacheRepo,
		loggers.MakeAddUserProductLogger(logger),
		conf.CacheConf.DefaultTTL,
		backendEventsRepository,
		conf.BackendEventsConf.Enabled,
	)

	setPartialConfigInteractor := usecases.MakeSetPartialConfigInteractor(
		productRepo,
		cacheRepo,
		loggers.MakeSetPartialConfigLogger(logger),
		conf.CacheConf.DefaultTTL,
	)

	setConfigInteractor := usecases.MakeSetConfigInteractor(
		productRepo,
		cacheRepo,
		loggers.MakeSetConfigLogger(logger),
		conf.CacheConf.DefaultTTL,
	)

	getUserProductsInteractor := usecases.MakeGetUserProductsInteractor(
		productRepo,
		loggers.MakeGetUserProductsLogger(logger),
	)

	getReportInteractor := usecases.MakeGetReportInteractor(
		productRepo,
		loggers.MakeGetReportLogger(logger),
	)

	expireProductsInteractor := usecases.MakeExpireProductsInteractor(
		productRepo,
		loggers.MakeExpireProductsLogger(logger),
	)

	// UserAdsHandler
	getUserAdsHandler := handlers.GetUserAdsHandler{
		Interactor:          getUserAdsInteractor,
		GetAdInteractor:     getAdInteractor,
		UnitOfAccountSymbol: conf.AdConf.UnitOfAccountSymbol,
		CurrencySymbol:      conf.AdConf.CurrencySymbol,
	}

	addUserProductHandler := handlers.AddUserProductHandler{
		Interactor: addUserProductInteractor,
	}

	getUserProductsHandler := handlers.GetUserProductsHandler{
		Interactor: getUserProductsInteractor,
	}

	getReportHandler := handlers.GetReportHandler{
		Interactor: getReportInteractor,
	}

	setPartialConfigHandler := handlers.SetPartialConfigHandler{
		Interactor: setPartialConfigInteractor,
	}

	setConfigHandler := handlers.SetConfigHandler{
		Interactor: setConfigInteractor,
	}

	expireProductsHandler := handlers.ExpireProductsHandler{
		Interactor: expireProductsInteractor,
	}

	// HealthHandler
	var healthHandler handlers.HealthHandler

	useBrowserCache := handlers.Cache{
		MaxAge:  conf.BrowserCacheConf.MaxAge,
		Etag:    conf.BrowserCacheConf.Etag,
		Enabled: conf.BrowserCacheConf.Enabled,
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
						Pattern: "/related/{listID:[0-9]+}",
						Handler: &getUserAdsHandler,
					},
					{
						Name:    "Add product",
						Method:  "POST",
						Pattern: "/assigns",
						Handler: &addUserProductHandler,
					},
					{
						Name:    "Get user products",
						Method:  "GET",
						Pattern: "/assigns",
						Handler: &getUserProductsHandler,
					},
					{
						Name:    "Set user product config",
						Method:  "PUT",
						Pattern: "/assigns/{ID:[0-9]+}",
						Handler: &setConfigHandler,
					},
					{
						Name:    "Set partial user product config",
						Method:  "PATCH",
						Pattern: "/assigns/{ID:[0-9]+}",
						Handler: &setPartialConfigHandler,
					},
					{
						Name:    "Get report",
						Method:  "GET",
						Pattern: "/report",
						Handler: &getReportHandler,
					},
					{
						Name:    "Expire products",
						Method:  "GET",
						Pattern: "/expire-products",
						Handler: &expireProductsHandler,
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

// Autoexecute database migrations
func setupMigrations(conf infrastructure.Config, dbHandler *infrastructure.PgsqlHandler, logger loggers.Logger) {
	driver, err := mpgsql.WithInstance(dbHandler.Conn, &mpgsql.Config{})
	if err != nil {
		logger.Error("Error to instance migrations %v", err)
		return
	}
	mig, err := migrate.NewWithDatabaseInstance(
		"file://"+conf.DatabaseConf.MgFolder,
		conf.DatabaseConf.MgDriver,
		driver,
	)
	if err != nil {
		logger.Error("Consume migrations sources err %#v", err)
		return
	}
	version, _, _ := mig.Version()
	logger.Info("Migrations Actual Version %d", version)
	err = mig.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Info("Migration message: %v", err)
		return
	}
	version, _, _ = mig.Version()
	logger.Info("Migrations upgraded to version %d", version)
}
