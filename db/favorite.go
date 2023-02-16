package db

import (
	"context"
)

func GetFavoriteRelation(ctx context.Context, uid int64, vid int64) (*Video, error) {
	user := new(User)
	if err := DB.WithContext(ctx).First(user, uid).Error; err != nil {
		return nil, err
	}

	video := new(Video)
	// if err := DB.WithContext(ctx).First(&video, vid).Error; err != nil {
	// 	return nil, err
	// }

	if err := DB.WithContext(ctx).Model(&user).Association("FavoriteVideos").Find(&video, vid); err != nil {
		return nil, err
	}
	return video, nil
}
