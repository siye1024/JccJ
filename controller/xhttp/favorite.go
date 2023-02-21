/*
	resp.StatusCode		resp.StatusMsg

- 	0					success
-	70001				PLease Log In First!
-	70002				You have Liked!
-	70003				Invalid Token or User ID
-	70004				Invalid Action
-	70005				Invalid Request
-	70006				Database Error

-	70001				Error Video ID
-	70002				Error Action Type
-	70003				Invalid Token or User ID
-	70004				Invalid Token
-	70005				Invalid Request
-	70006				Database Error
- 	-1					Favorite Operation Error
*/

package xhttp

import (
	"dousheng/controller/xrpc"
	"dousheng/rpcserver/kitex_gen/favorite"
	_ "github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/kerrors"
	_ "github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

// 传递 点赞操作 的上下文至 Favorite 服务的 RPC 客户端, 并获取相应的响应.
func FavoriteAction(c *gin.Context) {
	var (
		param          FavoriteActionParam
		respStatusCode = -1
		respStatusMsg  = "Favorite Operation Error"
	)
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")

	if len(token) == 0 {
		respStatusCode = 70001
		respStatusMsg = "PLease Log In First!"
		SendResponse(c, gin.H{
			"status_code": respStatusCode,
			"status_msg":  respStatusMsg,
		})
		return
	}

	vid, err := strconv.Atoi(video_id)
	if err != nil {
		log.Println(err)
		SendResponse(c, gin.H{
			"status_code": respStatusCode,
			"status_msg":  respStatusMsg,
		})
		return
	}
	act, err := strconv.Atoi(action_type)
	if err != nil {
		log.Println(err)
		SendResponse(c, gin.H{
			"status_code": respStatusCode,
			"status_msg":  respStatusCode,
		})
		return
	}

	param.Token = token
	param.VideoId = int64(vid)
	param.ActionType = int32(act)

	resp, err := xrpc.FavoriteAction(c, &favorite.DouyinFavoriteActionRequest{
		VideoId:    param.VideoId,
		Token:      param.Token,
		ActionType: param.ActionType,
	})
	bizErr, isBizErr := kerrors.FromBizStatusError(err)
	if isBizErr == true || err != nil {
		if isBizErr == false { // if it is not business error, return -1 default error
			log.Println(err.Error())
		} else { // business err
			respStatusCode = int(bizErr.BizStatusCode())
			respStatusMsg = bizErr.BizMessage()
		}
		SendResponse(c, gin.H{
			"status_code": respStatusCode,
			"status_msg":  respStatusMsg,
		})

		return
	}

	SendResponse(c, resp)
}

// 传递 获取点赞列表操作 的上下文至 Favorite 服务的 RPC 客户端, 并获取相应的响应.
func FavoriteList(c *gin.Context) {
	var (
		param          UserParam
		respStatusCode = -1
		respStatusMsg  = "Get Favorite List Error"
	)
	userid, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": respStatusCode,
			"status_msg":  respStatusMsg,
		})
		return
	}
	param.UserId = int64(userid)
	param.Token = c.Query("token")

	if param.UserId <= 0 {
		SendResponse(c, gin.H{
			"status_code": 70003,
			"status_msg":  "Invalid Token or User ID",
		})
		return
	}

	resp, err := xrpc.FavoriteList(c, &favorite.DouyinFavoriteListRequest{
		UserId: param.UserId,
		Token:  param.Token,
	})
	bizErr, isBizErr := kerrors.FromBizStatusError(err)
	if isBizErr == true || err != nil {
		if isBizErr == false { // if it is not business error, return -1 default error
			log.Println(err.Error())
		} else { // business err
			respStatusCode = int(bizErr.BizStatusCode())
			respStatusMsg = bizErr.BizMessage()
		}
		SendResponse(c, gin.H{
			"status_code": respStatusCode,
			"status_msg":  respStatusMsg,
		})

		return
	}

	SendResponse(c, resp)
}
