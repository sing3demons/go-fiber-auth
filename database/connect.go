package database

import (
	"app/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
	connection, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	connection.AutoMigrate(&models.User{}, &models.PasswordReset{})
	// connection.Migrator().DropTable(&models.User{})
	db = connection

}

func GetDB() *gorm.DB {
	return db
}
