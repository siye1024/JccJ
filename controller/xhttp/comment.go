/*
	resp.StatusCode		resp.StatusMsg

- 	0					success
-	30001				PLease Log In First!
-	30002				Invalid Action
-	30003				Invalid Token or User ID
-	30004				Invalid Comment ID
-	30005				Invalid Video
-	30009				Database Error
- 	-1					Comment Service Error
*/
package xhttp

import (
	"dousheng/controller/xrpc"
	"dousheng/rpcserver/kitex_gen/comment"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

// CommentAction deliver the context of the comment Op to the client of the RPC, and get the response
func CommentAction(c *gin.Context) {
	var (
		param          CommentActionParam
		respStatusCode = -1
		respStatusMsg  = "Comment Service Error"
	)
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")

	if len(token) == 0 {
		respStatusCode = 30001
		respStatusMsg = "PLease Log In First!"
		SendResponse(c, gin.H{
			"status_code": respStatusCode,
			"status_msg":  respStatusMsg,
		})
		return
	}

	vid, err := strconv.Atoi(video_id)
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": respStatusCode,
			"status_msg":  respStatusMsg,
		})
		return
	}
	action, err := strconv.Atoi(action_type)
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": respStatusCode,
			"status_msg":  respStatusMsg,
		})
		return
	}

	param.Token = token
	param.VideoId = int64(vid)
	param.ActionType = int32(action)

	rpcReq := comment.DouyinCommentActionRequest{
		VideoId:    param.VideoId,
		Token:      param.Token,
		ActionType: param.ActionType,
	}

	if action == 1 {
		comment_text := c.Query("comment_text")
		rpcReq.CommentText = &comment_text
	} else if action == 2 {
		comment_id := c.Query("comment_id")
		com_id, err := strconv.Atoi(comment_id)
		if err != nil {
			SendResponse(c, gin.H{
				"status_code": respStatusCode,
				"status_msg":  respStatusMsg,
			})
			return
		}
		com_id64 := int64(com_id)
		rpcReq.CommentId = &com_id64
	} else {
		SendResponse(c, gin.H{
			"status_code": 30002,
			"status_msg":  "Invalid Action",
		})
		return
	}

	resp, err := xrpc.CommentAction(c, &rpcReq)
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

// CommentList deliver the context of the "get comments" Op to the client of the RPC, and get the response
func CommentList(c *gin.Context) {
	var (
		param          CommentListParam
		respStatusCode = -1
		respStatusMsg  = "Comment Service Error"
	)
	videoid, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": respStatusCode,
			"status_msg":  respStatusMsg,
		})
		return
	}
	param.VideoId = int64(videoid)
	param.Token = c.Query("token")

	if param.VideoId <= 0 {
		SendResponse(c, gin.H{
			"status_code": 30005,
			"status_msg":  "Invalid Video",
		})
		return
	}

	resp, err := xrpc.CommentList(c, &comment.DouyinCommentListRequest{
		VideoId: param.VideoId,
		Token:   param.Token,
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
