package db

import (
	"context"
	"errors"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"gorm.io/gorm"
)

func GetFavoriteRelation(ctx context.Context, uid int64, vid int64) (bool, error) {
	//make sure uid and vid exist in the table
	user := new(User)
	if err := DB.WithContext(ctx).First(&user, uid).Error; err != nil {
		return false, err
	}

	video := new(Video)
	if err := DB.WithContext(ctx).Model(&user).Association("FavoriteVideos").Find(&video, vid); err != nil {
		return false, err
	}
	//please note that Asscoation(xx).Find(&xx,xx) wouldn't return gorm.ErrRecordNotFound if no finding

	if video == nil || video.AuthorID == 0 {
		return false, nil
	}
	return true, nil
}

// Favorite new favorite data.
func Favorite(ctx context.Context, uid int64, vid int64) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		//1. 新增点赞数据
		user := new(User)
		if err := tx.WithContext(ctx).First(user, uid).Error; err != nil {
			return err
		}

		// double check unexpected action, if user has liked the video, do nothing
		isFav, err := GetFavoriteRelation(ctx, uid, vid)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if isFav == true {
			return kerrors.NewBizStatusError(70002, "You have Liked! Invalid Action")
		}

		video := new(Video)
		if err := tx.WithContext(ctx).First(video, vid).Error; err != nil {
			return err
		}

		if err := tx.WithContext(ctx).Model(&user).Association("FavoriteVideos").Append(video); err != nil {
			return err
		}
		//2.改变 video 表中的 favorite count
		res := tx.Model(video).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			err := kerrors.NewBizStatusError(70006, "Database Error")
			return err
		}

		return nil
	})

	return err
}

// DisFavorite deletes the specified favorite from the database
func DisFavorite(ctx context.Context, uid int64, vid int64) error {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		//1. 删除点赞数据
		user := new(User)
		if err := tx.WithContext(ctx).First(user, uid).Error; err != nil {
			return err
		}
		// double check unexpected action
		isFav, err := GetFavoriteRelation(ctx, uid, vid)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if isFav == false {
			return kerrors.NewBizStatusError(70004, "Invalid Action")
		}

		video := new(Video)
		if err := tx.WithContext(ctx).First(video, vid).Error; err != nil {
			return err
		}
		err = tx.Unscoped().WithContext(ctx).Model(&user).Association("FavoriteVideos").Delete(video)
		if err != nil {
			return err
		}

		//2.改变 video 表中的 favorite count
		if video.FavoriteCount > 0 {
			res := tx.Model(video).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
			if res.Error != nil {
				return res.Error
			}

			if res.RowsAffected != 1 {
				err := kerrors.NewBizStatusError(70006, "Database Error")
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	//double check if likes are negative, if yes , roll back
	video := new(Video)
	if err := DB.WithContext(ctx).First(video, vid).Error; err != nil {
		return err
	}
	if video.FavoriteCount < 0 {
		//TODO: ADD RETRY Method
		err = Favorite(ctx, uid, vid)
		if err != nil {
			logger.Errorf("video %d has negative favorite count, rollback error")
		}
		return kerrors.NewBizStatusError(70004, "Invalid Action")
	}

	return nil
}

// FavoriteList returns a list of Favorite videos.
func FavoriteList(ctx context.Context, uid int64) ([]Video, error) {
	user := new(User)
	if err := DB.WithContext(ctx).First(user, uid).Error; err != nil {
		return nil, err
	}

	videos := []Video{}

	if err := DB.WithContext(ctx).Model(&user).Association("FavoriteVideos").Find(&videos); err != nil {
		return nil, err
	}
	return videos, nil
}
