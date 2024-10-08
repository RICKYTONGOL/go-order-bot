// Package main implements the API server entry point.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dongle/go-order-bot/internal/api/svc"
	v1 "github.com/dongle/go-order-bot/internal/api/v1"
	"github.com/dongle/go-order-bot/internal/config"
	"github.com/dongle/go-order-bot/internal/logger"
	"github.com/dongle/go-order-bot/internal/repository"
	"github.com/dongle/go-order-bot/internal/utilservice"
	"github.com/rs/cors"
)

// apiServer implements the API server application
type apiServer struct {
	cfg *config.Config
	log logger.Logger
	// api          resolvers.ApiResolver
	srv          *http.Server
	isVersionReq bool
	isInit       bool
}

// init initializes the server
func (app *apiServer) init() {
	// make sure to capture version and rescan depth
	flag.BoolVar(&app.isVersionReq, "v", false, "get the server version")
	flag.BoolVar(&app.isInit, "init", false, "initialize api server")

	// get the configuration including parsing the calling flags
	var err error
	app.cfg, err = config.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	// configure logger based on the configuration
	app.log = logger.New(app.cfg)

	svc.SetConfig(app.cfg)
	svc.SetLogger(app.log)
	repository.SetConfig(app.cfg)
	repository.SetLogger(app.log)

	utilservice.SetConfig(app.cfg)
	utilservice.SetLogger(app.log)

	// make the HTTP server
	app.makeHttpServer()
}

// run executes the server function.
func (app *apiServer) run() {
	PrintVersion(app.cfg)
	if app.isVersionReq {
		return
	}

	if app.isInit {
		app.log.Infof("initializing API server")
		app.initApiServer()
		return
	}

	// make sure to capture terminate signals
	app.observeSignals()

	svc.Manager().Run()
	utilservice.Instance().Run()

	app.log.Infof("listening for requests on %s", app.cfg.Server.BindAddress)

	err := app.srv.ListenAndServe()
	if err != nil {
		app.log.Errorf(err.Error())
	}

	app.terminate()
}

// makeHttpServer creates and configures the HTTP server to be used to serve incoming requests
func (app *apiServer) makeHttpServer() {
	// create request MUXer
	srvMux := new(http.ServeMux)

	// create HTTP server to handle our requests
	app.srv = &http.Server{
		Addr:              app.cfg.Server.BindAddress,
		ReadTimeout:       time.Second * time.Duration(app.cfg.Server.ReadTimeout),
		WriteTimeout:      time.Second * time.Duration(app.cfg.Server.WriteTimeout),
		IdleTimeout:       time.Second * time.Duration(app.cfg.Server.IdleTimeout),
		ReadHeaderTimeout: time.Second * time.Duration(app.cfg.Server.HeaderTimeout),
		Handler:           srvMux,
	}

	// setup handlers
	app.setupHandlers(srvMux)
}

// setupHandlers initializes an array of handlers for our HTTP API end-points.
func (app *apiServer) setupHandlers(mux *http.ServeMux) {
	corsHandler := cors.New(corsOptions(app.cfg))
	v1.SetupHandlers(mux, app.cfg, app.log, corsHandler)
}

// corsOptions constructs new set of options for the CORS handler based on provided configuration.
func corsOptions(cfg *config.Config) cors.Options {
	return cors.Options{
		AllowedOrigins:     cfg.Server.CorsOrigin,
		AllowedMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE", "UPDATE"},
		AllowedHeaders:     []string{"Origin", "Accept", "Content-Type", "application/json", "*", "Authorization"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		MaxAge:             300,
	}
}

// observeSignals setups terminate signals observation.
func (app *apiServer) observeSignals() {
	// log what we do
	app.log.Info("os signals captured")

	// make the signal consumer
	ts := make(chan os.Signal, 1)
	signal.Notify(ts, syscall.SIGINT, syscall.SIGTERM)

	// start monitoring
	go func() {
		// wait for the signal
		<-ts

		// terminate HTTP responder
		app.log.Notice("closing HTTP server")
		if err := app.srv.Close(); err != nil {
			app.log.Errorf("could not terminate HTTP listener")
			os.Exit(0)
		}
	}()
}

// initApiServer initialize the api server
func (app *apiServer) initApiServer() {
	// init setting
	app.log.Info("initialize api server")

}

// terminate modules of the server.
func (app *apiServer) terminate() {
	app.log.Notice("closing services")
	if mgr := svc.Manager(); mgr != nil {
		mgr.Close()
	}

	// terminate connections to DB, blockchain, etc.
	app.log.Notice("closing repository")
	if repo := repository.R(); repo != nil {
		repo.Close()
	}

	app.log.Notice("closing utilservice server")
	if utilServiceInstance := utilservice.Instance(); utilServiceInstance != nil {
		utilServiceInstance.Close()
	}
}
