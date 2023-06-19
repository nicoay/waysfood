package database

import (
	"fmt"
	"mytask/models"
	"mytask/pkg/mysql"
)

func RunMigration() {
	err := mysql.DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Transaction{},
		&models.Carts{},
		&models.Order{},
	)

	if err != nil {
		panic("Migration Failed")
	}

	fmt.Println("Migration Success👍👍")
}
