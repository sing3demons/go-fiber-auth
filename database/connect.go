package database

import (
	"app/models"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	name := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(db)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, name)

	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// connection, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	connection.AutoMigrate(&models.User{}, &models.PasswordReset{})
	// connection.Migrator().DropTable(&models.User{}, &models.PasswordReset{})
	db = connection

}

func GetDB() *gorm.DB {
	return db
}
