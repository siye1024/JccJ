package api

import (
	"context"
	"dousheng/db"
	"dousheng/pkg/pack"
	"dousheng/rpcserver/kitex_gen/comment"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type CommentListService struct {
	ctx context.Context
}

// NewCommentActionService new CommentActionService
func NewCommentListService(ctx context.Context) *CommentListService {
	return &CommentListService{
		ctx: ctx,
	}
}

// CommentList return comment list
func (s *CommentListService) CommentList(req *comment.DouyinCommentListRequest, fromID int64) ([]*comment.Comment, error) {
	Comments, err := db.GetVideoComments(s.ctx, req.VideoId)
	if err != nil {
		err := kerrors.NewBizStatusError(30006, "Get the Comment List Error")
		return nil, err
	}

	comments, err := pack.Comments(s.ctx, Comments, fromID)
	if err != nil {
		err := kerrors.NewBizStatusError(30006, "Get the Comment List Error")
		return nil, err
	}

	return comments, nil
}
