package model

type User struct {
	ID            int    `gorm:"primaryKey" json:"id"`
	Email         string `gorm:"column:email" json:"email"`
	Password      string `gorm:"column:password" json:"password"`
	RememberToken string `gorm:"column:remember_token" json:"remember_token"`
	Common
}

func (User) TableName() string {
	return "users"
}
