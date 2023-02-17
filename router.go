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
	apiRouter.GET("/user/", xhttp.GetUserById)
	apiRouter.POST("/user/register/", xhttp.Register)
	apiRouter.POST("/user/login/", xhttp.Login)

	apiRouter.GET("/feed/", xhttp.Feed)
	//apiRouter.POST("/publish/action/", xhttp.Publish)
	//apiRouter.GET("/publish/list/", xhttp.PublishList)

	// extra apis - I
	//apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	//apiRouter.GET("/favorite/list/", controller.FavoriteList)
	//apiRouter.POST("/comment/action/", controller.CommentAction)
	//apiRouter.GET("/comment/list/", controller.CommentList)

}
