package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.elastic.co/apm/module/apmgorilla"
	"go.elastic.co/apm/module/apmhttp"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// HTTPRoute model for http routes
type HTTPRoute struct {
	Method string
	Path   string
}

// LOGGER initializing the logger
var LOGGER *zap.SugaredLogger

// HTTPController defined http controllers
type HTTPController interface {
	GET(http.ResponseWriter, *http.Request)
	LIST(http.ResponseWriter, *http.Request)
	GetPaths() []HTTPRoute
}

// FrontierWebServer frontier webserver model and required params
type FrontierWebServer struct {
	host string
	port int

	routes     *mux.Router
	srv        *http.Server
	rootRouter *mux.Router
	log        *zap.SugaredLogger

}

// NewFrontierWebServer Frontier web server required parameters initialization
func NewFrontierWebServer(logger *zap.SugaredLogger, host string, port int) *FrontierWebServer{
	LOGGER = logger

	rootRouter := mux.NewRouter()
	rootRouter.Use(apmgorilla.Middleware())

	ws := &FrontierWebServer{
		host: host, port: port, log: logger,
		routes: rootRouter.PathPrefix("/api/v1").Subrouter(),
	}
	rootRouter.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	ws.rootRouter = rootRouter
	return ws
}

/*Start Method  reads the parameters from configuration file and assign those parameters to  http server to create the
    http server for frontier service , and finally ListenAndServer method starts the http port to receive http requests
    by listening to tcp
 */
func (f *FrontierWebServer) Start() error{
	addr := fmt.Sprintf("%v:%v", f.host, f.port)
	f.log.Infof("Starting web server @ -> Host: %v, Port: %v", f.host, f.port)
	f.srv = &http.Server{
		Handler: apmhttp.Wrap(f.rootRouter),
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return f.srv.ListenAndServe()
}

func (f *FrontierWebServer) Stop() error {
  return f.srv.Close()
}

// SetRoute routing the http methods
func (f *FrontierWebServer) SetRoute(ctrl HTTPController)  {
	for _, route := range ctrl.GetPaths() {
		switch route.Method {
		case "GET":
			f.routes.HandleFunc(route.Path, ctrl.GET).Methods(route.Method)
		case "LIST":
			f.routes.HandleFunc(route.Path, ctrl.LIST).Methods("GET")
		}
	}
}