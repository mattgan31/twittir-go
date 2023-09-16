package database

import (
	"fmt"
	"log"
	"os"
	"twittir-go/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = os.Getenv("PGHOST")
	user     = os.Getenv("PGUSER")
	password = os.Getenv("PGPASSWORD")
	dbport   = os.Getenv("PGPORT")
	dbname   = os.Getenv("PGDATABASE")
	db       *gorm.DB
	err      error
)

func StartDB() {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, dbport)
	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	db.Debug().AutoMigrate(models.User{})
	db.Debug().AutoMigrate(models.Relationship{})
	db.Debug().AutoMigrate(models.Post{})
	db.Debug().AutoMigrate(models.Comment{})
	db.Debug().AutoMigrate(models.Likes{})
}

func GetDB() *gorm.DB {
	return db
}
