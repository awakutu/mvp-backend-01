package controller

import (
	"CoCreate/app/model"
	"CoCreate/app/utils"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
	"gorm.io/gorm/clause"

	"github.com/gin-gonic/gin"
)

type Verif struct {
	Email    string `json:"email"`
	Username string `json:"Username"`
}

//buat akun baru
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

//buat akun sementara
func CreateAccountTEmp(c *gin.Context) {

	var account model.UserTemporary
	if err := c.Bind(&account); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	model.InsertNewAccountTemp(account)
}

//login
func Login(c *gin.Context) {
	var auth model.Auth
	var account model.UserTemporary
	var account1 model.User
	var account2 model.UserReject
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

	//usename direject hubunngi admin
	q2 := model.DB.Where("username=?", auth.Username).First(&account2)
	b2 := q2.RowsAffected
	if b2 == 1 {
		utils.WrapAPIError(c, "Username tidak diterima hubungi admin di mvpkelompok1@gmail.com", http.StatusOK)
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

//tampilkan list kategori
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

//create kategori pilihan user
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

//kirim email
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

//kirim veirifkasi
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

	err1 := model.DB.Model(&u).Where("email= ?", uID).Update("verifikasi", "Ya")
	if err1 != nil {
		/*utils.WrapAPIData(c, map[string]interface{}{
			"Email": uID,
		}, http.StatusOK, "success")*/
		c.Redirect(http.StatusPermanentRedirect, "http://localhost:3000/Verified")
		return
	} else {
		utils.WrapAPIError(c, "err1.Error()", http.StatusBadRequest)
		return
	}
}

//tmapilkan profil
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

//update profil
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

//tampilkan user yang di akan diapprove
func GetListUser(c *gin.Context) {
	var usr []model.UserTemporary
	res := model.GetLUser(usr)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": res,
	}, http.StatusOK, "success")
}

//tampilkan user yg diblacklist
func GetListUserRej(c *gin.Context) {
	var usr []model.UserReject
	res := model.GetLUserRE(usr)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": res,
	}, http.StatusOK, "success")
}

//buat admin
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

//login admin
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

//hapus dari tabel approve >> bisa login
func AccepAdmin(c *gin.Context) {
	var accountemp []model.UserTemporary

	if err := c.Bind(&accountemp); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	//fmt.Println(account)

	//q := model.DB.Save(&account)

	//approve akun
	q := model.DB.Delete(&accountemp)
	row := q.RowsAffected

	utils.WrapAPIData(c, map[string]interface{}{
		"Data":         accountemp,
		"Row affected": row,
	}, http.StatusOK, "success")

}

//masuk ke tabel blacklist
func RejectAd(c *gin.Context) {
	var accountemp []model.UserTemporary
	//var account2 []model.UserReject

	if err := c.Bind(&accountemp); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	//reject akkun
	q2 := model.DB.Model(model.UserReject{}).Create(&accountemp)
	row := q2.RowsAffected

	//hapus dari tabel sementara
	//q :=
	model.DB.Delete(&accountemp)
	//b := q.RowsAffected

	utils.WrapAPIData(c, map[string]interface{}{
		"Data":         accountemp,
		"Row affected": row,
	}, http.StatusOK, "success")
}

//reject ke approv
func RejectoApprov(c *gin.Context) {
	var account2 []model.UserReject

	if err := c.Bind(&account2); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	//reject akkun
	//q2 := model.DB.Model(model.UserReject{}).Create(&account1)
	//br := q2.RowsAffected

	//hapus dari tabel reject
	q := model.DB.Delete(&account2)
	b := q.RowsAffected

	utils.WrapAPIData(c, map[string]interface{}{
		"Data":         account2,
		"Row affected": b,
	}, http.StatusOK, "success")
}

//get all postingan
func GetAllListPost(c *gin.Context) {
	var Post []model.Posting

	model.DB.Preload(clause.Associations).Find(&Post)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": Post,
	}, http.StatusOK, "success")
}

type Viewrespon struct {
	Sumview int `json:"view"`
}

//detail posting
func GetDetailPost(c *gin.Context) {
	var Post model.Posting

	parm := c.Param(":id")

	var viewrespon Viewrespon

	model.DB.Where("id=? ", parm).Find(&Post)

	model.DB.Model(&Post).Where("id= ?", parm).Update("view", Post.View+1).Scan(&viewrespon)

	model.DB.Preload("Comment").Find(&Post)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": Post,
		"View": viewrespon,
	}, http.StatusOK, "success")
}

//insert posting
func InserPost(c *gin.Context) {
	var Post model.Posting
	//ar account model.Admin
	if err := c.Bind(&Post); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	//secs := now.UnixNano()

	Post.Tgl_pos = &now

	flag, err := model.InsertPost(Post)
	if flag {
		utils.WrapAPIData(c, map[string]interface{}{
			//"Data":  flag,
			"Data": Post,
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

//increase like
func IncLike(c *gin.Context) {
	var Post model.Posting

	if err := c.Bind(&Post); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	sk := c.Param("id")

	//	var S1 SumLike
	var Sumlikerespon SumLikeRes
	model.DB.Where("id=? ", sk).Find(&Post)

	model.DB.Model(&Post).Where("id= ?", sk).Update("like", Post.Like+1).Scan(&Sumlikerespon)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": Sumlikerespon,
	}, http.StatusOK, "success")
}

//decrease like
func DecLike(c *gin.Context) {
	var Post model.Posting

	if err := c.Bind(&Post); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	sk := c.Param("id")

	var Sumlikerespon SumLikeRes
	model.DB.Where("id=? ", sk).Find(&Post)
	nilaiLike := float64(Post.Like)

	res := math.Signbit(nilaiLike - 1)
	if res == false {
		model.DB.Model(&Post).Where("id= ?", sk).Update("like", Post.Like-1).Scan(&Sumlikerespon)
		utils.WrapAPIData(c, map[string]interface{}{
			"Data": Sumlikerespon,
		}, http.StatusOK, "success")
	} else {
		//TSumlikerespon := Sumlikerespon
		utils.WrapAPIData(c, map[string]interface{}{
			"Data": Sumlikerespon,
		}, http.StatusOK, "success")
	}
}

//funct insert comment
func InsertComment(c *gin.Context) {
	var co model.Comment
	//var account model.User

	if err := c.Bind(&co); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	var parm string
	parm = c.Param("id") //idpostting

	now := time.Now()
	co.Tgl_co = &now

	parm2, err := strconv.Atoi(parm)
	//fmt.Println(parm2, err)

	co.ID_posting = parm2

	flag, err := model.InsertCommentm(co)
	if flag {
		utils.WrapAPIData(c, map[string]interface{}{
			"ID Postingan": co.ID_posting,
		}, http.StatusOK, "success")
		return
	} else {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
}

func UpdateComment(c *gin.Context) {
	var co model.Comment
	//var account model.User

	if err := c.Bind(&co); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	//var parm string
	//parm = c.Param("id") //idpostting

	now := time.Now()
	co.Tgl_co = &now

	//parm2, err := strconv.Atoi(parm)
	//fmt.Println(parm2, err)

	//co.ID_posting = parm2

	model.DB.Save(&co)
	//if flag {
	utils.WrapAPIData(c, map[string]interface{}{
		"Data": co,
	}, http.StatusOK, "success")
}

/*func updateUser(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		title := c.Param("title")
		deskripsi := c.Param("deksripsi")

		var posting model.Posting
		db.Where("name=?", name).Find(&user)
		user.Email = email
		db.Save(&user)
		return c.String(http.StatusOK, name+" user successfully updated")
	}
}*/

func CheckIdPost(c *gin.Context) {
	var post model.Posting
	var PostTemporary model.Posting
	var PostNonTemp model.Posting

	if err := c.Bind(&post); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	model.DB.Where("title = ? and deskripsi = ?", post.Title, post.Deskripsi).Find(&PostTemporary)
	fmt.Println(PostNonTemp.ID)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": PostTemporary.ID,
	}, http.StatusOK, "success")
}

func UpdatePosting(c *gin.Context) {
	var Post model.Posting
	//ar account model.Admin
	if err := c.Bind(&Post); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()

	Post.Tgl_pos = &now

	model.DB.Model(&Post).UpdateColumns(&Post)

	utils.WrapAPIData(c, map[string]interface{}{
		//"Data":  flag,
		"Data": Post,
	}, http.StatusOK, "success")
}

//
func DeletePosting(c *gin.Context) {
	var Post model.Posting
	//ar account model.Admin
	if err := c.Bind(&Post); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	//secs := now.UnixNano()

	Post.Tgl_pos = &now

	model.DB.Delete(&Post)

	utils.WrapAPIData(c, map[string]interface{}{
		//"Data":  flag,
		"Data": Post,
	}, http.StatusOK, "success")
}

func DeleteComment(c *gin.Context) {
	var co model.Comment
	//var account model.User

	if err := c.Bind(&co); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	var parm string
	parm = c.Param("id") //idpostting

	now := time.Now()
	co.Tgl_co = &now

	parm2, err := strconv.Atoi(parm)
	fmt.Println(parm2, err)

	co.ID_posting = parm2

	model.DB.Delete(&co)
	//if flag {
	utils.WrapAPIData(c, map[string]interface{}{
		"Data": co,
	}, http.StatusOK, "success")
	//	return
	//} else {
	//	utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
	//	return
	//}
}

//get list commen dalam post
func GetListComInPost(c *gin.Context) {
	var Post model.Posting
	var usr model.User

	if err := c.Bind(&usr); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	parm := c.Param("id") //idposting

	var viewrespon Viewrespon

	model.DB.Where("id=? ", parm).Find(&Post)

	model.DB.Model(&Post).Where("id= ?", parm).Update("view", Post.View+1).Scan(&viewrespon)

	model.DB.Preload(clause.Associations).Preload("Comment", "id_posting", parm).Find(&Post) //buat mencari posting di dalam  comment dimana where id positng adalah paramater

	var v bool
	if usr.Username == Post.Username {
		v = true
	} else {
		v = false
	}

	utils.WrapAPIData(c, map[string]interface{}{
		"Posting":  Post,
		"View":     viewrespon,
		"Pemilik ": v,
	}, http.StatusOK, "success")
}

func InsertKat(c *gin.Context) {
	var kateg model.Kategores
	//var account model.User

	if err := c.Bind(&kateg); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	var parm string
	parm = c.Param("id") //idpostting

	param_tempry, err := strconv.Atoi(parm)
	fmt.Println(param_tempry, err)

	model.DB.Create(kateg)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data":         kateg,
		"ID Postingan": parm,
	}, http.StatusOK, "success")
}

/*
//increase like
func DIncLike(c *gin.Context) {
	var li model.Posting

	if err := c.Bind(&li); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	sk := c.Param("id")

	//	var S1 SumLike
	var S2 SumDisLikeRes
	model.DB.Where("id=? ", sk).Find(&li)

	model.DB.Model(&li).Where("id= ?", sk).Update("dislike", li.Dislike+1).Scan(&S2)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": S2,
	}, http.StatusOK, "success")
}

type SumDisLike struct {
	Like int `json:"dislike"`
}

type SumDisLikeRes struct {
	Like int `json:"dislike"`
}

//decrease like
func DDecLike(c *gin.Context) {
	var li model.Posting

	if err := c.Bind(&li); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	sk := c.Param("id")

	//var S1 SumLike
	var S2 SumDisLikeRes
	model.DB.Where("id=? ", sk).Find(&li)
	model.DB.Model(&li).Where("id= ?", sk).Update("dislike", li.Dislike-1).Scan(&S2)
	utils.WrapAPIData(c, map[string]interface{}{
		"Data": S2,
	}, http.StatusOK, "success")
}*/

func FilterTampilJenisKat(c *gin.Context) {
	var Post []model.Posting

	sk := c.Param("jenis_kategori")

	model.DB.Preload(clause.Associations).Where("kategorip = ?", sk).Find(&Post)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": Post,
	}, http.StatusOK, "success")
}

func FilterTampilAllwTypost(c *gin.Context) {
	var Post []model.Posting

	parm := c.Param("jenisposting")

	q := model.DB.Preload(clause.Associations).Where("typep = ?", parm).Find(&Post)
	fmt.Println(q)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": Post,
	}, http.StatusOK, "success")
}

func FilterTampilAllwKatUser(c *gin.Context) {
	var Post []model.Posting
	var usr model.User

	if err := c.Bind(&usr.Username); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	model.DB.Where("kategori = ?", usr).Preload("Comment").Find(&Post)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": Post,
	}, http.StatusOK, "success")
}

func Tampilkanlistkategoriuser(c *gin.Context) {
	var kategori []model.Detail_category
	//var
	var usr model.User

	var parm string
	parm = c.Param("username") //idpostting
	param_tempry, err := strconv.Atoi(parm)
	fmt.Println(err)

	model.DB.Where("username = ?", param_tempry).Find(&usr)
	model.DB.Where("id_u= ?", usr.ID).Find(&kategori)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": kategori,
	}, http.StatusOK, "success")
}

type Trending_membership struct {
	Username string `json:"username"`
	Trending int    `json:"trending_membership"`
}

func TrendingMembership(c *gin.Context) {
	var trending_membership []Trending_membership

	model.DB.Raw("Select username, count(username) as trending from postings group by username order by trending desc limit 3").Scan(&trending_membership)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": trending_membership,
	}, http.StatusOK, "success")
}

type Trending_artikel struct {
	Title    string `json:"title"`
	Trending int    `json:"trending_artikel"`
}

func TrendingArtikel(c *gin.Context) {
	var trending_artikel []Trending_artikel

	model.DB.Raw("Select title, max(view) as trending from postings group by title order by trending desc limit 3").Scan(&trending_artikel)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": trending_artikel,
	}, http.StatusOK, "success")
}

func InsertProject(c *gin.Context) {
	var pro model.Project
	var pro1 model.Project
	var grup1 model.GrupProject

	if err := c.Bind(&pro); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	pro.Tgl_pos = &now
	pro.Tgl_edit = &now

	pro.SumAnggota = 1
	model.InsertProj(pro)

	model.DB.Model(&pro).Where("title=?", pro.Title).Scan(&pro1)

	grup1.IDP = pro1.ID

	fmt.Println(pro1.ID) //lihat nilai pro id  //nah save baru run mba //oke  //sdh bisa????? //SUUUDAHHHHHH  //oke  //ku out ya //okaayyy

	grup1.Role = "admin"
	grup1.Date_join = &now
	grup1.IDU = pro.IDU
	grup1.Username = pro.Username
	grup1.Projectname = pro.Title

	model.InsertGrupProj(grup1)

	utils.WrapAPIData(c, map[string]interface{}{
		"role":                      grup1.Role,
		"Data Insert Project":       pro,
		"Data Insert Group Project": grup1,
	}, http.StatusOK, "success")

}

type SumGroup struct {
	SumAnggota int `json:"sum_anggota"`
}

func InsertGroupProj(c *gin.Context) {
	var grup model.GrupProject
	var prj model.Project
	//ar account model.Admin
	if err := c.Bind(&grup); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	sk := c.Param("id")

	var sg SumGroup
	model.DB.Where("id=? ", sk).Find(&prj)
	model.DB.Model(&prj).Where("id= ?", sk).Update("sum_anggota", prj.SumAnggota+1).Scan(&sg)

	now := time.Now()
	//secs := now.UnixNano()
	grup.Date_join = &now

	flag, err := model.InsertGrupProj(grup)
	if flag {
		utils.WrapAPIData(c, map[string]interface{}{
			"Data": grup,
		}, http.StatusOK, "success")
		return
	} else {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
}

//list semua projek + komen
func GetListProjAll(c *gin.Context) {
	/* var listproj []model.Project

	res := model.GetProj(listproj)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": res,
	}, http.StatusOK, "success") */

	var proj1 []model.Project

	model.DB.Preload(clause.Associations).Find(&proj1)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": proj1,
	}, http.StatusOK, "success")

}

//list projek berdasarkan username
func GetListProj(c *gin.Context) {
	//var gp []model.GrupProject
	//var gp1 model.GrupProject
	var gp3 []model.Project

	/*if err := c.Bind(&u); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	//log.Println("LOGIN") */

	//di save >> save

	aID := c.Param("username")

	model.DB.Preload("GrupProject", "username", aID).Find(&gp3)

	//model.DB.Raw("Select * from projct group by username order by trending desc limit 3").Scan(&trending_membership)

	//model.DB.Where("grup_projects.username = ?", aID).Preload("GrupProject", "username", aID).Find(&gp3)

	//model.DB.Model(&gp1).Where("username=?", aID).Scan(&gp)
	//model.DB.Model(&gp).Where("id_project=?", gp.id_project).Scan(&gp3)

	utils.WrapAPIData(c, map[string]interface{}{
		"Anggota": gp3,
	}, http.StatusOK, "success")
}

func GetProj(c *gin.Context) {
	//var ka []model.Kategori
	var p model.Project

	pID := c.Param("id")

	model.DB.Where("id=?", pID).Preload(clause.Associations).Find(&p)

	fmt.Println(&pID, p.ID)
	model.DB.Find(&p)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": p,
	}, http.StatusOK, "success")
}

func GetListAnggota(c *gin.Context) {
	var gp []model.GrupProject
	var gp1 model.GrupProject

	/*if err := c.Bind(&u); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}
	//log.Println("LOGIN") */

	aID := c.Param("id")

	model.DB.Model(&gp1).Where("id_p=?", aID).Scan(&gp)

	utils.WrapAPIData(c, map[string]interface{}{
		"Anggota": gp,
	}, http.StatusOK, "success")
}

func UpdateProj(c *gin.Context) {
	var uproj model.Project

	if err := c.Bind(&uproj); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	sk := c.Param("id")

	var sk1 model.Project
	model.DB.Where("id=?", sk).Find(&sk1)

	fmt.Println(sk1.ID)

	now := time.Now()
	uproj.Tgl_edit = &now

	result := model.DB.Model(model.Project{}).Where("id = ?", sk1.ID).Updates(uproj)
	b := result.RowsAffected

	utils.WrapAPIData(c, map[string]interface{}{
		"Data":        &uproj,
		"Rows_update": b,
	}, http.StatusOK, "success")
}

func DeleteProj(c *gin.Context) {
	var membergrup []model.GrupProject
	var pro model.Project

	if err := c.Bind(&pro); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	q := model.DB.Where("id=?", pro.ID).Find(&pro)
	y := model.DB.Where("id_p=? ", pro.ID).Find(&membergrup)
	model.DB.Delete(&pro)
	model.DB.Delete(&membergrup) //slh
	fmt.Println(q)
	fmt.Println(y)

	utils.WrapAPIData(c, map[string]interface{}{
		"Project":    pro,
		"Membergrup": membergrup,
	}, http.StatusOK, "success")
}

func DeleteAnggota(c *gin.Context) {
	var grup model.GrupProject
	var pro model.Project
	var sg SumGroup

	if err := c.Bind(&grup); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	q := model.DB.Where("id_p=? AND username=?", grup.IDP, grup.Username).Find(&grup)
	y := model.DB.Where("id=? ", grup.IDP).Find(&pro)
	model.DB.Model(&pro).Where("id= ?", grup.IDP).Update("sum_anggota", pro.SumAnggota-1).Scan(&sg)
	model.DB.Delete(&grup) //slh
	fmt.Println(q)
	fmt.Println(y)

	utils.WrapAPIData(c, map[string]interface{}{
		"Delete": grup,
		"tes":    sg,
	}, http.StatusOK, "success")
}

func InsertTask(c *gin.Context) {
	var task1 model.Task

	if err := c.Bind(&task1); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	task1.Tgl_pos = &now
	task1.Tgl_edit = &now

	model.InsertTask(task1)

	utils.WrapAPIData(c, map[string]interface{}{
		"task": task1,
	}, http.StatusOK, "success")

}

func GetTask(c *gin.Context) {
	//var ka []model.Kategori
	var t []model.Task
	var t1 model.Task

	tID := c.Param("id")

	model.DB.Model(&t1).Where("id_p=?", tID).Scan(&t)

	fmt.Println(&tID, t1.ID)

	utils.WrapAPIData(c, map[string]interface{}{
		"Data": t,
	}, http.StatusOK, "success")
}

func UpdateTask(c *gin.Context) {
	var utask model.Task

	if err := c.Bind(&utask); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	tk := c.Param("id")

	var tk1 model.Task
	model.DB.Where("id=?", tk).Find(&tk1)

	fmt.Println(utask.ID)

	now := time.Now()
	utask.Tgl_edit = &now

	result := model.DB.Model(model.Task{}).Where("id= ?", tk1.ID).Updates(utask)
	b := result.RowsAffected

	utils.WrapAPIData(c, map[string]interface{}{
		"Data":        &utask,
		"Rows_update": b,
	}, http.StatusOK, "success")
}

func DeleteTask(c *gin.Context) {
	var dtask model.Task

	if err := c.Bind(&dtask); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	dt := c.Param("id")

	y := model.DB.Where("id=? ", dt).Find(&dtask)
	model.DB.Delete(&dtask)
	fmt.Println(y)

	utils.WrapAPIData(c, map[string]interface{}{
		"Delete": dtask,
	}, http.StatusOK, "success")
}
