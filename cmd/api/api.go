package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct{
	config config 
}

type config struct{
	addr string 
}

// Mount sets up the routes for the application.
func (app *application) mount() http.Handler{
 	r := chi.NewRouter()

// Middleware base   
    r.Use(middleware.RequestID) // important for rate limiting
	r.Use(middleware.RealIP) // important for rate limiting and analytics and tracing
    r.Use(middleware.Recoverer) // recover from crashes 
    r.Use(middleware.Logger) // log requests 
    
// Set a timeout value on the request context (ctx), that will signal
// through ctx.Done() that the request has timed out and further
// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
    
    r.Route("/v1", func(r chi.Router){
     r.Get("/v1/health", app.healthCheckHandler)
     })
	
	return r
}

// Run starts the server.
func (app *application) run(mux http.Handler) error{
	
	srv := &http.Server{
		Addr: app.config.addr,
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 10 * time.Second,

	}
	
	fmt.Printf("Your server is running on %s", app.config.addr)
	
	return srv.ListenAndServe()
}