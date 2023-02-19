package api

import (
	"context"
	"dousheng/db"
	"dousheng/pkg/pack"
	"dousheng/rpcserver/kitex_gen/relation"
	"dousheng/rpcserver/kitex_gen/user"
)

type FollowingListOp struct {
	ctx context.Context
}

// NewFollowingListOp creates a new FollowingListService
func NewFollowingListOp(ctx context.Context) *FollowingListOp {
	return &FollowingListOp{
		ctx: ctx,
	}
}

// FollowingList returns the following lists
func (s *FollowingListOp) FollowingList(req *relation.DouyinRelationFollowListRequest, fromID int64) ([]*user.User, error) {
	FollowingUser, err := db.FollowingList(s.ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return pack.FollowingList(s.ctx, FollowingUser, fromID)
}

type FollowerListOp struct {
	ctx context.Context
}

// NewFollowerListOp creates a new FollowerListService
func NewFollowerListOp(ctx context.Context) *FollowerListOp {
	return &FollowerListOp{
		ctx: ctx,
	}
}

// FollowerList returns the Follower Lists
func (s *FollowerListOp) FollowerList(req *relation.DouyinRelationFollowerListRequest, fromID int64) ([]*user.User, error) {
	FollowerUser, err := db.FollowerList(s.ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return pack.FollowerList(s.ctx, FollowerUser, fromID)
}
