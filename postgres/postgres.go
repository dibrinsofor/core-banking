package postgres

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConstructDatabaseURI() string {
	dburl := os.Getenv("DATABASE_URL")
	if dburl != "" {
		return dburl
	}

	USER := os.Getenv("POSTGRES_USER")
	PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	HOST := os.Getenv("POSTGRES_HOST")
	DBNAME := os.Getenv("POSTGRES_DB")
	PORT := os.Getenv("POSTGRES_PORT")
	SSLMODE := os.Getenv("SSLMODE")
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", USER, PASSWORD, HOST, PORT, DBNAME, SSLMODE)
}

func Init() (*gorm.DB, error) {

	databaseUrl := ConstructDatabaseURI()

	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
