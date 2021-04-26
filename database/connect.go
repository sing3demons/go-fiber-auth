package database

import (
	"app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
	dsn := "user=postgres password=123456 dbname=go_auth port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
