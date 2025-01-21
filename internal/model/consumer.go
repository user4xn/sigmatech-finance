package model

import "time"

type Consumer struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	NIK        string    `json:"nik" gorm:"unique"`
	FullName   string    `json:"full_name"`
	LegalName  string    `json:"legal_name"`
	BirthPlace string    `json:"birth_place"`
	BirthDate  time.Time `json:"birth_date"`
	Salary     float64   `json:"salary"`
	KtpURL     string    `json:"ktp_url"`
	SelfieURL  string    `json:"selfie_url"`
	Common
	Limits       []Limit       `json:"limits,omitempty"`
	Transactions []Transaction `json:"transactions,omitempty"`
}

func (Consumer) TableName() string {
	return "consumers"
}
