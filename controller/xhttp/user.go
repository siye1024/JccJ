/*
	resp.StatusCode		resp.StatusMsg

- 	0					success
-	10001					Empty Username or Password
-	10002					User Already Exist
-	-1					Service Process Error
*/
package xhttp

import (
	"dousheng/controller/xrpc"
	"dousheng/rpcserver/user/kitex_gen/user"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/gin-gonic/gin"
	"log"
)

func Register(c *gin.Context) {
	//client request metric
	var (
		registerMsg    UserRegisterParam
		respStatusCode int
		respStatusMsg  string
	)
	registerMsg.UserName = c.Query("username")
	registerMsg.PassWord = c.Query("password")
	// if username == empty or password == empty, Actually this case has been processed by user client
	if len(registerMsg.UserName) == 0 || len(registerMsg.PassWord) == 0 {
		SendResponse(c, gin.H{
			"status_code": 10001,
			"status_msg":  "Empty Username or Password",
		})
		return
	}

	// transport to rpc client
	resp, err := xrpc.Register(c, &user.DouyinUserRegisterRequest{
		Username: registerMsg.UserName,
		Password: registerMsg.PassWord,
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

	SendResponse(c, resp) // service success
}
