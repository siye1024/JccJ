package main

import (
	comment "dousheng/rpcserver/comment/kitex_gen/comment/commentsrv"
	"log"
)

func main() {
	svr := comment.NewServer(new(CommentSrvImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
