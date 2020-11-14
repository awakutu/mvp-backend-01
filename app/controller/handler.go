package controller

import (
	"CoCreate/app/model"
	"CoCreate/app/utils"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/gomail.v2"

	"github.com/gin-gonic/gin"
)

type Verif struct {
	Email string `json:"email"`
}

func CreateAccount(c *gin.Context) {

	var account model.User
	if err := c.Bind(&account); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	pass, err := utils.HashGenerator(account.Password)
	if err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	account.Password = pass
	flag, err := model.InsertNewAccount(account)
	if flag {
		utils.WrapAPISuccess(c, "success", http.StatusOK)
		return
	} else {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
}

func Login(c *gin.Context) {
	var auth model.Auth
	if err := c.Bind(&auth); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("LOGIN")
	flag, err, token := model.Login(auth)
	if flag {
		utils.WrapAPIData(c, map[string]interface{}{
			"token": token,
		}, http.StatusOK, "success")
	} else {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
	}
}

func GetKategori(c *gin.Context) {
	var ka []model.Kategori
	var u model.User
	if err := c.Bind(&u); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	//log.Println("LOGIN")
	uID := c.Param("id")

	q := model.DB.Where("username=?", uID).Find(&u)
	fmt.Println(q, &uID, u.ID)
	if u.Username == "" {
		c.JSON(http.StatusNotFound, gin.H{"MESSAGE ": http.StatusNotFound, "Result": "tidak ditemukan"})
	}

	res := model.GetKateogi(ka)

	utils.WrapAPIData(c, map[string]interface{}{
		"ID":       u.ID,
		"Username": u.Username,
		"Data":     res,
	}, http.StatusOK, "success")
}

func CreateUserKag(c *gin.Context) {
	var usk model.Detail_category
	//usk.IDU := c.Param("id")

	if err := c.Bind(&usk); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	flag, err := model.UserIKat(usk)
	if flag {
		utils.WrapAPISuccess(c, "success", http.StatusOK)
		return
	} else {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
}

//func main() {
// sender configuration.
const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "CoCreate <mvpkelompok1@gmail.com>"
const CONFIG_AUTH_EMAIL = "mvpkelompok1@gmail.com"
const CONFIG_AUTH_PASSWORD = "14112020mvp"

func verif(c *gin.Context) {

	var v Verif
	if err := c.Bind(&v); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", v.Email)
	//mailer.SetAddressHeader("Cc", "tralalala@gmail.com", "Tra Lala La")
	mailer.SetHeader("Subject", "Test mail")
	mailer.SetBody("text/html", "Verifikasi email")
	//mailer.Attach("./sample.png")

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	//flag, err := model.UserIKat(usk)
	err := dialer.DialAndSend(mailer)
	/*if flag {
		utils.WrapAPISuccess(c, "success", http.StatusOK)
		return
	} else {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}*/

	if err != nil {
		log.Fatal(err.Error())
	} else {
		utils.WrapAPISuccess(c, "success", http.StatusOK)
		log.Println("Mail sent!")
	}

}
