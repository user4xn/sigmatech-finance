package model

type User struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	UserType   string `json:"user_type" gorm:"type:enum('admin','consumer');default:'consumer'"`
	ConsumerID *int64 `json:"consumer_id"`
	Email      string `gorm:"column:email" json:"email"`
	Password   string `gorm:"column:password" json:"password"`
	Common
	Consumer *Consumer `json:"consumer,omitempty" gorm:"foreignKey:ConsumerID"`
}

func (User) TableName() string {
	return "users"
}
