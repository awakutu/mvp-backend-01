package main

import (
	"CoCreate/app/controller"
	"CoCreate/app/middleware"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	token, err := controller.RandToken(64)
	if err != nil {
		log.Fatal("unable to generate random token: ", err)
	}
	store := sessions.NewCookieStore([]byte(token))

	//-----------------------------------------------------
	//router.
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(sessions.Sessions("cocreate", store))

	router.GET("/auth/google/callback", controller.AuthHandler) //redirect
	//router.GET("/auth/google/callback", controller.LoginHandler) //aws
	router.GET("/google", controller.LoginHandler) //localhost

	//-----------------------------------------------------
	//authorized := router.Group("/battle")
	//authorized.Use(middleware.AuthorizeRequest())
	//{
	//	authorized.GET("/field", controller.FieldHandler)
	//}

	//user
	//------------------------------------------------------------------
	router.POST("/api/register", controller.CreateAccount)
	router.POST("/api/login", controller.Login)
	router.POST("/api/verifikasi", controller.Verifikasi)
	router.GET("/api/verifikasi/:email", controller.VerifikasiSent)
	router.GET("/api/pref/:username", middleware.Auth, controller.GetKategori)
	router.POST("/api/prefInsert", middleware.Auth, controller.CreateUserKag)
	router.GET("/api/profile/:username", middleware.Auth, controller.GetProfil)
	router.POST("/api/profile/:username/update", middleware.Auth, controller.UpdateProfil)
	router.POST("/api/insertGDB")

	//admin
	//------------------------------------------------------------------
	router.POST("/api/admin/register", controller.CreateAdmin)
	router.POST("/api/admin/login", controller.LoginAdmin)
	router.GET("/api/admin/listuser", middleware.Auth, controller.GetListUser)
	router.POST("/api/admin/updateuser", middleware.Auth, controller.AccepAdmin)
	router.POST("/api/admin/reject", middleware.Auth, controller.RejectAd)
	router.POST("/api/admin/rejectoap", middleware.Auth, controller.RejectoApprov)
	//baru
	router.GET("/api/admin/listuserej", middleware.Auth, controller.GetListUserRej)

	//posting
	//------------------------------------------------------------------
	router.GET("/api/Dashboard", middleware.Auth, controller.GetListPost)
	router.GET("/api/Dashboard/all", middleware.Auth, controller.GetAllListPost)
	router.POST("/api/Dashboard", middleware.Auth, controller.InserPost)

	//like
	//---------------------------------------------------------------------
	router.POST("/api/likei/:id", middleware.Auth, controller.IncLike)
	router.POST("/api/liked/:id", middleware.Auth, controller.DecLike)

	//comment
	//---------------------------------------------------------------------

	router.GET("/api/posting/:id", middleware.Auth, controller.GetListComInPost)

	router.POST("/api/comment/:id", middleware.Auth, controller.InsertCo)

	router.Run(":8084") //port server utama
	//port server sign in google // router.Run(":8085")
}
