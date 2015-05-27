package main

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"html/template"
	"log"
	"net/http"
)

var signedOutTemplate = template.Must(template.New("").Parse(`
<html><body>
Sign in with Google:
<form action="/authorize" method="POST"><input type="submit" value="Signin"/></form>
</body></html>
`))

var userInfoTemplate = template.Must(template.New("").Parse(`
<html><body>
Successfully signed in.</br>
UserInfo:
{{.}}
</body></html>
`))

var code = ""
var token = ""

var conf = &oauth2.Config{
	ClientID:     "651227985508-15l25gtcb20055kssnv3v84t42921q4q.apps.googleusercontent.com",
	ClientSecret: "LsjDS8jeOxfFajvf4InibCwt",
	RedirectURL:  "http://localhost:8080/callback",
	Scopes: []string{
		"profile",
	},
	Endpoint: google.Endpoint,
}

const userInfoURL = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/authorize", handleAuthorize)
	http.HandleFunc("/callback", handleCallback)

	http.ListenAndServe("localhost:8080", nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	signedOutTemplate.Execute(w, nil)
}

func handleAuthorize(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, conf.AuthCodeURL("state"), http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}
	client := conf.Client(oauth2.NoContext, tok)

	resp, _ := client.Get(userInfoURL)
	buf := make([]byte, 1024)
	resp.Body.Read(buf)
	userInfoTemplate.Execute(w, string(buf))
}
