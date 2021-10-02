package twit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	twitt "github.com/Bwise1/zuri_bot/go-twitter/twitter"
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

type SendTweetResp struct {
	StatusCode int         `json:"status_code"`
	IsSent     bool        `json:"is_sent"`
	Message    interface{} `json:"message"`
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
	rw.Header().Set("Content-Type", "application/json")
	var bod map[string]string
	err := json.NewDecoder(req.Body).Decode(&bod)
	if err != nil {
		rw.WriteHeader(400)
		respBody := SendTweetResp{StatusCode: 400, IsSent: false, Message: err}
		if err := json.NewEncoder(rw).Encode(respBody); err != nil {
			log.Printf("Error Sending Response %v", err)
		}
		return
	}
	if !isAllowedLength(bod["message"]) {
		respBody := SendTweetResp{StatusCode: 400, IsSent: false, Message: "Error Sending Tweet, length is more that 257 characters"}
		if err := json.NewEncoder(rw).Encode(respBody); err != nil {
			log.Printf("Error Sending Response, %v", err)
		}
		return
	}

	twClient := ConnTwitter(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	comp, err := SendTweetText(twClient, bod["message"])
	if err != nil {
		respBody := SendTweetResp{StatusCode: 419, IsSent: comp, Message: err}
		if err := json.NewEncoder(rw).Encode(respBody); err != nil {
			log.Printf("Error Sending Response %v", err)
		}
		return
	}
	respBody := SendTweetResp{StatusCode: 200, IsSent: comp, Message: "Successfully Sent Tweet"}
	if err := json.NewEncoder(rw).Encode(respBody); err != nil {
		log.Printf("Error Sending Response %v", err)
	}

}

func CreateNewTweetMedia(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	twClient := ConnTwitter(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))

	err := req.ParseMultipartForm(200000)
	if err != nil {
		fmt.Println("error here", err)
		respBody := SendTweetResp{StatusCode: 400, IsSent: false, Message: "Error Processign Form"}
		if err := json.NewEncoder(rw).Encode(respBody); err != nil {
			log.Printf("Error Sending Response %v", err)
			return
		}

	}
	formdata := req.MultipartForm
	files := formdata.File["media"]
	message := req.FormValue("message")
	fmt.Println("Hello", message)
	if !isAllowedLength(message) {
		respBody := SendTweetResp{StatusCode: 400, IsSent: false, Message: "Error Sending Tweet, length is more that 257 characters"}
		if err := json.NewEncoder(rw).Encode(respBody); err != nil {
			log.Printf("Error Sending Response %v", err)
		}
		return
	}

	var mediaIDs []int64
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
	fmt.Println(mediaIDs)
	comp, err := SendTweetMedia(twClient, mediaIDs, message)
	if err != nil {
		respBody := SendTweetResp{StatusCode: 419, IsSent: comp, Message: err}
		if err := json.NewEncoder(rw).Encode(respBody); err != nil {
			log.Printf("Error Sending Response %v", err)
		}
	}
	respBody := SendTweetResp{StatusCode: 200, IsSent: comp, Message: "Successfully Sent Tweet"}
	if err := json.NewEncoder(rw).Encode(respBody); err != nil {
		log.Printf("Error Sending Response %v", err)
	}
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
	return len(text) <= 257
}

func ConnTwitter(accessToken string, accessSecret string) *twitt.Client {
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitt.NewClient(httpClient)
	return client
}

func SendTweetText(client *twitt.Client, message string) (comp bool, err error) {
	tweet, _, err := client.Statuses.Update(message, nil)

	if err != nil {
		fmt.Println(err)
		return false, err
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
	fmt.Println(tweet)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}
