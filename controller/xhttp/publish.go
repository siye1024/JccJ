/*
	resp.StatusCode		resp.StatusMsg

- 	0					success
-	-1					Service Process Error
*/

package xhttp

import (
	"bytes"
	"dousheng/controller/xrpc"
	"dousheng/rpcserver/publish/kitex_gen/publish"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/gin-gonic/gin"
	"io"
	"log"
)

func PublishAction(c *gin.Context) {
	var (
		paramVar       PublishActionParam
		token, title   string
		respStatusCode = -1
		respStatusMsg  = "Video Publish Error"
	)
	token = c.PostForm("token")
	title = c.PostForm("title")

	_, fileHeader, err := c.Request.FormFile("data")
	if err != nil {
		log.Println(err)
		SendResponse(c, gin.H{"status_code": respStatusCode, "status_msg": respStatusMsg})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		log.Println(err)
		SendResponse(c, gin.H{"status_code": respStatusCode, "status_msg": respStatusMsg})
		return
	}

	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Println(err)
		SendResponse(c, gin.H{"status_code": respStatusCode, "status_msg": respStatusMsg})
		return
	}

	paramVar.Token = token
	paramVar.Title = title

	resp, err := xrpc.PublishAction(c, &publish.DouyinPublishActionRequest{
		Title: paramVar.Title,
		Token: paramVar.Token,
		Data:  buf.Bytes(),
	})
	bizErr, isBizErr := kerrors.FromBizStatusError(err)
	if isBizErr == true || err != nil {
		if isBizErr == false { // if it is not business error
			log.Println(err)
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

	SendResponse(c, resp) // service success
}

func PublishList(c *gin.Context) {

}