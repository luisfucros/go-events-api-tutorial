package main

import (
	"net/http"
	"fmt"
	"time"
)

func (app *application) serve() error {
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", app.config.Server.Port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Infof("starting server on port %d", app.config.Server.Port)

	return server.ListenAndServe()
}