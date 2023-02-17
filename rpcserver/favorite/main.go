package main

import (
	favorite "dousheng/rpcserver/favorite/kitex_gen/favorite/favoritesrv"
	"log"
)

func main() {
	svr := favorite.NewServer(new(FavoriteSrvImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}