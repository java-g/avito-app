package models

import (
	"time"

	pq "github.com/lib/pq"
	"gorm.io/datatypes"
)

type Banner struct {
	BannerID  uint           `json:"banner_id" gorm:"primarykey"`
	FeatureID uint           `json:"feature_id"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	TagIds    pq.Int64Array  `json:"tag_ids" gorm:"type:integer[]"`
	Content   datatypes.JSON `json:"content"`
}
