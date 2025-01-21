package model

type Transaction struct {
	ID                 int64   `json:"id" gorm:"primaryKey"`
	ConsumerID         int64   `json:"consumer_id"`
	ContractNo         string  `json:"contract_no"`
	CTR                float64 `json:"ctr"`
	AdminFee           float64 `json:"admin_fee"`
	InstallmentMonth   int8    `json:"installment_month"`
	InterestPercentage float64 `json:"interest_percentage"`
	AssetName          string  `json:"asset_name"`
	Status             string  `json:"status" gorm:"type:enum('pending','process','completed','denied');default:'pending'"`
	Common
	Consumer     Consumer      `json:"consumer,omitempty" gorm:"foreignKey:ConsumerID"`
	Installments []Installment `json:"installments,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}
