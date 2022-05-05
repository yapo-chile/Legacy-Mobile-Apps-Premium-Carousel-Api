package infrastructure

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ServiceConf holds configuration for this Service
type ServiceConf struct {
	Host      string `env:"HOST" envDefault:":8080"`
	Profiling bool   `env:"PROFILING" envDefault:"true"`
}

// LoggerConf holds configuration for logging
// LogLevel definition:
//   0 - Debug
//   1 - Info
//   2 - Warning
//   3 - Error
//   4 - Critic
type LoggerConf struct {
	SyslogIdentity string `env:"SYSLOG_IDENTITY"`
	SyslogEnabled  bool   `env:"SYSLOG_ENABLED" envDefault:"false"`
	StdlogEnabled  bool   `env:"STDLOG_ENABLED" envDefault:"true"`
	LogLevel       int    `env:"LOG_LEVEL" envDefault:"2"`
}

// PrometheusConf holds configuration to report to Prometheus
type PrometheusConf struct {
	Port    string `env:"PORT" envDefault:"8877"`
	Enabled bool   `env:"ENABLED" envDefault:"false"`
}

// RuntimeConfig config to start the app
type RuntimeConfig struct {
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	Port int    `env:"PORT" envDefault:"8080"`
}

// KafkaProducerConf holds configurations to connect and produce to kafka
type KafkaProducerConf struct {
	Host              string `env:"HOST"`
	Port              int    `env:"PORT" envDefault:"9092"`
	Acks              string `env:"ACKS" envDefault:"1"`
	CompressionType   string `env:"COMPRESSION_TYPE" envDefault:"none"`
	Retries           int    `env:"RETRIES" envDefault:"10"`
	LingerMS          int    `env:"LINGER_MS" envDefault:"0"`
	RequestTimeoutMS  int    `env:"REQUEST_TIMEOUT_MS" envDefault:"30000"`
	EnableIdempotence bool   `env:"ENABLE_IDEMPOTENCE" envDefault:"false"`
}

// BackendEventsConf holds backend events configurations
type BackendEventsConf struct {
	PremiumProductsTopic string `env:"PREMIUM_PRODUCTS_TOPIC" envDefault:"premium_product"`
	Enabled              bool   `env:"ENABLED" envDefault:"false"`
}

// DatabaseConf holds configuration for postgres database connection
type DatabaseConf struct {
	Host        string `env:"HOST" envDefault:"db"`
	Port        int    `env:"PORT" envDefault:"5432"`
	Dbname      string `env:"NAME" envDefault:"pgdb"`
	DbUser      string `env:"USER" envDefault:"postgres"`
	DbPasswd    string `env:"PASSWORD" envDefault:"postgres"`
	Sslmode     string `env:"SSL_MODE" envDefault:"disable"`
	MaxIdle     int    `env:"MAX_IDLE" envDefault:"10"`
	MaxOpen     int    `env:"MAX_OPEN" envDefault:"100"`
	MgFolder    string `env:"MIGRATIONS_FOLDER" envDefault:"migrations"`
	MgDriver    string `env:"MIGRATIONS_DRIVER" envDefault:"postgres"`
	ConnRetries int    `env:"CONN_RETRIES" envDefault:"3"`
}

// GomsClientConf holds configuration regarding to our http client (premium-carousel-api itself in this case)
type GomsClientConf struct {
	TimeOut            int    `env:"TIMEOUT" envDefault:"30"`
	GetHealthcheckPath string `env:"HEALTH_PATH" envDefault:"/get/healthcheck"`
}

// EtcdConf configure how to read configuration from remote Etcd service
type EtcdConf struct {
	Host       string `env:"HOST" envDefault:"http://lb:2397"`
	LastUpdate string `env:"LAST_UPDATE" envDefault:"/last_update"`
	Prefix     string `env:"PREFIX" envDefault:"/v2/keys"`
	RegionPath string `env:"REGION_PATH" envDefault:"/public/location/regions.json"`
}

// CorsConf holds cors headers
type CorsConf struct {
	Enabled bool   `env:"ENABLED" envDefault:"false"`
	Origin  string `env:"ORIGIN" envDefault:"*"`
	Methods string `env:"METHODS" envDefault:"GET, OPTIONS"`
	Headers string `env:"HEADERS" envDefault:"Accept,Content-Type,Content-Length,If-None-Match,Accept-Encoding,User-Agent"`
}

// GetHeaders return map of cors used
func (cc CorsConf) GetHeaders() map[string]string {
	if !cc.Enabled {
		return map[string]string{}
	}

	return map[string]string{
		"Origin":  cc.Origin,
		"Methods": cc.Methods,
		"Headers": cc.Headers,
	}
}

// BrowserCacheConf Used to handle browser cache
type BrowserCacheConf struct {
	Enabled bool `env:"ENABLED" envDefault:"false"`
	//Cache max age in secs(use browser cache)
	MaxAge time.Duration `env:"MAX_AGE" envDefault:"720h"`
	Etag   int64
}

// CacheConf holds cache configurations
type CacheConf struct {
	Host       string        `env:"HOST" envDefault:"cache:6379"`
	Prefix     string        `env:"PREFIX" envDefault:"cache"`
	Password   string        `env:"PASSWORD"`
	DB         int           `env:"DB"`
	DefaultTTL time.Duration `env:"DEFAULT_TTL" envDefault:"1h"`
}

// ControlPanelConf holds Control Panel configurations
type ControlPanelConf struct {
	ResultsPerPage int `env:"RESULTS_PER_PAGE" envDefault:"50"`
}

// InitEtag use current epoc to config etag
func (chc *BrowserCacheConf) InitEtag() int64 {
	chc.Etag = time.Now().Unix()
	return chc.Etag
}

// AdConf contains search-ms configuration params
type AdConf struct {
	Host                string `env:"HOST" envDefault:"http://10.15.1.78"`
	Port                string `env:"PORT" envDefault:"19200"`
	Username            string `env:"USERNAME" envDefault:"user"`
	Password            string `env:"PASSSWORD" envDefault:"pass"`
	Index               string `env:"PATH" envDefault:"ads"`
	ImageServerURL      string `env:"IMAGE_SERVER_URL" envDefault:"https://img.yapo.cl/%s/%s/%s.jpg"`
	CurrencySymbol      string `env:"CURRENCY_SYMBOL" envDefault:"$"`
	UnitOfAccountSymbol string `env:"UNIT_OF_ACCOUNT_SYMBOL" envDefault:"UF"`
	MaxAdsToDisplay     int    `env:"MAX_ADS_TO_DISPLAY" envDefault:"15"`
	MinAdsToDisplay     int    `env:"MIN_ADS_TO_DISPLAY" envDefault:"2"`
}

// Config holds all configuration for the service
type Config struct {
	ServiceConf       ServiceConf       `env:"SERVICE_"`
	PrometheusConf    PrometheusConf    `env:"PROMETHEUS_"`
	LoggerConf        LoggerConf        `env:"LOGGER_"`
	Runtime           RuntimeConfig     `env:"APP_"`
	GomsClientConf    GomsClientConf    `env:"GOMS_"`
	EtcdConf          EtcdConf          `env:"ETCD_"`
	CorsConf          CorsConf          `env:"CORS_"`
	BrowserCacheConf  BrowserCacheConf  `env:"BROWSER_CACHE_"`
	CacheConf         CacheConf         `env:"CACHE_"`
	DatabaseConf      DatabaseConf      `env:"DATABASE_"`
	AdConf            AdConf            `env:"AD_"`
	ControlPanelConf  ControlPanelConf  `env:"CP_"`
	KafkaProducerConf KafkaProducerConf `env:"KAFKA_PRODUCER_"`
	BackendEventsConf BackendEventsConf `env:"BACKEND_EVENTS_"`
}

// LoadFromEnv loads the config data from the environment variables
func LoadFromEnv(data interface{}) {
	load(reflect.ValueOf(data), "", "")
}

// valueFromEnv lookup the best value for a variable on the environment
func valueFromEnv(envTag, envDefault string) string {
	// Maybe it's a secret and <envTag>_FILE points to a file with the value
	// https://rancher.com/docs/rancher/v1.6/en/cattle/secrets/#docker-hub-images
	if fileName, ok := os.LookupEnv(fmt.Sprintf("%s_FILE", envTag)); ok {
		// filepath.Clean() will clean the input path and remove some unnecessary things
		// like multiple separators doble "." and others
		// if for some reason you are having troubles reaching your file, check the
		// output of the Clean function and test if its what you expect
		// you can find more info here: https://golang.org/pkg/path/filepath/#Clean
		b, err := ioutil.ReadFile(filepath.Clean(fileName))
		if err == nil {
			return string(b)
		}

		fmt.Print(err)
	}
	// The value might be set directly on the environment
	if value, ok := os.LookupEnv(envTag); ok {
		return value
	}
	// Nothing to do, return the default
	return envDefault
}

// load the variable defined in the envTag into Value
func load(conf reflect.Value, envTag, envDefault string) { //nolint: gocyclo, gocognit
	if conf.Kind() == reflect.Ptr {
		reflectedConf := reflect.Indirect(conf)
		// Only attempt to set writeable variables
		if reflectedConf.IsValid() && reflectedConf.CanSet() {
			value := valueFromEnv(envTag, envDefault)
			// Print message if config is missing
			if envTag != "" && value == "" && !strings.HasSuffix(envTag, "_") {
				fmt.Printf("Config for %s missing\n", envTag)
			}

			switch reflectedConf.Interface().(type) {
			case int:
				if value, err := strconv.ParseInt(value, 10, 32); err == nil {
					reflectedConf.Set(reflect.ValueOf(int(value)))
				}
			case int64:
				if value, err := strconv.ParseInt(value, 10, 64); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			case uint32:
				if value, err := strconv.ParseUint(value, 10, 32); err == nil {
					reflectedConf.Set(reflect.ValueOf(uint32(value)))
				}
			case float64:
				if value, err := strconv.ParseFloat(value, 64); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			case string:
				reflectedConf.Set(reflect.ValueOf(value))
			case bool:
				if value, err := strconv.ParseBool(value); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			case time.Time:
				if value, err := time.Parse(time.RFC3339, value); err == nil {
					reflectedConf.Set(reflect.ValueOf(value))
				}
			case time.Duration:
				if t, err := time.ParseDuration(value); err == nil {
					reflectedConf.Set(reflect.ValueOf(t))
				}
			}

			if reflectedConf.Kind() == reflect.Struct {
				// Recursively load inner struct fields
				for i := 0; i < reflectedConf.NumField(); i++ {
					if tag, ok := reflectedConf.Type().Field(i).Tag.Lookup("env"); ok {
						def, _ := reflectedConf.Type().Field(i).Tag.Lookup("envDefault")
						load(reflectedConf.Field(i).Addr(), envTag+tag, def)
					}
				}
			}
		}
	}
}
