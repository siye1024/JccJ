package db

import (
	"context"
	"gorm.io/gorm"
	"time"
)

// Video Gorm Data Structures
// Video 属于 Author, AuthorID是外键(belongs to)
// caution: UpdatedAt will be covered by model variable 'updated_at', if you want to use create time ,please use `created_at`
type Video struct {
	gorm.Model
	//UpdatedAt     time.Time `gorm:"column:update_time;not null;index:idx_update" `
	Author        User   `gorm:"foreignkey:AuthorID"`
	AuthorID      int64  `gorm:"index:idx_authorid;not null"`
	PlayUrl       string `gorm:"type:varchar(255);not null"`
	CoverUrl      string `gorm:"type:varchar(255)"`
	FavoriteCount int64  `gorm:"default:0"`
	CommentCount  int64  `gorm:"default:0"`
	Title         string `gorm:"type:varchar(50);not null"`
}

func (Video) TableName() string {
	return "videos"
}
func MGetVideos(ctx context.Context, limit int, latestTime *int64) ([]*Video, error) {
	videos := make([]*Video, 0)

	conn := DB.WithContext(ctx)

	if err := conn.Limit(limit).Order("created_at desc").Find(&videos, "created_at < ?", time.UnixMilli(*latestTime)).Error; err != nil {
		return nil, err
	}
	return videos, nil
}
