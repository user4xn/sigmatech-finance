package seeder

import (
	"clean-arch/database"
	"clean-arch/internal/model"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed() {
	var (
		err error
	)

	db := database.GetConnection()

	fmt.Println("executing user seed...")
	err = UserSeed(db)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func UserSeed(db *gorm.DB) error {
	var (
		InsertModel []model.User
	)

	passwordAdmin := []byte("demouser123")
	hashedPasswordAdmin, err := bcrypt.GenerateFromPassword(passwordAdmin, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return err
	}

	InsertModel = []model.User{
		{
			UserType: "admin",
			Email:    "demouser@gmail.com",
			Password: string(hashedPasswordAdmin),
		},
	}

	for _, data := range InsertModel {
		model := model.User{}
		err := db.Where("email = ?", data.Email).First(&model).Error
		if err != nil && err == gorm.ErrRecordNotFound {
			if err := db.Create(&data).Error; err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	return nil
}
