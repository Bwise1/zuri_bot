package twit

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"

	twitt "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/gologin/v2/twitter"
	"github.com/dghubble/oauth1"
	"github.com/dghubble/sessions"
)

const (
	sessionName     = "example-twtter-app"
	sessionSecret   = "example cookie signing secret"
	sessionUserKey  = "twitterID"
	sessionUsername = "twitterUsername"
)

// sessionStore encodes and decodes session data stored in signed cookies
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

func IssueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		twitterUser, err := twitter.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(twitterUser)
		// 2. Implement a success handler to issue some form of session
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = twitterUser.ID
		session.Values[sessionUsername] = twitterUser.ScreenName
		session.Save(w)
		http.Redirect(w, req, "/", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

func CreateNewTweet(rw http.ResponseWriter, req *http.Request) {
	comp, err := SendTweet(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	println(comp)
	if err != nil {
		fmt.Fprintln(rw, "Hello world!")
	}
	fmt.Fprintln(rw, comp)

}

func RandomString(n int) string {
	var output string

	ascii := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ascii_arr := strings.Split(ascii, "")
	for i := 1; i < n; i++ {
		randInt := rand.Intn(len(ascii_arr))
		output = output + ascii_arr[randInt]
	}
	return output
}

func SendTweet(accessToken string, accessSecret string) (comp bool, err error) {
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitt.NewClient(httpClient)
	tweet, resp, err := client.Statuses.Update("just setting up tinz", nil)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}
	println(tweet, err)
	return true, nil
}
