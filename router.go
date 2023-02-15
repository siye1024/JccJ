package main

import (
	"dousheng/controller/xhttp"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	//apiRouter.GET("/user/", http.UserInfo)
	apiRouter.POST("/user/register/", xhttp.Register)
	//apiRouter.POST("/user/login/", http.Login)

}
