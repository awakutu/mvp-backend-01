package model

import (
	"fmt"

	"CoCreate/app/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type User struct {
	ID       int    `gorm:"primary_key";auto_increment;not_null json:"-"`
	Username string `json:"Username"`
	Password string `json:"password,omitempty"`
	Nama     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Ttl      string `json:"ttl,omitempty"`
	Foto     string `json:"foto,omitempty"`
	Status   bool
}

type Auth struct {
	Username string `json:"name"`
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

func InsertNewAccount(account User) (bool, error) {

	if err := DB.Create(&account).Error; err != nil {
		return false, errors.Errorf("invalid prepare statement :%+v\n", err)
	}
	return true, nil
}

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

func UserIKat(UIK Detail_category) (bool, error) {

	if err := DB.Create(&UIK).Error; err != nil {
		return false, errors.Errorf("invalid prepare statement :%+v\n", err)
	}
	return true, nil
}

/*
func Uprof(up User) (bool, error) {

	//if err := DB.Where(i).Error; err != nil {
		return false, errors.Errorf("invalid prepare statement :%+v\n", err)
	}
	return true, nil
}*/
