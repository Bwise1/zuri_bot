package twit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

var allowedMimeTypes = []string{"application/pdf",
	"image/png", "image/jpg", "text/plain", "image/jpeg",
	"video/mp4", "video/mpeg", "video/ogg", "video/quicktime",
	"application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"application/vnd.android.package-archive", "application/octet-stream",
	"application/x-rar-compressed", " application/octet-stream", " application/zip", "application/octet-stream", "application/x-zip-compressed", "multipart/x-zip",
}

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

func CreateNewTweetText(rw http.ResponseWriter, req *http.Request) {
	var bod map[string]string
	err := json.NewDecoder(req.Body).Decode(&bod)
	if err != nil {
		rw.WriteHeader(400)
		return
	}
	if !isAllowedLength(bod["message"]) {
		rw.WriteHeader(400)
		return
	}

	twClient := ConnTwitter(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	comp, err := SendTweetText(twClient, bod["message"])
	println(comp)
	if err != nil {
		fmt.Fprintln(rw, err)
	}
	fmt.Fprintln(rw, comp)

}

func CreateNewTweetMedia(rw http.ResponseWriter, req *http.Request) {
	twClient := ConnTwitter(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))

	err := req.ParseMultipartForm(200000)
	if err != nil {
		rw.WriteHeader(400)
		return
	}
	formdata := req.MultipartForm
	files := formdata.File["media"]
	message := req.FormValue("message")
	if !isAllowedLength(message) {
		rw.WriteHeader(400)
		return
	}
	fmt.Println(message)
	mediaIDs := make([]int64, 4)
	for i := 0; i < 4; i++ {
		mimeType := files[i].Header.Get("Content-Type")
		if contains(mimeType, allowedMimeTypes) {
			buf := bytes.NewBuffer(nil)
			c, err := files[i].Open()
			if err != nil {
				fmt.Println(err)
			}
			if _, err := io.Copy(buf, c); err != nil {
				fmt.Println(err)
			}
			mediaID, err := UploadMedia(twClient, buf.Bytes(), mimeType)
			if err != nil {
				fmt.Println(err)
			}
			mediaIDs = append(mediaIDs, mediaID)

		}
	}
	SendTweetMedia(twClient, mediaIDs, message)

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
func contains(v string, a []string) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}
	return false
}

func isAllowedLength(text string) bool {
	if len(text) > 257 {
		return false
	}
	return true
}

func ConnTwitter(accessToken string, accessSecret string) *twitt.Client {
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitt.NewClient(httpClient)
<<<<<<< HEAD
	return client
}

func SendTweetText(client *twitt.Client, message string) (comp bool, err error) {
	tweet, resp, err := client.Statuses.Update(message, nil)

=======
	tweet, resp, err := client.Statuses.Update("just setting up tinz", nil)
>>>>>>> 4d5e647fb8a4f31ba8a1b4ab0f06af5fdcbd8a96
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

func UploadMedia(client *twitt.Client, byts []byte, mType string) (int64, error) {
	media, _, err := client.Media.Upload(byts, mType)
	if err != nil {
		return 0, err
	}
	return media.MediaID, nil

}

func SendTweetMedia(client *twitt.Client, mediaIDS []int64, message string) (comp bool, err error) {
	tweet, _, err := client.Statuses.Update(message, &twitt.StatusUpdateParams{MediaIds: mediaIDS})

	if err != nil {
		return false, err
	}
	println(tweet.Text)
	return true, nil
}
