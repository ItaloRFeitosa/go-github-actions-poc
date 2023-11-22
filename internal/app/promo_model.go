package app

import "time"

type PromoModel struct {
	ID          uint      `json:"id" uri:"id"`
	Link        string    `json:"link"`
	ProductName string    `json:"productName"`
	CreatedAt   time.Time `json:"createdAt" gorm:"<-:create"`
}

func (p *PromoModel) TableName() string {
	return "promos"
}
