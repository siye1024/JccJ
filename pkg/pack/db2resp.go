package pack

import (
	"context"
	"dousheng/db"
	"dousheng/rpcserver/feed/kitex_gen/feed"
	"dousheng/rpcserver/feed/kitex_gen/user"
	"errors"
	"gorm.io/gorm"
)

// fromID is the uid of the request user, not video author
func PackVideos(ctx context.Context, vs []*db.Video, fromID *int64) ([]*feed.Video, error) {
	videos := make([]*feed.Video, 0)
	for _, v := range vs {
		video2, err := PackVideo(ctx, v, *fromID)
		if err != nil {
			return nil, err
		}

		videos = append(videos, video2)
	}

	return videos, nil
}

func PackVideo(ctx context.Context, v *db.Video, fromID int64) (*feed.Video, error) {
	if v == nil {
		return nil, nil
	}
	video_author, err := db.GetUserByID(ctx, int64(v.AuthorID))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	author, err := PackUser(ctx, video_author, fromID)
	if err != nil {
		return nil, err
	}
	favorite_count := int64(v.FavoriteCount)
	comment_count := int64(v.CommentCount)

	isFavorite := false
	if fromID != 0 { // tourist uid = 0, login user uid != 0
		results, err := db.GetFavoriteRelation(ctx, fromID, int64(v.ID))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		} else if results != nil && results.AuthorID != 0 {
			isFavorite = true
		}
	}
	return &feed.Video{
		Id:            int64(v.ID),
		Author:        author,
		PlayUrl:       v.PlayUrl,
		CoverUrl:      v.CoverUrl,
		FavoriteCount: favorite_count,
		CommentCount:  comment_count,
		IsFavorite:    isFavorite,
		Title:         v.Title,
	}, nil
}

func PackUsers(ctx context.Context, us []*db.User, fromID int64) ([]*user.User, error) {
	users := make([]*user.User, 0)
	for _, u := range us {
		user2, err := PackUser(ctx, u, fromID)
		if err != nil {
			return nil, err
		}

		if user2 != nil {
			users = append(users, user2)
		}
	}
	return users, nil
}
func PackUser(ctx context.Context, u *db.User, fromID int64) (*user.User, error) {
	if u == nil {
		return &user.User{
			Name: "non-exist user",
		}, nil
	}

	follow_count := int64(u.FollowCount)
	follower_count := int64(u.FollowerCount)

	// true means fromID has followed u.ID
	isFollow := false
	relation, err := db.GetRelation(ctx, fromID, int64(u.ID))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if relation != nil {
		isFollow = true
	}

	return &user.User{
		Id:            int64(u.ID),
		Name:          u.UserName,
		FollowCount:   &follow_count,
		FollowerCount: &follower_count,
		IsFollow:      isFollow,
	}, nil
}
