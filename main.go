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
	//	curl localhost:8084/api/login -H 'content-type:application/json' -d '{"username": "farhani", "password":"farhan"}'

	router.GET("/api/pref/:id", middleware.Auth, controller.GetKategori)
	//curl localhost:8084/api/pref/farhani -H 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiIn0.-CmeD9djX3ZzMWQ6kmE_W11Cbk1ZmZCSqtl_bgk_GNU'

	router.POST("/api/prefInsert", middleware.Auth, controller.CreateUserKag)
	//curl localhost:8084/api/prefInsert -H 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiIn0.-CmeD9djX3ZzMWQ6kmE_W11Cbk1ZmZCSqtl_bgk_GNU' -H 'content-type:application/json' -d '{"id_user":1,"id_kategori":1, "jenis_kategori":"Keuangan"}'

	router.GET("/api/profil/:username", middleware.Auth, controller.GetProfil)
	//curl 13.250.111.2:8084/api/profil/farhani -H 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiIn0.-CmeD9djX3ZzMWQ6kmE_W11Cbk1ZmZCSqtl_bgk_GNU'''

	router.POST("/api/profil/:username/update", middleware.Auth, controller.UpdateProfil)
	//curl  localhost:8084/api/profil/farhani/update -H 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiIn0.-CmeD9djX3ZzMWQ6kmE_W11Cbk1ZmZCSqtl_bgk_GNU' -H 'content-type:application/json' -d '{"name":"Farhan", "phone":"082251983584"}'

	//admin
	//------------------------------------------------------------------
	router.POST("/api/admin/register", controller.CreateAdmin)
	//curl localhost:8084/api/admin/register -H 'content-type:application/json' -d '{"username": "admin", "password":"123456"}'
	//output :{"code":200,"status":"success"}

	router.POST("/api/admin/login", controller.LoginAdmin)
	//curl localhost:8084/api/admin/login -H 'content-type:application/json' -d '{"username": "admin", "password":"123456"}'
	//{"code":200,"data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiIn0.-CmeD9djX3ZzMWQ6kmE_W11Cbk1ZmZCSqtl_bgk_GNU"},"status":"success"}

	router.GET("/api/admin/listuser", middleware.Auth, controller.GetListUser)
	//curl  localhost:8084/api/profil/farhani/update -H 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiIn0.-CmeD9djX3ZzMWQ6kmE_W11Cbk1ZmZCSqtl_bgk_GNU' | json_pp

	router.Run(":8084")
}
