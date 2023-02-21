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
	apiRouter.POST("/publish/action/", xhttp.PublishAction)
	apiRouter.GET("/publish/list/", xhttp.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", xhttp.FavoriteAction)
	apiRouter.GET("/favorite/list/", xhttp.FavoriteList)
	apiRouter.POST("/comment/action/", xhttp.CommentAction)
	apiRouter.GET("/comment/list/", xhttp.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", xhttp.RelationAction)
	apiRouter.GET("/relation/follow/list/", xhttp.FollowList)
	apiRouter.GET("/relation/follower/list/", xhttp.FollowerList)
	apiRouter.GET("/relation/friend/list/", xhttp.FriendList)
	//apiRouter.GET("/message/chat/", controller.MessageChat)
	//apiRouter.POST("/message/action/", controller.MessageAction)

}
