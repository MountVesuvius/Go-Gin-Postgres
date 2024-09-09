package initialize

import (
	"fmt"
	"os"

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

// Function of mass destruction 
// func dropAllTables() error {
//     var tables []string
//     log.Println("Dropped all tables")
// 	DB.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables)
//
// 	for _, table := range tables {
// 		err := DB.Migrator().DropTable(table)
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	return nil
// }
