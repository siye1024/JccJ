/*
	resp.StatusCode		resp.StatusMsg

- 	0					success
-	10001				Empty Username or Password
-	10002				User Already Exist
-	10003				That's not even a token
-	10004				Token expired
-	10005				Token is not active yet
-	10006				Couldn't handle this token
-	10007				Invalid Username
-	10008				Invalid Username or Password
-	10009				User Already Withdraw
-	10010				Encoded Hash isn't in the correct format
-	10011				Encoded Hash isn't in the correct version
-	10012				Signature is Invalid
-	10013				Register Failed
-	10014				Too Long Username or Password
-	-1					Service Process Error
*/
package xhttp

import (
	"dousheng/controller/xrpc"
	"dousheng/rpcserver/kitex_gen/user"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func Register(c *gin.Context) {
	//client request metric
	var (
		registerMsg    UserRegisterParam
		respStatusCode int    = -1
		respStatusMsg  string = "Register Error"
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
	// len of username and password can't exceed 32 bit
	if len(registerMsg.UserName) > 32 || len(registerMsg.PassWord) > 32 {
		SendResponse(c, gin.H{
			"status_code": 10014,
			"status_msg":  "Too Long Username or Password",
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

	SendResponse(c, resp) // service success
}

// Login Op
func Login(c *gin.Context) {
	//client login metric
	var (
		logMsg         UserRegisterParam
		respStatusCode = -1
		respStatusMsg  = "Login Error"
	)
	logMsg.UserName = c.Query("username")
	logMsg.PassWord = c.Query("password")
	// username or password is empty
	if len(logMsg.UserName) == 0 || len(logMsg.PassWord) == 0 {
		SendResponse(c, gin.H{
			"status_code": 10001,
			"status_msg":  "Empty Username or Password",
		})
		return
	}
	// len of the username or password exceed the 32
	if len(logMsg.UserName) > 32 || len(logMsg.PassWord) > 32 {
		SendResponse(c, gin.H{
			"status_code": 10014,
			"status_msg":  "Too Long Username or Password",
		})
		return
	}

	resp, err := xrpc.Login(c, &user.DouyinUserRegisterRequest{
		Username: logMsg.UserName,
		Password: logMsg.PassWord,
	})
	bizErr, isBizErr := kerrors.FromBizStatusError(err)
	if isBizErr == true || err != nil {
		if isBizErr == false { //  if it is not business error, return -1 default error
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

// Get User Info By ID
func GetUserById(c *gin.Context) {
	//client request metric
	var (
		getUserByIdMsg UserParam
		respStatusCode = -1
		respStatusMsg  = "Get User Info Error"
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

	getUserByIdMsg.UserId = int64(userid)
	getUserByIdMsg.Token = c.Query("token")

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
		if isBizErr == false { //  if it is not business error, return -1 default error
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
