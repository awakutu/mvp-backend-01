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
	ID       int    `gorm:"primary_key";auto_increment;not_null json:"-"`
	Username string `json:"Username"`
	Password string `json:"password"`
	Nama     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Ttl      string `json:"ttl"`
	Foto     []byte `json:"foto"`
	Status   string `json:"status"`
}

type UserTemporary struct {
	ID       int    `gorm:"primary_key";auto_increment;not_null json:"-"`
	Username string `json:"Username"`
	Password string `json:"password"`
	Nama     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Ttl      string `json:"ttl"`
	Foto     []byte `json:"foto"`
	Status   string `json:"status"`
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
	ID        int        `gorm:"primary_key";auto_increment;not_null json:"-"`
	Title     string     `json:"title"`
	Deskripsi string     `json:"deskripsi"`
	Like      int        `json:"like"`
	Share_pos string     `json:"share_pos"`
	Tgl_pos   *time.Time `json:"tgl_pos"`
	Username  string     `json:"username"`
}

type Comment struct {
	ID         int        `gorm:"primary_key";auto_increment;not_null json:"-"`
	Isi_co     string     `json:"title"`
	Deskripsi  string     `json:"deskripsi"`
	Tgl_co     *time.Time `json:"tgl_pos"`
	Username   string     `json:"username"`
	ID_posting int        `json:"id_post"`
}

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

//admin melihat list user
func GetLUser(ul []UserTemporary) []UserTemporary {
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
