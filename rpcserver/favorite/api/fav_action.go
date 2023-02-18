package api

import (
	"context"
	"dousheng/db"
	"dousheng/rpcserver/kitex_gen/favorite"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type FavoriteActionOp struct {
	ctx context.Context
}

// NewFavoriteActionOp new FavoriteActionOp
func NewFavoriteActionOp(ctx context.Context) *FavoriteActionOp {
	return &FavoriteActionOp{ctx: ctx}
}

// FavoriteAction action favorite.
func (s *FavoriteActionOp) FavoriteAction(req *favorite.DouyinFavoriteActionRequest) error {
	// 1-点赞
	if req.ActionType == 1 {
		return db.Favorite(s.ctx, req.UserId, req.VideoId)
	}
	// 2-取消点赞
	if req.ActionType == 2 {
		return db.DisFavorite(s.ctx, req.UserId, req.VideoId)
	}

	err := kerrors.NewBizStatusError(70005, "Error occurred while binding the request body to the struct")
	return err
}
