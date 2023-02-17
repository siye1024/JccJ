package main

import (
	publish "dousheng/rpcserver/publish/kitex_gen/publish/publishsrv"
	"log"
)

func main() {
	svr := publish.NewServer(new(PublishSrvImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
