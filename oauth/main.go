package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var OauthGoogle = &oauth2.Config{
	ClientID:     "<CLIENT_ID>",
	ClientSecret: "<CLIENT_SECRET>",
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Endpoint:     google.Endpoint,
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
}

func googleLogin(w http.ResponseWriter, req *http.Request) {
	URL, err := url.Parse(OauthGoogle.Endpoint.AuthURL)
	if err != nil {
		log.Println(err)
		return
	}

	parameters := url.Values{}
	parameters.Add("client_id", OauthGoogle.ClientID)
	parameters.Add("scope", strings.Join(OauthGoogle.Scopes, " "))
	parameters.Add("redirect_uri", OauthGoogle.RedirectURL)
	// parameters.Add("response_type", "token")
	parameters.Add("response_type", "code")
	parameters.Add("state", "dynamicrandom")
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}

func googleCallback(w http.ResponseWriter, req *http.Request) {
	// state := req.FormValue("state")
	code := req.FormValue("code")
	if code == "" {
		w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := req.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission.."))
		}
		return
	}

	ctx := context.Background()
	token, err := OauthGoogle.Exchange(ctx, code)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("AccessToken: ", token.AccessToken)
	log.Println("Expiration Time:", token.Expiry.String())

	// token := req.FormValue("access_token")
	w.Write(googleGetUserInfo(token.AccessToken))

}

func googleGetUserInfo(token string) []byte {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token))
	if err != nil {
		log.Println("Get: " + err.Error() + "\n")
		return nil
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	return response
}

func main() {
	http.HandleFunc("/auth/google", googleLogin)
	http.HandleFunc("/auth/google/callback", googleCallback)
	http.ListenAndServe(":8080", nil)
}
