package model

type Limit struct {
	ID         int64   `json:"id" gorm:"primaryKey"`
	ConsumerID int64   `json:"consumer_id"`
	Amount     float64 `json:"amount"`
	Tenor      int8    `json:"tenor"`
	Common
	Consumer Consumer `json:"consumer,omitempty" gorm:"foreignKey:ConsumerID"`
}

func (Limit) TableName() string {
	return "limits"
}
