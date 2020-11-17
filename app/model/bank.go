package model

import (
	"fmt"
	"time"

	"CoCreate/app/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	ID       int    `gorm:"primary_key";auto_increment;not_null json:"-"`
	Username string `json:"Username"`
	Password string `json:"password"`
	Nama     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Ttl      string `json:"ttl"`
	Foto     []byte `json:"foto"`
	Status   bool   `json:status`
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
	Comment   string     `json:"comment"`
	Like      int        `json:"like"`
	Share_pos string     `json:"share_pos"`
	Tgl_pos   *time.Time `json:"tgl_pos"`
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

//get list kategori
func GetKateogi(kat []Kategori) []Kategori {

	//var account User
	//res := map[string]interface{}{}
	DB.Find(&kat)
	/*if err := DB.Where(&User{Username: auth.Username}).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, errors.Errorf("Account not found"), ""
		}
	}*/

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
func GetLUser(ul []User) []User {
	DB.Find(&ul)
	fmt.Println(ul)
	return ul
}

//akun admin baru
func InsertNewAdmin(account Admin) (bool, error) {
	if err := DB.Create(&account).Error; err != nil {
		return false, errors.Errorf("invalid prepare statement :%+v\n", err)
	}
	//DB.Create(&account)
	return true, nil
}
