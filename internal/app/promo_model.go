package app

import (
	"time"

	"gorm.io/gorm"
)

type PromoModel struct {
	ID          uint           `json:"id" uri:"id"`
	Link        string         `json:"link"`
	ProductName string         `json:"productName"`
	CreatedAt   time.Time      `json:"createdAt" gorm:"<-:create"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt"`
}

func (p *PromoModel) TableName() string {
	return "promos"
}
