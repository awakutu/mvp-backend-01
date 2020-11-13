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

	router.GET("/api/pref", middleware.Auth, controller.GetKategori)

	//server.AssignHandler("/api/pref", controller.GetK)

	router.Run(":8084")
}
