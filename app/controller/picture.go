package controller

import (
	"CoCreate/app/model"
	"CoCreate/app/utils"
	"bufio"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func TerimaUploadJPGFoto(c *gin.Context) {
	var foto model.Foto

	if err := c.Bind(&foto); err != nil {
		utils.WrapAPIError(c, err.Error(), http.StatusBadRequest)
		return
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(foto.Value))
	m, formatString, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()
	fmt.Println(bounds, formatString)

	_, err = os.Stat("image/upload/profile")

	if os.IsNotExist(err) {
		errDir := os.MkdirAll("image/upload/profile", 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}

	parm := c.Param("username")
	//Encode from image format to writer
	pngFilename := "image/upload/profile/" + parm + ".jpg"

	f, err := os.OpenFile(pngFilename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
		return
	}
	model.DB.Model(&model.User{}).Where("username = ?", parm).Update("foto", pngFilename)
	//q := model.DB.Exec("UPDATE users SET ")

	err = png.Encode(f, m)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Png file", pngFilename, "created")
}

func GetProfilJPGtobase64(c *gin.Context) {
	var foto model.Foto
	var usr model.User

	parm := c.Param("username")

	fileName := "image/upload/profile/" + parm + ".jpg"

	model.DB.Where("username = ?", parm).Find(&usr)

	foto.Username = usr.Nama

	//foto =

	imgFile, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer imgFile.Close()

	// create a new buffer base on file size
	fInfo, _ := imgFile.Stat()
	var size = fInfo.Size()
	buf := make([]byte, size)

	// read file content into buffer
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)

	foto.Username = parm
	//
	//foto.Value = usr.Foto

	foto.Value = base64.StdEncoding.EncodeToString(buf)
	//fmt.Println("Base64 string is:", imgBase64Str)
	//return imgBase64Str
	utils.WrapAPIData(c, map[string]interface{}{
		"Data": foto,
	}, http.StatusOK, "success")

}
