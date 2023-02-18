/*
	resp.StatusCode		resp.StatusMsg

- 	0					success
-	30001				Error Video ID
-	30002				Error Action Type
-	30003				Error Comment ID
-	30004				Comment Error
-	30005				Error Token or VideoID
-	30006				Get the Comment List Error
-	30007				Get Token Error
-	30008				Comment Operation Error
-	30009				Database Error
-	30010				Error occurred while binding the request body to the struct
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
		respStatusCode int
		respStatusMsg  string
	)
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")

	vid, err := strconv.Atoi(video_id)
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": 30001,
			"status_msg":  "Error Video ID",
		})
		return
	}
	action, err := strconv.Atoi(action_type)
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": 30002,
			"status_msg":  "Error Action Type",
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
	} else {
		comment_id := c.Query("comment_id")
		com_id, err := strconv.Atoi(comment_id)
		if err != nil {
			SendResponse(c, gin.H{
				"status_code": 30003,
				"status_msg":  "Error Comment ID",
			})
			return
		}
		com_id64 := int64(com_id)
		rpcReq.CommentId = &com_id64
	}

	resp, err := xrpc.CommentAction(c, &rpcReq)
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

// CommentList deliver the context of the "get comments" Op to the client of the RPC, and get the response
func CommentList(c *gin.Context) {
	var (
		param          CommentListParam
		respStatusCode int
		respStatusMsg  string
	)
	videoid, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": 30001,
			"status_msg":  "Error Video ID",
		})
		return
	}
	param.VideoId = int64(videoid)
	param.Token = c.Query("token")

	if len(param.Token) == 0 || param.VideoId < 0 {
		SendResponse(c, gin.H{
			"status_code": 30005,
			"status_msg":  "Error Token or VideoID",
		})
		return
	}

	resp, err := xrpc.CommentList(c, &comment.DouyinCommentListRequest{
		VideoId: param.VideoId,
		Token:   param.Token,
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
