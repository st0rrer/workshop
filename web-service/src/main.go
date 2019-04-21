package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/st0rrer/workshop/web-service/src/config"
	"github.com/st0rrer/workshop/web-service/src/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func encodingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Encoding") == "gzip" {
			r.Body = &gzipReader{
				body: r.Body,
			}
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	logger := log.New(os.Stderr, "main ", log.LstdFlags)

	r := mux.NewRouter()

	r.HandleFunc("/api/{report:(?:activity|visit)}/v1", handler.GetReportHandler()).
		Methods("POST").
		Headers("Content-Type", "application/json",
			"Content-Encoding", "gzip")

	r.Use(encodingMiddleware)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", config.GetConfig().Port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      r,
	}

	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatalln(err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		srv.Shutdown(ctx)
		cancelFunc()
		os.Exit(0)
	}()

	<-ctx.Done()
	logger.Println("Gracefully stop application")
}
