/*
	resp.StatusCode		resp.StatusMsg

- 	0					success
-	20001				Convert Error
-	20002				Videos Pack Error
-	-1					Service Process Error
*/

package xhttp

import (
	"dousheng/controller/xrpc"
	"dousheng/rpcserver/feed/kitex_gen/feed"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// 传递 获取用户视频流操作 的上下文至 Feed 服务的 RPC 客户端, 并获取相应的响应.
func Feed(c *gin.Context) {
	var (
		feedVar        FeedParam
		latestTime     int64
		token          string = "dsdasd"
		respStatusCode        = 20001
		respStatusMsg         = "Convert Error"
	)
	//check latest time here because we need to do Atoi
	lastst_time := c.Query("latest_time")
	if len(lastst_time) != 0 {
		if parsetime, err := strconv.Atoi(lastst_time); err != nil {
			SendResponse(c, gin.H{"status_code": respStatusCode, "status_msg": respStatusMsg})
			return
		} else { // valid latest time
			latestTime = int64(parsetime)
		}
	} else { // empty latest time, choose current time
		latestTime = time.Now().UnixMilli()
	}

	feedVar.LatestTime = &latestTime

	token = c.Query("token")
	feedVar.Token = &token

	resp, err := xrpc.Feed(c, &feed.DouyinFeedRequest{
		LatestTime: feedVar.LatestTime,
		Token:      feedVar.Token,
	})
	_, isBizErr := kerrors.FromBizStatusError(err)
	if isBizErr == true || err != nil {

		SendResponse(c, gin.H{"status_code": respStatusCode, "status_msg": respStatusMsg})
		return
	}
	SendResponse(c, resp)
}
