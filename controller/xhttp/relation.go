/*
	resp.StatusCode		resp.StatusMsg

- 	0					success
-	40001				PLease Log In First!
-	40002				Invalid To_UserId
-	40003				Error Token or UserID
-	40004				Invalid Action Type
-	40006				Database Error
-	-1					Default Error
*/

package xhttp

import (
	"dousheng/controller/xrpc"
	"dousheng/rpcserver/kitex_gen/relation"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

// 传递 关注操作 的上下文至 Relation 服务的 RPC 客户端, 并获取相应的响应.
func RelationAction(c *gin.Context) {
	var (
		param          RelationActionParam
		respStatusCode = -1
		respStatusMsg  = "Relation Action Error"
	)
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")

	if len(token) == 0 {
		respStatusCode = 40001
		respStatusMsg = "PLease Log In First!"
		SendResponse(c, gin.H{
			"status_code": respStatusCode,
			"status_msg":  respStatusMsg,
		})
		return
	}

	tid, err := strconv.Atoi(to_user_id)
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": respStatusCode,
			"status_msg":  respStatusMsg,
		})
		return
	}
	if tid <= 0 {
		SendResponse(c, gin.H{
			"status_code": 40002,
			"status_msg":  "Invalid To_UserId",
		})
		return
	}

	act, err := strconv.Atoi(action_type)
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": respStatusCode,
			"status_msg":  respStatusMsg,
		})
		return
	}

	param.Token = token
	param.ToUserId = int64(tid)
	param.ActionType = int32(act)

	rpcReq := relation.DouyinRelationActionRequest{
		ToUserId:   param.ToUserId,
		Token:      param.Token,
		ActionType: param.ActionType,
	}

	resp, err := xrpc.RelationAction(c, &rpcReq)
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

// 传递 获取正在关注列表操作 的上下文至 Relation 服务的 RPC 客户端, 并获取相应的响应.
func FollowList(c *gin.Context) {
	var (
		param          UserParam
		respStatusCode = -1
		respStatusMsg  = "Get Follow List Error"
	)
	uid, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": respStatusCode,
			"status_msg":  respStatusMsg,
		})
		return
	}
	param.UserId = int64(uid)
	param.Token = c.Query("token")

	if param.UserId < 0 {
		SendResponse(c, gin.H{
			"status_code": 40003,
			"status_msg":  "Error Token or UserID",
		})
		return
	}

	resp, err := xrpc.RelationFollowList(c, &relation.DouyinRelationFollowListRequest{
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

// 传递 获取粉丝列表操作 的上下文至 Relation 服务的 RPC 客户端, 并获取相应的响应.
func FollowerList(c *gin.Context) {
	var (
		param          UserParam
		respStatusCode = -1
		respStatusMsg  = "Get Follower List Error"
	)
	uid, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": respStatusCode,
			"status_msg":  respStatusMsg,
		})
		return
	}
	param.UserId = int64(uid)
	param.Token = c.Query("token")

	if param.UserId < 0 {
		SendResponse(c, gin.H{
			"status_code": 40003,
			"status_msg":  "Error Token or UserID",
		})
		return
	}

	resp, err := xrpc.RelationFollowerList(c, &relation.DouyinRelationFollowerListRequest{
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

func FriendList(c *gin.Context) {
	var (
		param          UserParam
		respStatusCode = -1
		respStatusMsg  = "Get Friend List Error"
		userid         int
		err            error
	)
	user_id := c.Query("user_id")
	if len(user_id) > 0 {
		userid, err = strconv.Atoi(user_id)
		if err != nil {
			SendResponse(c, gin.H{
				"status_code": respStatusCode,
				"status_msg":  respStatusMsg,
			})
			return
		}
	} else {
		userid = 0
	}

	param.UserId = int64(userid)
	param.Token = c.Query("token")

	if len(param.Token) == 0 || param.UserId <= 0 { //can't be friend with tourist
		SendResponse(c, gin.H{
			"status_code": 40003,
			"status_msg":  "Error Token or UserID",
		})
		return
	}
	resp, err := xrpc.RelationFriendList(c, &relation.DouyinRelationFriendListRequest{
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
