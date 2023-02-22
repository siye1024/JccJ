package main

import (
	"dousheng/controller/xrpc"
	"dousheng/db"
	"dousheng/pkg/minio"
	commentsrv "dousheng/rpcserver/comment"
	favoritesrv "dousheng/rpcserver/favorite"
	feedsrv "dousheng/rpcserver/feed"
	publishsrv "dousheng/rpcserver/publish"
	relationsrv "dousheng/rpcserver/relation"
	usersrv "dousheng/rpcserver/user"

	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

func main() {

	var (
		wg sync.WaitGroup
		ip = "localhost"
	)
	wg.Add(9)

	go func() { // INIT HTTP
		defer wg.Done()
		fmt.Println("Start Gin HTTP service...")
		r := gin.Default()
		initRouter(r)
		r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	}()

	go func() { // INIT DB & Minio Client
		defer wg.Done()
		db.InitDB()
		xminio.InitMInio(ip)
	}()

	go func() { // INIT User RPC server
		defer wg.Done()
		var userServer usersrv.UserSrvImpl
		defer userServer.Stop()
		userServer.Start()
	}()

	go func() { // INIT Feed RPC server
		defer wg.Done()
		var feedServer feedsrv.FeedSrvImpl
		defer feedServer.Stop()
		feedServer.Start()
	}()

	go func() { // INIT Comment RPC server
		defer wg.Done()
		var commentServer commentsrv.CommentSrvImpl
		defer commentServer.Stop()
		commentServer.Start()
	}()

	go func() { // INIT Publish RPC server
		defer wg.Done()
		var publishServer publishsrv.PublishSrvImpl
		defer publishServer.Stop()
		publishServer.Start()
	}()

	go func() { // INIT Relation RPC server
		defer wg.Done()
		var relationServer relationsrv.RelationSrvImpl
		defer relationServer.Stop()
		relationServer.Start()
	}()

	go func() { // INIT Favorite RPC server
		defer wg.Done()
		var favoriteServer favoritesrv.FavoriteSrvImpl
		defer favoriteServer.Stop()
		favoriteServer.Start()
	}()

	go func() { // INIT All RPC client
		defer wg.Done()
		time.Sleep(time.Second * 2)
		xrpc.InitRpcClient()
	}()

	wg.Wait()

}
