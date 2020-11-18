package controller

import (
	"CoCreate/app/model"
	"CoCreate/app/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/gin-gonic/gin"
)

type Verif struct {
	Email    string `json:"email"`
	Username string `json:"Username"`
}

func CreateAccount(c *gin.Context) {
	var at model.User
	var account model.User
	var accountT model.UserTemporary
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

	q := model.DB.Where("username=?", account.Username).First(&account)

	b := q.RowsAffected
	if b == 1 {
		utils.WrapAPIError(c, "Username Sudah ada", http.StatusOK)
		return
	}

	q2 := model.DB.Where("email=?", account.Email).First(&account)

	b2 := q2.RowsAffected
	if b2 == 1 {
		utils.WrapAPIError(c, "Email Sudah Ada", http.StatusOK)
		return
	}

	flag, err := model.InsertNewAccount(account)

	accountT.ID = account.ID
	accountT.Email = account.Email
	accountT.Nama = account.Nama
	accountT.Password = account.Password
	accountT.Phone = account.Phone
	accountT.Status = account.Status
	accountT.Ttl = account.Ttl
	accountT.Username = account.Username
	model.InsertNewAccountTemp(accountT)

	model.DB.Where("email=?", account.Email).First(&at)

	if flag {
		utils.WrapAPIData(c, map[string]interface{}{
			"Account ID": at.ID,
			"Username":   account.Username,
		}, http.StatusOK, "success")
		return
	} else {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
}

func CreateAccountTEmp(c *gin.Context) {

	var account model.UserTemporary
	if err := c.Bind(&account); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	model.InsertNewAccountTemp(account)
}

func Login(c *gin.Context) {
	var auth model.Auth
	var account model.UserTemporary
	var account1 model.User
	if err := c.Bind(&auth); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("LOGIN")

	q := model.DB.Where("username=?", auth.Username).First(&account)
	model.DB.Where("username=?", auth.Username).First(&account1)

	b := q.RowsAffected
	if b == 1 {
		utils.WrapAPIError(c, "Username Belum diapprove", http.StatusOK)
		return
	}

	flag, err, token := model.Login(auth)
	if flag {
		utils.WrapAPIData(c, map[string]interface{}{
			"token":    token,
			"username": auth.Username,
			"ID":       account1.ID,
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
	uID := c.Param("username")

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
	var usk []model.Detail_category
	//usk.IDU := c.Param("id")

	if err := c.Bind(&usk); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	flag, err := model.UserIKat(usk)
	if flag {
		utils.WrapAPIData(c, map[string]interface{}{
			"Data": usk,
		}, http.StatusOK, "Success")
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

func Verifikasi(c *gin.Context) {

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
	mailer.SetBody("text/html", "Verifikasi email : http://3.15.137.94:8084/api/verifikasi/"+v.Email)
	//mailer.Attach("./sample.png")

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	//flag, err := model.UserIKat(usk)
	err := dialer.DialAndSend(mailer)

	if err != nil {
		log.Fatal(err.Error())
	} else {
		utils.WrapAPIData(c, map[string]interface{}{
			"Email": v.Email,
		}, http.StatusOK, "success")
		log.Println("Mail sent!")
	}
}

func VerifikasiSent(c *gin.Context) {

	var v Verif
	var u model.User
	if err := c.Bind(&v); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	uID := c.Param("email")

	q := model.DB.Where("email=?", uID).Find(&v)
	fmt.Println(q, &uID, v.Email)
	/*if &uID == "" {
		c.JSON(http.StatusNotFound, gin.H{"MESSAGE ": http.StatusNotFound, "Result": "Tidak ada email tersebut"})
	}*/

	err1 := model.DB.Model(&u).Where("email= ?", uID).Update("verifikasi", "aktif")
	if err1 != nil {
		utils.WrapAPIData(c, map[string]interface{}{
			"Email": uID,
		}, http.StatusOK, "success")
		return
	} else {
		utils.WrapAPIError(c, "err1.Error()", http.StatusBadRequest)
		return
	}
}

func GetProfil(c *gin.Context) {
	//var ka []model.Kategori
	var u model.User

	uID := c.Param("username")

	q := model.DB.Where("username=?", uID).Find(&u)
	fmt.Println(q, &uID, u.ID)
	model.DB.Find(&u)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": u,
	}, http.StatusOK, "success")
}

func UpdateProfil(c *gin.Context) {
	var usk model.User

	if err := c.Bind(&usk); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	sk := c.Param("username")

	var sk1 model.User
	model.DB.Where("username=?", sk).Find(&sk1)

	fmt.Println(sk1.ID)

	pass, err := utils.HashGenerator(usk.Password)
	if err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	usk.Password = pass

	result := model.DB.Model(model.User{}).Where("id = ?", sk1.ID).Updates(usk)
	b := result.RowsAffected

	utils.WrapAPIData(c, map[string]interface{}{
		"Data":        &usk,
		"Rows_update": b,
	}, http.StatusOK, "success")
}

func GetListUser(c *gin.Context) {
	var usr []model.UserTemporary
	res := model.GetLUser(usr)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": res,
	}, http.StatusOK, "success")
}

func CreateAdmin(c *gin.Context) {

	var account model.Admin
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
	flag, err := model.InsertNewAdmin(account)
	if flag {
		utils.WrapAPIData(c, map[string]interface{}{
			"Data": account,
		}, http.StatusOK, "success")
		return
	} else {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
}

func LoginAdmin(c *gin.Context) {
	var auth model.Auth
	if err := c.Bind(&auth); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("LOGIN")
	flag, err, token := model.LoginAdmin(auth)
	if flag {
		utils.WrapAPIData(c, map[string]interface{}{
			"token":    token,
			"username": auth.Username,
		}, http.StatusOK, "success")
	} else {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
	}
}

func AccepAdmin(c *gin.Context) {
	var account1 []model.UserTemporary

	if err := c.Bind(&account1); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	//fmt.Println(account)

	//q := model.DB.Save(&account)
	q := model.DB.Delete(&account1)
	b := q.RowsAffected

	utils.WrapAPIData(c, map[string]interface{}{
		"Data":         account1,
		"Row affected": b,
	}, http.StatusOK, "success")

}

//POST
func GetListPost(c *gin.Context) {
	var usr []model.Posting

	res := model.GetAllPost(usr)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": res,
	}, http.StatusOK, "success")
}

func InserPost(c *gin.Context) {
	var usr model.Posting
	//ar account model.Admin
	if err := c.Bind(&usr); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	//secs := now.UnixNano()

	usr.Tgl_pos = &now

	flag, err := model.InsertPost(usr)
	if flag {
		utils.WrapAPIData(c, map[string]interface{}{
			"Data":  flag,
			"Data2": usr,
		}, http.StatusOK, "success")
		return
	} else {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
}

//like
type SumLike struct {
	Like int `json:"like"`
}

type SumLikeRes struct {
	Like int `json:"like"`
}

func IncLike(c *gin.Context) {
	var li model.Posting

	if err := c.Bind(&li); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	sk := c.Param("id")

	//	var S1 SumLike
	var S2 SumLikeRes
	model.DB.Where("id=? ", sk).Find(&li)

	model.DB.Model(&li).Where("id= ?", sk).Update("like", li.Like+1).Scan(&S2)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": S2,
	}, http.StatusOK, "success")
}

func DecLike(c *gin.Context) {
	var li model.Posting

	if err := c.Bind(&li); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	sk := c.Param("id")

	//var S1 SumLike
	var S2 SumLikeRes
	model.DB.Where("id=? ", sk).Find(&li)

	model.DB.Model(&li).Where("id= ?", sk).Update("like", li.Like-1).Scan(&S2)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": S2,
	}, http.StatusOK, "success")
}

func InsertCo(c *gin.Context) {
	var co model.Comment
	//var account model.User

	if err := c.Bind(&co); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	var idp string
	idp = c.Param("id") //idpostting

	now := time.Now()
	co.Tgl_co = &now

	idp2, err := strconv.Atoi(idp)
	fmt.Println(idp2, err)

	co.ID_posting = idp2

	flag, err := model.InsertCom(co)
	if flag {
		utils.WrapAPIData(c, map[string]interface{}{
			"Data":         flag,
			"ID Postingan": co.ID,
		}, http.StatusOK, "success")
		return
	} else {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
}

func GetListComInPost(c *gin.Context) {
	var co []model.Comment
	var co1 model.Comment
	var p model.Posting
	var p1 model.Posting

	//var idp string
	idp := c.Param("id") //idposting

	model.DB.Model(&p).Where("id=?", idp).Scan(&p1)
	model.DB.Model(&co1).Where("id_posting=?", idp).Scan(&co)

	utils.WrapAPIData(c, map[string]interface{}{
		"Posting": p1,
		"Comment": co,
	}, http.StatusOK, "success")
}
