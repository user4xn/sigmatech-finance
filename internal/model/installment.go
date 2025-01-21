package model

import "time"

type Installment struct {
	ID              int64      `json:"id" gorm:"primaryKey"`
	TransactionID   int64      `json:"transaction_id"`
	ConsumerID      int64      `json:"consumer_id"`
	Label           string     `json:"label"`
	Amount          float64    `json:"amount"`
	PaymentDeadline time.Time  `json:"payment_deadline"`
	PaidAt          *time.Time `json:"paid_at"`
	Common
	Transaction Transaction          `json:"transaction,omitempty" gorm:"foreignKey:TransactionID"`
	Consumer    Consumer             `json:"consumer,omitempty" gorm:"foreignKey:ConsumerID"`
	Penalties   []TransactionPenalty `json:"penalties,omitempty"`
}

func (Installment) TableName() string {
	return "installments"
}
