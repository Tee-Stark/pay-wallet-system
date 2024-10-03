package config

import (
	"fmt"
	"pay-system/domain"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func NewPostgresDB() *gorm.DB {
	host := PG_HOST
	port := PG_PORT
	dbUser := PG_USER
	password := PG_PASSWORD
	dbName := PG_DB

	conn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host,
		port,
		dbName,
		dbUser,
		password,
	)

	db, err := gorm.Open("postgres", conn)
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&domain.User{}, &domain.Payment{}, &domain.Wallet{})

	return db
}
