package model

import "time"

type TransactionPenalty struct {
	ID                int64      `json:"id" gorm:"primaryKey"`
	TransactionID     int64      `json:"transaction_id"`
	InstallmentID     int64      `json:"installment_id"`
	PenaltyPercentage float64    `json:"penalty_percentage"`
	PaidAt            *time.Time `json:"paid_at"`
	Common
	Transaction Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID"`
	Installment Installment `json:"installment,omitempty" gorm:"foreignKey:InstallmentID"`
}

func (TransactionPenalty) TableName() string {
	return "transaction_penalties"
}
