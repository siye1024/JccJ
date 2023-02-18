/*
	resp.StatusCode		resp.StatusMsg

- 	0					success
-	40001				Error UserID
-	40002				Error Action Type
-	40003				Error Token or UserID
-	40004				Invalid Token
-	40005				Error occurred while binding the request body to the struct
-	40006				Database Error
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
		respStatusCode int
		respStatusMsg  string
	)
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")

	tid, err := strconv.Atoi(to_user_id)
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": 40001,
			"status_msg":  "Error UserID",
		})
		return
	}

	act, err := strconv.Atoi(action_type)
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": 40002,
			"status_msg":  "Error Action Type",
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
		if isBizErr == false { // if it is not business error
			respStatusCode = -1
			respStatusMsg = "Service Process Error"
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
func RelationFollowList(c *gin.Context) {
	var (
		param          UserParam
		respStatusCode int
		respStatusMsg  string
	)
	uid, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": 40001,
			"status_msg":  "Error UserID",
		})
		return
	}
	param.UserId = int64(uid)
	param.Token = c.Query("token")

	if len(param.Token) == 0 || param.UserId < 0 {
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
		if isBizErr == false { // if it is not business error
			respStatusCode = -1
			respStatusMsg = "Service Process Error"
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
func RelationFollowerList(c *gin.Context) {
	var (
		param          UserParam
		respStatusCode int
		respStatusMsg  string
	)
	uid, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": 40001,
			"status_msg":  "Error UserID",
		})
		return
	}
	param.UserId = int64(uid)
	param.Token = c.Query("token")

	if len(param.Token) == 0 || param.UserId < 0 {
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
		if isBizErr == false { // if it is not business error
			respStatusCode = -1
			respStatusMsg = "Service Process Error"
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
