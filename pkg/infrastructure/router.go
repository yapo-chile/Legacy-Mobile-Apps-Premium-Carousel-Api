package infrastructure

import (
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gorilla/context"
	"github.mpi-internal.com/Yapo/goms/pkg/interfaces/handlers"
	"github.mpi-internal.com/Yapo/goms/pkg/interfaces/loggers"
	"gopkg.in/gorilla/mux.v1"
)

// Route stands for an http endpoint description
type Route struct {
	Name      string
	Method    string
	Pattern   string
	Handler   handlers.Handler
	UseCache  bool
	TimeCache time.Duration
}

type routeGroups struct {
	Prefix string
	Groups []Route
}

// WrapperFunc defines a type for functions that wrap an http.HandlerFunc
// to modify its behaviour
type WrapperFunc func(pattern string, handler http.HandlerFunc) http.HandlerFunc

// Routes is an array of routes with a common prefix
type Routes []routeGroups

// RouterMaker gathers route and wrapper information to build a router
type RouterMaker struct {
	Logger        loggers.Logger
	WrapperFuncs  []WrapperFunc
	WithProfiling bool
	Routes        Routes
	Cors          handlers.Cors
	Cache         handlers.Cache
}

// NewRouter setups a Router based on the provided routes
func (maker *RouterMaker) NewRouter() http.Handler {
	router := mux.NewRouter()
	for _, routeGroup := range maker.Routes {
		subRouter := router.PathPrefix(routeGroup.Prefix).Subrouter()
		for _, route := range routeGroup.Groups {
			hLogger := loggers.MakeJSONHandlerLogger(maker.Logger)
			hInputHandler := NewInputHandler()
			cache := &handlers.Cache{}
			if route.UseCache {
				cache.Enabled = maker.Cache.Enabled
				cache.Etag = maker.Cache.Etag
				cache.MaxAge = maker.Cache.MaxAge
				if route.TimeCache > 0 {
					// Replace default max age
					cache.MaxAge = route.TimeCache
				}
			}
			handler := handlers.MakeJSONHandlerFunc(route.Handler, hLogger, hInputHandler, maker.Cors, cache)
			for _, wrapFunc := range maker.WrapperFuncs {
				handler = wrapFunc(route.Pattern, handler)
			}
			subRouter.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)
		}
	}
	if maker.WithProfiling {
		router.HandleFunc("/debug/pprof/", pprof.Index)
		router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		router.HandleFunc("/debug/pprof/trace", pprof.Trace)

		router.Handle("/debug/pprof/block", pprof.Handler("block"))
		router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		router.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
		router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	}
	return context.ClearHandler(router)
}
