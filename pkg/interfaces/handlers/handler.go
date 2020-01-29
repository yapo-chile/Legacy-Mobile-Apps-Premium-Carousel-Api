package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Yapo/goutils"
)

// HandlerInput is a placeholder for whatever input a handler may need.
type HandlerInput interface{}

// InputGetter defines a type for all functions that, when called, will attempt
// to retrieve and parse the input of a request and return it. Should any error
// happen, a goutils.Response must be filled with an adequate message and code
type InputGetter func() (HandlerInput, *goutils.Response)

// Handler is the interface for the objects that should process web requests.
// Input() must return a fresh struct to be filled with the request input
// Execute(input) receives a filled input struct to handle the request
type Handler interface {
	// Input should return a pointer to the struct that this handler will need
	// to be filled with the user input for a request
	Input(InputRequest) HandlerInput
	// Execute is the actual handler code. The InputGetter can be used to retrieve
	// the request's input at any time (or not at all).
	Execute(InputGetter) *goutils.Response
}

// InputHandler defines what methods an input handler should have
type InputHandler interface {
	NewInputRequest(*http.Request) InputRequest
	SetInputRequest(InputRequest, HandlerInput)
	Input() (HandlerInput, *goutils.Response)
}

// InputRequest defines what methods an input handler should have
type InputRequest interface {
	Set(interface{}) TargetRequest
}

// TargetRequest defines what methods an output request should have
type TargetRequest interface {
	FromJSONBody() TargetRequest
	FromRawBody() TargetRequest
	FromPath() TargetRequest
	FromQuery() TargetRequest
	FromHeaders() TargetRequest
	FromCookies() TargetRequest
	FromForm() TargetRequest
}

// Cors methods to configure cache and cors
type Cors interface {
	// GetHeaders should return the map of headers using key > value format
	GetHeaders() map[string]string
}

// Cache used to work with
type Cache struct {
	// MaxAge is used to know how much time the response is valid at
	// browser level
	MaxAge time.Duration
	// Etag contains the identifier of current running version
	Etag int64
	// Enable allows use or ignore the feature
	Enabled bool
}

// MakeJSONHandlerFunc wraps a Handler on a json-over-http context, returning
// a standard http.HandlerFunc
func MakeJSONHandlerFunc(h Handler, l JSONHandlerLogger, ih InputHandler, crs Cors, cache *Cache) http.HandlerFunc {
	jh := jsonHandler{handler: h, logger: l, inputHandler: ih, cors: crs, cache: cache}
	return jh.run
}

// JSONHandlerLogger defines all the events a jsonHandler can report
type JSONHandlerLogger interface {
	LogRequestStart(r *http.Request)
	LogRequestEnd(*http.Request, *goutils.Response)
	LogRequestPanic(*http.Request, *goutils.Response, interface{})
}

// jsonHandler provides an http.HandlerFunc that reads its input and formats
// its output as json
type jsonHandler struct {
	handler      Handler
	logger       JSONHandlerLogger
	inputHandler InputHandler
	cors         Cors
	cache        *Cache
}

func (jh *jsonHandler) setupCors(w *http.ResponseWriter) {
	for key, value := range jh.cors.GetHeaders() {
		(*w).Header().Set("Access-Control-Allow-"+key, value)
	}
}

func (jh *jsonHandler) inBrowserCache(w http.ResponseWriter, r *http.Request) bool {
	if jh.cache.Enabled {
		key := strconv.FormatInt(jh.cache.Etag, 10)
		seconds := fmt.Sprintf("%.0f", jh.cache.MaxAge.Seconds())
		e := `"` + key + `"`
		w.Header().Set("Etag", e)
		w.Header().Set("Cache-Control", "max-age="+seconds)
		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, e) {
				return true
			}
		}
	}
	return false
}

// run will prepare the input for the actual handler and format the response
// as json. Also, request information will be logged. It's an instance of
// http.HandlerFunc
func (jh *jsonHandler) run(w http.ResponseWriter, r *http.Request) {
	jh.logger.LogRequestStart(r)
	jh.setupCors(&w)
	// Default response
	response := &goutils.Response{
		Code: http.StatusInternalServerError,
	}
	// Function the request can call to retrieve its input
	ri := jh.inputHandler.NewInputRequest(r)
	input := jh.handler.Input(ri)
	jh.inputHandler.SetInputRequest(ri, input)
	// Format the output and send it down the writer
	outputWriter := func() {
		goutils.CreateJSON(response)
		goutils.WriteJSONResponse(w, response)
	}
	// Handle panicking handlers and report errors
	errorHandler := func() {
		if err := recover(); err != nil {
			jh.logger.LogRequestPanic(r, response, err)
		}
	}
	// Setup before calling the actual handler
	defer outputWriter()
	defer errorHandler()

	if jh.inBrowserCache(w, r) {
		response = &goutils.Response{
			Code: http.StatusNotModified,
		}
	} else {
		// Do the Harlem Shake
		response = jh.handler.Execute(jh.inputHandler.Input)
	}
	jh.logger.LogRequestEnd(r, response)
}
