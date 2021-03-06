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

	//user
	//------------------------------------------------------------------

	router.POST("/api/register", controller.CreateAccount)
	router.POST("/api/login", controller.Login)
	router.POST("/api/loginGoogle", controller.AuthLoginGoogle)
	router.POST("/api/verifikasi", controller.Verifikasi)
	router.GET("/api/verifikasi/:email", controller.VerifikasiSent)
	router.GET("/api/pref/:username", middleware.Auth, controller.GetKategori)
	router.POST("/api/prefInsert", middleware.Auth, controller.CreateUserKag)
	router.GET("/api/profile/:username", middleware.Auth, controller.GetProfil)
	router.POST("/api/profile/:username/update", middleware.Auth, controller.UpdateProfil)

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

	//---------------------------------------------------------------------
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

	//trending
	//---------------------------------------------------------------------
	router.GET("/api/dashboard/trending_artikel", middleware.Auth, controller.TrendingArtikel)
	router.GET("/api/dashboard/trending_membership", middleware.Auth, controller.TrendingMembership)

	//proyekinovasi
	//project
	//---------------------------------------------------------------------
	router.POST("/api/project", middleware.Auth, controller.InsertProject)             //create project
	router.GET("/api/project", middleware.Auth, controller.GetListProjAll)             //dapatkan semua project
	router.GET("/api/project/list/:username", middleware.Auth, controller.GetListProj) //Apinya salah soon updated
	router.POST("/api/project/update", middleware.Auth, controller.UpdateProj)         //edit suatu projek tertentu (id = id_projek)
	router.POST("/api/project/delete", middleware.Auth, controller.DeleteProj)         //delete suatu projek , akan di update besok routnya
	router.POST("/api/project/detail/:id", middleware.Auth, controller.GetProj)        // dapatkan detail projek dg id project tertentu

	router.POST("/api/project/groupinsert/:id", middleware.Auth, controller.InsertGroupProj) //masukkan anggota ke dalam grup projek (id=idprojek)
	router.GET("/api/project/groupanggota/:id", middleware.Auth, controller.GetListAnggota)  //Dapatkan list semua anggota dari satu proyek (id=idprojek)
	router.POST("/api/project/groupdelete", middleware.Auth, controller.DeleteAnggota)       //delete salah satu anggota di suatu proyek (edit soon)

	router.POST("/api/project/task", middleware.Auth, controller.InsertTask)             //buat task
	router.GET("/api/project/task/:id", middleware.Auth, controller.GetTask)             //dapatkan semua task dari satu projek id = id_project
	router.POST("/api/project/task/:id/edit", middleware.Auth, controller.UpdateTask)    //edit task  id = id_task
	router.POST("/api/project/task/:id/delete", middleware.Auth, controller.DeleteTask)  //delete id = id_task
	router.GET("/api/project/detailtask/:id", middleware.Auth, controller.GetDetailTask) //id = id_task

	//comment di task
	router.POST("/api/task/:id", middleware.Auth, controller.CommentTask)
	//router.POST("/api/commentask/:id")

	//upload foto
	//---------------------------------------------------------------------
	router.POST("/api/upload/profil/:username", controller.TerimaUploadJPGFoto)
	router.GET("/api/get/profil/:username", controller.GetProfilJPGtobase64)

	//upload foto postingan
	//---------------------------------------------------------------------
	router.POST("/api/upload/posting/:id", controller.TerimaUploadPostingFoto)
	router.GET("/api/get/posting/:id", controller.GetPostingJPGtobase64)

	//share postingan
	//---------------------------------------------------------------------
	router.POST("/api/sharepostingan/:id")
	router.GET("/api/sharepostingan/:id")

	//gogle
	//---------------------------------------------------------------------
	router.GET("/auth/google/callback", controller.AuthHandler) //redirect
	//router.GET("/google", controller.LoginHandler)
	//router.GET("/auth/google/callback", controller.LoginHandler) //aws
	//router.GET("/auth", controller.AuthHandler)
	router.GET("/google", controller.LoginHandler) //localhost

	//portofolio
	//---------------------------------------------------------------------
	router.GET("/api/portofolio/:username", middleware.Auth, controller.GetPortofolio)

	router.POST("/api/portofolio/insertExp/:username", middleware.Auth, controller.InsertExprience)
	router.POST("/api/portofolio/insertExt/:username", middleware.Auth, controller.InsertExpertise)

	router.POST("/api/portofolio/updateExp/:id", middleware.Auth, controller.UpdateExprience)
	router.POST("/api/portofolio/updateExt/:id", middleware.Auth, controller.UpdateExpertise)

	router.POST("/api/portofolio/deleteExp/:id", middleware.Auth, controller.DeleteExprience)
	router.POST("/api/portofolio/deleteExt/:id", middleware.Auth, controller.DeleteExpertise)

	//---------------------------------------------------------------------
	router.Run(":8084") //port server utama
}
