package controller

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Skarlso/google-oauth-go-sample/structs"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var cred Credentials
var conf *oauth2.Config

// Credentials which stores google ids.
type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

// RandToken generates a random @l length token.
func RandToken(l int) (string, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func init() {
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	if err := json.Unmarshal(file, &cred); err != nil {
		log.Println("unable to marshal data")
		return
	}

	conf = &oauth2.Config{
		ClientID:     cred.Cid,
		ClientSecret: cred.Csecret,
		RedirectURL:  "http://localhost:8084/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}

// IndexHandler handles the location /.
func IndexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"MESSAGE ": http.StatusOK, "Result": "Berhassil"})
	//c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

// AuthHandler handles authentication of a user and initiates a session.
func AuthHandler(c *gin.Context) {
	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	//queryState := c.Request.URL.Query().Get("state")
	if retrievedState != c.Query("state") {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
		return
	}
	tok, err := conf.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	client := conf.Client(oauth2.NoContext, tok)
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	u := structs.User{}
	if err = json.Unmarshal(data, &u); err != nil {
		log.Println(err)
		//c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error marshalling response. Please try agian."})
		c.JSON(http.StatusBadRequest, gin.H{"MESSAGE ": http.StatusBadRequest, "Result": "Bad Request"})
		return
	}
	session.Set("user-id", u.Email)
	err = session.Save()
	if err != nil {
		log.Println(err)
		//		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
		c.JSON(http.StatusBadRequest, gin.H{"MESSAGE ": http.StatusBadRequest, "Result": "Bad Request"})
		return
	}
	//seen := false

	c.JSON(http.StatusOK, gin.H{"email": u.Email})
}

func LoginHandler(ctx *gin.Context) {
	state, err := RandToken(32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error while generating random data."})
		return
	}
	session := sessions.Default(ctx)
	session.Set("state", state)
	session.Save()
	ctx.Writer.Write([]byte("<html><title>Golang Google</title> <body> <a href='" + GetLoginURL(state) + "'><button>Lanjutkan!</button> </a> </body></html>"))
}

func GetLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

// FieldHandler is a rudementary handler for logged in users.
func FieldHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user-id")
	c.JSON(http.StatusOK, gin.H{"user": userID})
}
