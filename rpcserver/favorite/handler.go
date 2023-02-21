package favoritesrv

import (
	"context"
	"dousheng/controller/xhttp"
	"dousheng/rpcserver/favorite/api"
	favorite "dousheng/rpcserver/kitex_gen/favorite"
	"dousheng/rpcserver/kitex_gen/favorite/favoritesrv"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"log"
	"net"
)

// FavoriteSrvImpl implements the last service interface defined in the IDL.
type FavoriteSrvImpl struct{}

// FavoriteAction implements the FavoriteSrvImpl interface.
func (s *FavoriteSrvImpl) FavoriteAction(ctx context.Context, req *favorite.DouyinFavoriteActionRequest) (resp *favorite.DouyinFavoriteActionResponse, err error) {
	var (
		respMsg = "Favortie Action Success"
	)
	claim, err := xhttp.Jwt.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}

	if claim.Id > 0 {
		req.UserId = claim.Id
	} else {
		return nil, kerrors.NewBizStatusError(70003, "Invalid Token or User ID")
	}

	if req.ActionType != 1 && req.ActionType != 2 || req.UserId <= 0 || req.VideoId <= 0 {
		err := kerrors.NewBizStatusError(70005, "Invalid Request")
		return nil, err
	}

	err = api.NewFavoriteActionOp(ctx).FavoriteAction(req)
	if err != nil {
		return nil, err
	}

	resp = &favorite.DouyinFavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  &respMsg,
	}
	return resp, nil
}

// FavoriteList implements the FavoriteSrvImpl interface.
func (s *FavoriteSrvImpl) FavoriteList(ctx context.Context, req *favorite.DouyinFavoriteListRequest) (resp *favorite.DouyinFavoriteListResponse, err error) {
	var (
		respMsg = "Get Favortie List Success"
		user_id int64
	)
	if len(req.Token) == 0 {
		user_id = 0
	} else {
		claim, err := xhttp.Jwt.ParseToken(req.Token)
		if err != nil {
			return nil, err
		}
		user_id = claim.Id
	}

	if req.UserId <= 0 || user_id < 0 { // if user can only see his own fav list, need req.UserId != claim.Id
		err := kerrors.NewBizStatusError(70003, "Invalid Token or User ID")
		return nil, err
	}

	videos, err := api.NewFavoriteListOp(ctx).FavoriteList(req)
	if err != nil {
		return resp, nil
	}

	resp = &favorite.DouyinFavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  &respMsg,
		VideoList:  videos,
	}
	return resp, nil
}

func (s *FavoriteSrvImpl) Start() {
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Panic(err)
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:0") //random empty port
	svr := favoritesrv.NewServer(new(FavoriteSrvImpl),
		server.WithServiceAddr(addr),                                       // address
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler),            //support kerrors
		server.WithRegistry(r),                                             // registry
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}), // limit
		server.WithMuxTransport(),                                          // Multiplex
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "favoriteService"}),
	)

	log.Println("Start Favorite RPC service...")
	err = svr.Run()
	if err != nil {
		log.Panic(err.Error())
	}

}
func (s *FavoriteSrvImpl) Stop() {
	log.Println("Stop Favorite  RPC service...")
}
