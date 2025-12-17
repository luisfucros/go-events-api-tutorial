package main

import (
	"net/http"
	"fmt"
	"time"
)

func (app *application) serve() error {
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", app.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Infof("starting server on port %d", app.port)

	return server.ListenAndServe()
}