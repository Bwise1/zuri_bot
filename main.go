package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Bwise1/zuri_bot/twit"
	oauth1Login "github.com/dghubble/gologin/v2/oauth1"
	"github.com/dghubble/gologin/v2/twitter"
	"github.com/dghubble/oauth1"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"

	"github.com/Bwise1/zuri_bot/mongo"
)

type App struct {
	*mux.Router
	*http.Server
	*mongo.DB
}

func main() {
	port := getPort()

	CLUSTER_URL := "mongodb+srv://admin:admin@cluster0.bahi3.mongodb.net/zuri_bot?retryWrites=true&w=majority"
	db, err := mongo.Connect(CLUSTER_URL)
	//db, err := mongo.Connect(os.Getenv("CLUSTER_URL"))
	defer db.Disconnect(context.Background())
	if err != nil {
		panic(err)
	}
	app := NewApp(db)
	log.Fatal(app.Run(port))
}

func NewApp(db *mongo.DB) *App {
	router := mux.NewRouter().StrictSlash(true)

	server := &http.Server{
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	app := &App{
		Router: router,
		Server: server,
		DB:     db,
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

	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "Hello world!")
	})
	// twit.SendTweet(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	router.Handle("/twitter/login", twitter.LoginHandler(oauth1Config, nil))
	router.Handle("/twitter/callback", twitter.CallbackHandler(oauth1Config, twit.IssueSession(a.SaveUserLogin), nil))
	router.HandleFunc("/twitter/post-text", twit.CreateNewTweetText).Methods("POST")
	router.HandleFunc("/twitter/post-media", twit.CreateNewTweetMedia).Methods("POST")

	c := cors.AllowAll().Handler(router)
	a.Server.Handler = handlers.LoggingHandler(os.Stdout, c)
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

func (a *App) SaveUserLogin(w http.ResponseWriter, r *http.Request) {
	collectionName := "twitter_users"
	token, secret, err := oauth1Login.AccessTokenFromContext(r.Context())
	tu, err := twitter.UserFromContext(r.Context())
	if err != nil {
		log.Printf("Error getting user tokens: %v", err)
		return
	}
	coll := a.DB.GetCollection(collectionName)
	doc := map[string]interface{}{
		"token":  token,
		"secret": secret,
		"user":   tu,
	}
	if _, err := coll.InsertOne(r.Context(), doc); err != nil {
		log.Printf("error: %v", err)
	}
}
