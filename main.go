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

	//user
	//------------------------------------------------------------------
	router.POST("/api/register", controller.CreateAccount)
	router.POST("/api/login", controller.Login)
	router.POST("/api/verifikasi", controller.Verifikasi)
	router.GET("/api/verifikasi/:email", controller.VerifikasiSent)
	router.GET("/api/pref/:id", middleware.Auth, controller.GetKategori)
	router.POST("/api/prefInsert", middleware.Auth, controller.CreateUserKag)
	router.GET("/api/profil/:username", middleware.Auth, controller.GetProfil)
	router.POST("/api/profil/:username/update", middleware.Auth, controller.UpdateProfil)

	//admin
	//------------------------------------------------------------------
	router.POST("/api/admin/register", controller.CreateAdmin)
	router.POST("/api/admin/login", controller.LoginAdmin)
	router.GET("/api/admin/listuser", middleware.Auth, controller.GetListUser)
	router.POST("/api/admin/updateuser", middleware.Auth, controller.AccepAdmin)

	//posting
	//------------------------------------------------------------------
	router.GET("/api/Dashboard", middleware.Auth, controller.GetListPost)
	router.POST("/api/Dashboard", middleware.Auth, controller.InserPost)

	//like
	//---------------------------------------------------------------------
	router.POST("/api/likei/:id", middleware.Auth, controller.IncLike)
	router.POST("/api/liked/:id", middleware.Auth, controller.DecLike)

	router.Run(":8084")
}
