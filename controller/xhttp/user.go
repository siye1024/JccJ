/*
	resp.StatusCode		resp.StatusMsg

- 	0					success
-	10001				Empty Username or Password
-	10002				User Already Exist
-	10003				JWT ERROR:That's not even a token
-	10004				JWT ERROR:Token expired
-	10005				JWT ERROR:Token is not active yet
-	10006				JWT ERROR:Couldn't handle this token
-	-1					Service Process Error
*/
package xhttp

import (
	"dousheng/controller/xrpc"
	"dousheng/rpcserver/user/kitex_gen/user"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
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

// Login Op
func Login(c *gin.Context) {
	//client login metric
	var (
		logMsg         UserRegisterParam
		respStatusCode int
		respStatusMsg  string
	)
	logMsg.UserName = c.Query("username")
	logMsg.PassWord = c.Query("password")

	if len(logMsg.UserName) == 0 || len(logMsg.PassWord) == 0 {
		SendResponse(c, gin.H{
			"status_code": 10001,
			"status_msg":  "Empty Username or Password",
		})
		return
	}

	resp, err := xrpc.Login(c, &user.DouyinUserRegisterRequest{
		Username: logMsg.UserName,
		Password: logMsg.PassWord,
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
		//log.Println("see here", resp, err, bizErr, isBizErr)
		//SendResponse(c, resp)
		return
	}

	SendResponse(c, resp) // service success

}

// Get User's ID OP
func GetUserById(c *gin.Context) {
	//client request metric
	var (
		getUserByIdMsg UserParam
		respStatusCode int
		respStatusMsg  string
	)
	userid, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		SendResponse(c, gin.H{
			"status_code": 10003,
			"status_msg":  "That's not even a token",
		})
	}
	getUserByIdMsg.UserId = int64(userid)
	getUserByIdMsg.Token = c.Query("token")

	// if username == empty or password == empty, Actually this case has been processed by dousheng client
	if len(getUserByIdMsg.Token) == 0 || getUserByIdMsg.UserId < 0 {
		SendResponse(c, gin.H{
			"status_code": 10006,
			"status_msg":  "Couldn't handle this token",
		})
		return
	}

	// transport to rpc client
	resp, err := xrpc.GetUserById(c, &user.DouyinUserRequest{
		UserId: getUserByIdMsg.UserId,
		Token:  getUserByIdMsg.Token,
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
		//log.Println("see here", resp, err, bizErr, isBizErr)
		//SendResponse(c, resp)
		return
	}

	SendResponse(c, resp) // service success
}
