package main

import (
	"CoCreate/app/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/register", controller.CreateAccount) //tanpa auth
	//router.POST("/guest/login", controller.Login)

	router.Run(":8084")
}
