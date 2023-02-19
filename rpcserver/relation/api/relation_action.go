package api

import (
	"context"
	"dousheng/db"
	"dousheng/rpcserver/kitex_gen/relation"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type RelationActionOp struct {
	ctx context.Context
}

// NewRelationActionOp new RelationActionOp
func NewRelationActionOp(ctx context.Context) *RelationActionOp {
	return &RelationActionOp{ctx: ctx}
}

// RelationAction action favorite.
func (s *RelationActionOp) RelationAction(req *relation.DouyinRelationActionRequest) error {
	// 1-关注
	if req.ActionType == 1 {
		return db.NewRelation(s.ctx, req.UserId, req.ToUserId)
	}
	// 2-取消关注
	if req.ActionType == 2 {
		return db.DisRelation(s.ctx, req.UserId, req.ToUserId)
	}

	return kerrors.NewBizStatusError(40004, "Error Action Type")
}
