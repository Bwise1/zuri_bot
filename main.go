package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

type App struct {
	*mux.Router
	*http.Server
	// Twitter credentials
}

func main() {
	port := getPort()
	app := NewApp()
	log.Fatal(app.Run(port))
}

func NewApp() *App {
	router := mux.NewRouter().StrictSlash(true)
	server := &http.Server{
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	app := &App{
		Router: router,
		Server: server,
	}
	app.RegisterRoutes()
	return app
}

func (a *App) RegisterRoutes() {
	router := a.Router

	router.HandleFunc("/", func(rw http.ResponseWriter,
		r *http.Request) {
		fmt.Fprintln(rw, "Hello world!")
	})

	a.Handler = handlers.LoggingHandler(os.Stdout, router)
}

func (a *App) Run(port ...string) error {
	if len(port) < 1 {
		a.Addr = ":8081"
	} else {
		a.Addr = ":" + port[0]
	}
	log.Printf("Server running on %s\n", a.Server.Addr)
	return a.ListenAndServe()
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	return port
}
