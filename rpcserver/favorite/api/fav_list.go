package api

import (
	"context"
	"dousheng/db"
	"dousheng/pkg/pack"
	"dousheng/rpcserver/kitex_gen/favorite"
	"dousheng/rpcserver/kitex_gen/feed"
)

type FavoriteListOp struct {
	ctx context.Context
}

// NewFavoriteListOp creates a new FavoriteListOp
func NewFavoriteListOp(ctx context.Context) *FavoriteListOp {
	return &FavoriteListOp{
		ctx: ctx,
	}
}

// FavoriteList returns a Favorite List
func (s *FavoriteListOp) FavoriteList(req *favorite.DouyinFavoriteListRequest) ([]*feed.Video, error) {
	FavoriteVideos, err := db.FavoriteList(s.ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	videos, err := pack.FavoriteVideos(s.ctx, FavoriteVideos, &req.UserId)
	if err != nil {
		return nil, err
	}
	return videos, nil
}
