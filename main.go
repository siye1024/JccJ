package main

import (
	"dousheng/controller/xrpc"
	"dousheng/db"
	usersrv "dousheng/rpcserver/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup
	wg.Add(4)

	go func() { // INIT HTTP
		defer wg.Done()
		fmt.Println("Start Gin HTTP service...")
		r := gin.Default()
		initRouter(r)
		r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	}()

	go func() { // INIT DB
		defer wg.Done()
		db.InitDB()
	}()

	go func() { // INIT User RPC server
		defer wg.Done()
		var userServer usersrv.UserSrvImpl
		defer userServer.Stop()
		userServer.Start()
	}()

	go func() { // INIT User RPC server
		defer wg.Done()
		time.Sleep(time.Second * 2)
		xrpc.InitRpcClient()
	}()

	wg.Wait()

}
