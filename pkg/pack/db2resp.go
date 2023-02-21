package pack

import (
	"context"
	"dousheng/db"
	"dousheng/rpcserver/kitex_gen/comment"
	"dousheng/rpcserver/kitex_gen/feed"
	"dousheng/rpcserver/kitex_gen/user"
	"errors"
	"gorm.io/gorm"
)

// fromID is the uid of the request user, possible not video author
func Videos(ctx context.Context, vs []*db.Video, fromID *int64) ([]*feed.Video, error) {
	videos := make([]*feed.Video, 0)
	for _, v := range vs {
		video2, err := Video(ctx, v, *fromID)
		if err != nil {
			return nil, err
		}

		videos = append(videos, video2)
	}

	return videos, nil
}

func Video(ctx context.Context, v *db.Video, fromID int64) (*feed.Video, error) {
	if v == nil {
		return nil, nil
	}
	video_author, err := db.GetUserByID(ctx, int64(v.AuthorID))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	author, err := User(ctx, video_author, fromID)
	if err != nil {
		return nil, err
	}
	favorite_count := int64(v.FavoriteCount)
	comment_count := int64(v.CommentCount)

	isFavorite := false
	if fromID > 0 { // tourist uid = 0, login user uid > 0
		isFav, err := db.GetFavoriteRelation(ctx, fromID, int64(v.ID))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if isFav == true {
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

func Users(ctx context.Context, us []*db.User, fromID int64) ([]*user.User, error) {
	users := make([]*user.User, 0)
	for _, u := range us {
		user2, err := User(ctx, u, fromID)
		if err != nil {
			return nil, err
		}

		if user2 != nil {
			users = append(users, user2)
		}
	}
	return users, nil
}
func User(ctx context.Context, u *db.User, fromID int64) (*user.User, error) {
	if u == nil {
		return &user.User{
			Name: "non-exist user",
		}, nil
	}

	follow_count := u.FollowCount
	follower_count := u.FollowerCount

	// true means fromID has followed u.ID
	isFollow := false
	if fromID == int64(u.ID) { // I have followed myself
		isFollow = true
	} else {
		relation, err := db.GetRelation(ctx, fromID, int64(u.ID))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if relation != nil && relation.UserID != 0 { //double check is necessary
			isFollow = true
		}
	}
	works, err := db.PublishList(ctx, int64(u.ID))
	if err != nil {
		return nil, err
	}
	workCount := int64(len(works))

	favs, err := db.FavoriteList(ctx, int64(u.ID))
	if err != nil {
		return nil, err
	}
	favCount := int64(len(favs))
	return &user.User{
		Id:             int64(u.ID),
		Name:           u.UserName,
		FollowCount:    &follow_count,
		FollowerCount:  &follower_count,
		IsFollow:       isFollow,
		WorkCount:      &workCount,
		FavoriteCount:  &favCount,
		TotalFavorited: &u.FavoritedCount,
	}, nil
}

// Comments pack Comments info.
func Comments(ctx context.Context, vs []*db.Comment, fromID int64) ([]*comment.Comment, error) {
	comments := make([]*comment.Comment, 0)
	for _, v := range vs {
		commentUser, err := db.GetUserByID(ctx, v.UserID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		packUser, err := User(ctx, commentUser, fromID)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &comment.Comment{
			Id:         int64(v.ID),
			User:       packUser,
			Content:    v.Content,
			CreateDate: v.CreatedAt.Format("01-02"),
		})
	}
	return comments, nil
}

// FollowingList pack lists of following info.
func FollowingList(ctx context.Context, vs []*db.Relation, fromID int64) ([]*user.User, error) {
	users := make([]*db.User, 0)
	for _, v := range vs {
		user2, err := db.GetUserByID(ctx, int64(v.ToUserID))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		users = append(users, user2)
	}

	return Users(ctx, users, fromID)
}

// FollowerList pack lists of follower info.
func FollowerList(ctx context.Context, vs []*db.Relation, fromID int64) ([]*user.User, error) {
	users := make([]*db.User, 0)
	for _, v := range vs {
		user2, err := db.GetUserByID(ctx, int64(v.UserID))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		users = append(users, user2)
	}

	return Users(ctx, users, fromID)
}

// FavoriteVideos pack favoriteVideos info.
func FavoriteVideos(ctx context.Context, vs []*db.Video, uid *int64) ([]*feed.Video, error) {
	videos := make([]*db.Video, 0) // db.Video -> * dbVideo
	for i, _ := range vs {
		videos = append(videos, vs[i])
	}

	packVideos, err := PackFaVideos(ctx, vs, uid)
	if err != nil {
		return nil, err
	}
	return packVideos, nil
}
func PackFaVideos(ctx context.Context, vs []*db.Video, fromID *int64) ([]*feed.Video, error) {
	videos := make([]*feed.Video, 0)
	for _, v := range vs {
		video_author, err := db.GetUserByID(ctx, int64(v.AuthorID))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		author, err := User(ctx, video_author, *fromID)
		if err != nil {
			return nil, err
		}
		favorite_count := int64(v.FavoriteCount)
		comment_count := int64(v.CommentCount)
		video2 := &feed.Video{
			Id:            int64(v.ID),
			Author:        author,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: favorite_count,
			CommentCount:  comment_count,
			IsFavorite:    true, //because it is favoritelist!!!
			Title:         v.Title,
		}

		videos = append(videos, video2)
	}

	return videos, nil
}
