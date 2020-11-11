package model

import (
	"CoCreate/app/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Account struct {
	ID        int    `gorm:"primary_key";auto_increment;not_null json:"-"`
	IdAccount string `json:"id_account,omitempty"`
	Username  string `json:"Username"`
	Password  string `json:"password,omitempty"`
	//AccountNumber int    `json:"account_number,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone"`
}

type Auth struct {
	Username string `json:"name"`
	Password string `json:"password"`
}

func Login(auth Auth) (bool, error, string) {
	var account Account
	if err := DB.Where(&Account{Username: auth.Username}).First(&account).Error; err != nil {
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

func InsertNewAccount(account Account) (bool, error) {
	//account.AccountNumber = utils.RangeIn(111111, 999999)
	//account.Saldo = 0
	//account.IdAccount = fmt.Sprintf("id-%d", utils.RangeIn(111, 999))
	if err := DB.Create(&account).Error; err != nil {
		return false, errors.Errorf("invalid prepare statement :%+v\n", err)
	}
	return true, nil
}
