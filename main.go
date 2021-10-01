package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Bwise1/zuri_bot/twit"

	"github.com/dghubble/gologin/v2/twitter"
	"github.com/dghubble/oauth1"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"
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
	oauth1Config := &oauth1.Config{
		ConsumerKey:    os.Getenv("CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("CONSUMER_SECRET"),
		CallbackURL:    "https://zuri-bot.herokuapp.com/twitter/callback",
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
	}
	router.HandleFunc("/", func(rw http.ResponseWriter,
		r *http.Request) {
		fmt.Fprintln(rw, "Hello world!")
	})
	// twit.SendTweet(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	router.Handle("/twitter/login", twitter.LoginHandler(oauth1Config, nil))
	router.Handle("/twitter/callback", twitter.CallbackHandler(oauth1Config, twit.IssueSession(), nil))
	router.HandleFunc("/twitter/post-text", twit.CreateNewTweetText).Methods("POST")
	router.HandleFunc("/twitter/post-media", twit.CreateNewTweetMedia).Methods("POST")

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
