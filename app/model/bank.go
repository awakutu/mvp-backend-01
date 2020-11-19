package model

import (
	"fmt"
	"time"

	"CoCreate/app/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	layoutDateTime = "2006-01-02 15:04:05"
)

type Goguser struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Name        struct {
		FamilyName string `json:"familyName"`
		GivenName  string `json:"givenName"`
	} `json:"name"`
	Emails []struct {
		Value    string `json:"value"`
		Verified bool   `json:"verified"`
	} `json:"emails"`
	Photos []struct {
		Value string `json:"value"`
	} `json:"photos"`
	Provider string `json:"provider"`
	Raw      string `json:"_raw"`
	JSON     struct {
		Sub           string `json:"sub"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Locale        string `json:"locale"`
	} `json:"_json"`
}

type User struct {
	ID         int    `gorm:"primary_key";auto_increment;not_null json:"-"`
	Username   string `json:"Username"`
	Password   string `json:"password"`
	Nama       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Ttl        string `json:"ttl"`
	Foto       []byte `json:"foto"`
	Status     string `json:"status"`
	Verifikasi string `json:"verifikasi"`
}

type UserTemporary struct {
	ID         int    `gorm:"primary_key";auto_increment;not_null json:"-"`
	Username   string `json:"Username"`
	Password   string `json:"password"`
	Nama       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Ttl        string `json:"ttl"`
	Foto       []byte `json:"foto"`
	Status     string `json:"status"`
	Verifikasi string `json:"verifikasi"`
}

type UserReject struct {
	ID         int    `gorm:"primary_key";auto_increment;not_null json:"-"`
	Username   string `json:"Username"`
	Password   string `json:"password"`
	Nama       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Ttl        string `json:"ttl"`
	Foto       []byte `json:"foto"`
	Status     string `json:"status"`
	Verifikasi string `json:"verifikasi"`
}
type Auth struct {
	Username string `json:"Username"`
	Password string `json:"password"`
}

type Kategori struct {
	ID            int    `gorm:"primary_key";auto_increment;not_null json:"-"`
	Nama_kategori string `json:"jenis_kategori"`
}

type Detail_category struct {
	IDU           int    `json:"id_user"`
	IDK           int    `json:"id_kategori"`
	Nama_kategori string `json:"jenis_kategori"`
}

type Admin struct {
	ID       int    `gorm:"primary_key";auto_increment;not_null json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Posting struct {
	ID        int        `gorm:"column:id";auto_increment;not_null json:"id_post"`
	Title     string     `gorm:"column:title" json:"title"`
	Deskripsi string     `gorm:"column:deskripsi" json:"deskripsi"`
	Like      int        `gorm:"column:like" json:"like"`
	Share_pos string     `gorm:"column:share_pos" json:"share_pos"`
	Tgl_pos   *time.Time `gorm:"column:tgl_pos" json:"tgl_pos"`
	Username  string     `gorm:"column:username" json:"username"`
	Comment   []Comment  `gorm:"ForeignKey:ID_posting;association_foreignKey:id"; json:"isi`
}

type Comment struct {
	ID         int        `gorm:"column:id";auto_increment;not_null json:"id_comment"`
	Isi_co     string     `gorm:"column:isi_co" json:"isi"`
	Tgl_co     *time.Time `gorm:"column:tgl_co" json:"tgl_pos"`
	Username   string     `gorm:"column:username" json:"username"`
	ID_posting int        `gorm:"column:id_posting" json:"id_post"`
	//`gorm:"column:id_post" json:"id_post"`
}

/*type Artikel struct {
	ID       int        `gorm:"column:id" json:"artikel_id"`
	Title    string     `gorm:"column:title" json:"title"`
	Kontent  string     `gorm:"column:kontent" json:"kontent"`
	Komentar []Komentar `gorm:"Foreignkey:Artikel_ID;association_foreignkey:ID;" json:"komentar"`
}

type Komentar struct {
	ID         int    `gorm:"column:id" json:"komentar_id"`
	Artikel_ID string `gorm:"column:artikel_id" json:"artikel_id"`
	Komentar   string `gorm:"column:komentar" json:"komentar"`
}*/

func checkMail(user User) (bool, error) {
	err := DB.Where(&User{Email: user.Email}).First(&user)
	if err.RowsAffected == 1 {
		return false, errors.Errorf("Account sudah terdaftar")
	}
	return true, nil
}

func checkUsername(user User) (bool, error) {
	err := DB.Where(&User{Email: user.Email}).First(&user)
	if err.RowsAffected == 1 {
		return false, errors.Errorf("Account sudah terdaftar")
	}
	return true, nil
}

func Login(auth Auth) (bool, error, string) {
	var account User
	if err := DB.Where(&User{Username: auth.Username}).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, errors.Errorf("Account not found"), ""
		}
	}

	err := utils.HashComparator([]byte(account.Password), []byte(auth.Password))
	if err != nil {
		return false, errors.Errorf("Incorrect Password"), ""
	} else {

		sign := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": auth.Username,
			//"account_number": account.AccountNumber,
		})

		token, err := sign.SignedString([]byte("secret"))
		if err != nil {
			return false, err, ""
		}
		return true, nil, token
	}
}

func InsertNewAccount(account User) (bool, error) {

	if err := DB.Create(&account).Error; err != nil {
		return false, errors.Errorf("invalid prepare statement :%+v\n", err)
	}
	return true, nil
}

func InsertNewAccountTemp(account UserTemporary) (bool, error) {

	if err := DB.Create(&account).Error; err != nil {
		return false, errors.Errorf("invalid prepare statement :%+v\n", err)
	}
	return true, nil
}

//get list kategori
func GetKateogi(kat []Kategori) []Kategori {
	DB.Find(&kat)
	fmt.Println(kat)
	return kat
}

//insert kategori from user
func UserIKat(UIK []Detail_category) (bool, error) {

	if err := DB.Create(&UIK).Error; err != nil {
		return false, errors.Errorf("invalid prepare statement :%+v\n", err)
	}
	return true, nil
}

//login admin
func LoginAdmin(auth Auth) (bool, error, string) {
	var account Admin
	if err := DB.Where(&Admin{Username: auth.Username}).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, errors.Errorf("Account not found"), ""
		}
	}
	err := utils.HashComparator([]byte(account.Password), []byte(auth.Password))
	if err != nil {
		return false, errors.Errorf("Incorrect Password"), ""
	} else {

		sign := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"name": auth.Username,
			//"account_number": account.AccountNumber,
		})

		token, err := sign.SignedString([]byte("secret"))
		if err != nil {
			return false, err, ""
		}
		return true, nil, token
	}
}

//admin melihat list user approve
func GetLUser(ul []UserTemporary) []UserTemporary {
	DB.Find(&ul)
	fmt.Println(ul)
	return ul
}

//admin melihat list user reject
func GetLUserRE(ul []UserReject) []UserReject {
	DB.Find(&ul)
	fmt.Println(ul)
	return ul
}

//akun admin baru
func InsertNewAdmin(account Admin) (bool, error) {
	if err := DB.Create(&account).Error; err != nil {
		return false, errors.Errorf("invalid prepare statement :%+v\n", err)
	}
	return true, nil
}

//posting
func Detailpost(post Posting) (bool, error) {
	//masih kosong
	return true, nil
}

func GetAllPost(post []Posting) []Posting {
	DB.Find(&post)
	fmt.Println(post)
	return post
}

func InsertPost(pos Posting) (bool, error) {
	if err := DB.Create(&pos).Error; err != nil {
		return false, errors.Errorf("invalid prepare statement :%+v\n", err)
	}
	return true, nil
}

func InsertCom(com Comment) (bool, error) {
	if err := DB.Create(&com).Error; err != nil {
		return false, errors.Errorf("invalid prepare statement :%+v\n", err)
	}
	return true, nil
}

func GetAllComPost(com []Comment) []Comment {
	DB.Find(&com)
	fmt.Println(com)
	return com
}
