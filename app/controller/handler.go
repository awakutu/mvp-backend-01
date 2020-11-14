package controller

import (
	"CoCreate/app/model"
	"CoCreate/app/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
