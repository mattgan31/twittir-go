package database

import (
	"fmt"
	"log"
	"twittir-go/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "alam"
	dbport   = "5432"
	dbname   = "twittir"
	db       *gorm.DB
	err      error
)

func StartDB() {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, dbport)
	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

}

func GetDB() *gorm.DB {
	return db
}

func Migrate() {
	db.AutoMigrate(&domain.User{}, &domain.Post{}, &domain.Comment{}, &domain.Likes{}, &domain.Relationship{})
}

func ConnectDB() error {
	var err error

	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, dbport)
	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	// Ping the database to ensure the connection is valid
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	log.Println("Database connection established")
	return nil
}
