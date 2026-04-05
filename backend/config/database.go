package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadDB() {

	log.Printf("log url %s\n", ENV.DBUrl)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		ENV.DBUrl,
		ENV.DBUsername,
		ENV.DBPassword,
		ENV.DBDatabase,
		ENV.DBPort,
	)

	log.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed connect database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed get sqlDB:", err)
	}

	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(4)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	DB = db
}
