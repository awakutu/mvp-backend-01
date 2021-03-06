package controller

import (
	"CoCreate/app/model"
	"CoCreate/app/utils"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var cred Credentials
var conf *oauth2.Config

var jwtKey = []byte("secret")

type UserGoogle struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
	Hd            string `json:"hd"`
}

// Credentials which stores google ids.
type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
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
	conf = &oauth2.Config{
		ClientID:     "913465578188-hai5duusvj9f2h6fv8do8hp79tkpqi5q.apps.googleusercontent.com",
		ClientSecret: "_hD3IHRGZEd2Bg5ICe4CWm7W",
		RedirectURL:  "http://kelompok1.dtstakelompok1.com/auth/google/callback",
		//RedirectURL: "http://localhost:8084/auth",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}
}

// IndexHandler handles the location /.
func IndexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"MESSAGE ": http.StatusOK, "Result": "Berhassil"})
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
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	//datatoken, _ :=
	ioutil.ReadAll(response.Body)
	//u := structs.User{}
	var u UserGoogle
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
		c.JSON(http.StatusBadRequest, gin.H{"MESSAGE ": http.StatusBadRequest, "Result": "Bad Request"})
		return
	}
	//seen := false

	fmt.Println(&u)
	var account model.User
	var accountT model.UserTemporary
	var accountR model.UserReject

	account.Email = u.Email
	accountT.Email = u.Email
	accountR.Email = u.Email

	account.Nama = u.Name
	accountT.Nama = u.Name
	accountT.Nama = u.Name

	account.Status = "Aktif"

	account.Username = u.Name
	accountT.Username = u.Name
	accountR.Username = u.Name
	//account.Status = "Tidak Aktif"

	account.Password = ""
	accountT.Password = ""
	accountR.Password = ""

	pass, err := utils.HashGenerator(account.Password)
	if err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	account.Password = pass
	accountT.Password = pass
	accountR.Password = pass

	e := model.DB.Where("email=?", u.Email).First(&account)
	if e.RowsAffected == 1 {
		utils.WrapAPIError(c, "Email Sudah Ada", http.StatusOK) //c.JSON(http.StatusOK, gin.H{"status check email": "ok"})
	}

	if err := model.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"eror": err})
	}

	if err := model.DB.Create(&accountT).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"eror": err})
	}

	if err := model.DB.Create(&accountR).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"eror": err})
	}
	model.DB.Where("email= ?", u.Email).Delete(&accountT)

	model.DB.Where("email= ?", u.Email).Delete(&accountR)

	var auth model.Auth

	pass = ""

	auth.Username = u.Name
	auth.Password = pass

	fmt.Println(u.Name)
	fmt.Println(account.Password)

	flag, err, token := model.Login(auth)
	if flag {
		session.Set("token", token)
		session.Save()
		c.Header("Authorization", token)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Authorization", token)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Set("Authorization", token)
		c.Copy().Set("Authorization", token)

		c.Redirect(http.StatusMovedPermanently, "http://localhost:3000/PrefCategory")
		c.Next()

	} else {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
	}

}

//fmt.Println(tokenString)

//

//c.Header("Content-Type", "application/json")

//c.Request.Write("Authorization:", tokenString)

//c.Request.Response.se
//c.Writer.Header().Set("Authorization:", tokenString)

//c.JSON(http.StatusOK, gin.H{"token JWT Generate": tokenString})
//c.JSON(http.StatusOK, gin.H{"token google": datatoken, "token JWT Generate": tokenString})
//c.Redirect(http.StatusTemporaryRedirect, "http://localhost:8084/api/pref/"+u.Name)

//	c.Redirect(http.StatusPermanentRedirect, "http://localhost:8084/api/pref/farhani")

func LoginHandler(ctx *gin.Context) {
	state, err := RandToken(32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error while generating random data."})
		return
	}
	session := sessions.Default(ctx)
	session.Set("state", state)
	session.Save()
	ctx.Writer.Write([]byte("<html><title>Lanjutkan Sign Google</title> <body> <a href='" + GetLoginURL(state) + "'><button>Lanjutkan!</button> </a> </body></html>"))
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
