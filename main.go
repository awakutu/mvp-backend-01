package main

import (
	"CoCreate/app/controller"
	"CoCreate/app/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/api/register", controller.CreateAccount) //tanpa auth
	router.POST("/api/login", controller.Login)

	router.POST("/api/verifikasi", controller.Verifikasi)

	router.GET("/api/verifikasi/:email", controller.VerifikasiSent)

	/*
		curl localhost:8084/api/login -H 'content-type:application/json' -d '{"username": "farhani", "password":"farhan"}'
	*/
	router.GET("/api/pref/:id", middleware.Auth, controller.GetKategori)
	/*
		curl localhost:8084/api/pref/farhani -H 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiIn0.-CmeD9djX3ZzMWQ6kmE_W11Cbk1ZmZCSqtl_bgk_GNU'
	*/

	router.POST("/api/prefInsert", middleware.Auth, controller.CreateUserKag)

	//curl localhost:8084/api/prefInsert -H 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiIn0.-CmeD9djX3ZzMWQ6kmE_W11Cbk1ZmZCSqtl_bgk_GNU' -H 'content-type:application/json' -d '{"id_user":1,"id_kategori":1, "jenis_kategori":"Keuangan"}'

	router.GET("/api/profil/:username", middleware.Auth, controller.GetProfil)

	router.POST("/api/profil/:username/update", middleware.Auth, controller.UpdateProfil)

	router.Run(":8084")
}
