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
func (s *CommentActionService) CommentAction(req *comment.DouyinCommentActionRequest) (int64, error) {
	// 1-评论
	if req.ActionType == 1 {
		commentId, err := db.NewComment(s.ctx, &db.Comment{
			UserID:  req.UserId,
			VideoID: req.VideoId,
			Content: *req.CommentText,
		})
		if err != nil {
			return commentId, err
		}
		return commentId, err
	}
	// 2-删除评论
	if req.ActionType == 2 {
		return *req.CommentId, db.DelComment(s.ctx, *req.CommentId, req.VideoId)
	}

	return 0, kerrors.NewBizStatusError(30002, "Invalid Action")
}
