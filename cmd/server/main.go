package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/pbaettig/moncron/internal/pkg/handlers"
	"github.com/pbaettig/moncron/internal/pkg/store"
	"github.com/pbaettig/moncron/internal/pkg/store/sqlite"
	log "github.com/sirupsen/logrus"
)

var (
	jobStore store.JobRunStorer
)

func setupServer(addr string, port int, r http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", addr, port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
}

func main() {
	log.SetLevel(log.DebugLevel)
	args := new(cmdlineArgs)
	if err := args.Parse(); err != nil {
		log.Fatalln(err.Error())
	}

	jobStore, err := sqlite.NewDB(args.dbPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("using DB %s", args.dbPath)

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	handlers.RegisterJobRunStorerApiRoutes(apiRouter, jobStore)
	handlers.RegisterJobRunStorerHtmlRoutes(router, jobStore)

	srv := setupServer(args.listenAddress, args.listenPort, handlers.LoggingHandler{gorillahandlers.CompressHandler(router)})
	go func() {
		log.Infof("starting server on %s:%d", args.listenAddress, args.listenPort)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), args.shutdownTimeout)
	defer cancel()

	srv.Shutdown(ctx)

	log.Info("shutting down")
	os.Exit(0)
}
