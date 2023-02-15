package db

import (
	"context"
	"gorm.io/gorm"
)

// User Gorm Data Structures
// User和Video是多对多关系，两个模型之间会有一个连接表 user_favorite_videos
type User struct {
	gorm.Model
	UserName       string  `gorm:"index:idx_username,unique;type:varchar(40);not null" json:"username"`
	Password       string  `gorm:"type:varchar(256);not null" json:"password"`
	FavoriteVideos []Video `gorm:"many2many:user_favorite_videos" json:"favorite_videos"`
	FollowCount    int     `gorm:"default:0" json:"follow_count"`
	FollowerCount  int     `gorm:"default:0" json:"follower_count"`
}

func (User) TableName() string {
	return "users"
}

// CreateUser create user info
func CreateUser(ctx context.Context, users []*User) error {
	return DB.WithContext(ctx).Create(users).Error
}

func QueryUser(ctx context.Context, userName string) ([]*User, error) {
	res := make([]*User, 0)

	if err := DB.WithContext(ctx).Where("user_name = ?", userName).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
