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
	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true
	cfg.AllowCredentials = true
	cfg.AllowMethods = []string{"GET", "POST"}
	cfg.AllowHeaders = []string{"Authorization", "Origin", "Accept", "X-Requested-With", " Content-Type", "Access-Control-Request-Method", "Access-Control-Request-Headers"}
	router.Use(cors.New(cfg))
	//router.Use(cors.Default())

	token, err := controller.RandToken(64)
	if err != nil {
		log.Fatal("unable to generate random token: ", err)
	}
	store := sessions.NewCookieStore([]byte(token))

	//-----------------------------------------------------

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(sessions.Sessions("cocreate", store))

	router.GET("/auth/google/callback", controller.AuthHandler) //redirect
	//router.GET("/auth/google/callback", controller.LoginHandler) //aws
	router.GET("/google", controller.LoginHandler)

	//router.GET("/auth", controller.AuthHandler)
	//router.GET("/google", controller.LoginHandler) //localhost

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
	router.GET("/api/admin/listuserej", middleware.Auth, controller.GetListUserRej)

	//posting
	//------------------------------------------------------------------
	//router.GET("/api/detailposting/:id", middleware.Auth, controller.GetDetailPost)
	router.GET("/api/dashboard/all", middleware.Auth, controller.GetAllListPost)
	router.POST("/api/dashboard", middleware.Auth, controller.InserPost)

	//router.POST("/api/kategori/:id", middleware.Auth, controller.InsertKat) //masukkan kategori di postingan

	//like
	//---------------------------------------------------------------------
	router.POST("/api/likei/:id", middleware.Auth, controller.IncLike)
	router.POST("/api/liked/:id", middleware.Auth, controller.DecLike)

	//-----
	//router.POST("/api/dislikei/:id", middleware.Auth, controller.DIncLike)
	//router.POST("/api/disliked/:id", middleware.Auth, controller.DDecLike)
	//comment
	//---------------------------------------------------------------------

	router.POST("/api/posting/:id", middleware.Auth, controller.GetListComInPost) //detail posting
	router.POST("/api/comment/:id", middleware.Auth, controller.InsertComment)    //masukkan komentar

	//filter tampilan

	router.GET("/api/dashboard/list/:username", middleware.Auth, controller.Tampilkanlistkategoriuser)
	router.POST("/api/dashboard/sort1/:jenis_kategori", middleware.Auth, controller.FilterTampilJenisKat)
	router.POST("/api/dashboard/sort2/:jenisposting", middleware.Auth, controller.FilterTampilAllwTypost)

	//router.POST("/api/Dashboard/:jenis_kategori_user", middleware.Auth, controller.FilterTampilJenisKat)

	//update dan delete
	router.POST("/api/dashboard/checkid", middleware.Auth, controller.CheckIdPost)
	router.POST("/api/dashboard/update", middleware.Auth, controller.UpdatePosting)
	router.POST("/api/dashboard/delete", middleware.Auth, controller.DeletePosting)

	//router.POST("/api/commentd", middleware.Auth, controller.DeleteComment)
	//router.POST("/api/commentu", middleware.Auth, controller.UpdateComment)

	router.Run(":8084") //port server utama
	//port server sign in google // router.Run(":8085")
}
