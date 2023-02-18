package api

import (
	"context"
	"dousheng/db"
	"dousheng/rpcserver/kitex_gen/comment"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type CommentActionService struct {
	ctx context.Context
}

// NewCommentActionService new CommentActionService
func NewCommentActionService(ctx context.Context) *CommentActionService {
	return &CommentActionService{ctx: ctx}
}

// CommentActionService action comment.
func (s *CommentActionService) CommentAction(req *comment.DouyinCommentActionRequest) error {
	// 1-评论
	if req.ActionType == 1 {
		return db.NewComment(s.ctx, &db.Comment{
			UserID:  int(req.UserId),
			VideoID: int(req.VideoId),
			Content: *req.CommentText,
		})
	}
	// 2-删除评论
	if req.ActionType == 2 {
		return db.DelComment(s.ctx, *req.CommentId, req.VideoId)
	}

	err := kerrors.NewBizStatusError(30010, "Error occurred while binding the request body to the struct")
	return err
}
