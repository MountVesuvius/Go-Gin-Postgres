package initialize

import (
	"os"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
    var err error

    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    name := os.Getenv("DB_NAME")
    // This is defined as an environment variable by docker compose, not the .env file
    host := os.Getenv("DB_HOST") 

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, name)
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        panic("Failed to connect to Database")
    }

}
