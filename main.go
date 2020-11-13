package main

import (
	"CoCreate/app/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/api/register", controller.CreateAccount) //tanpa auth
	router.POST("/api/login", controller.Login)

	router.Run(":80")
}
