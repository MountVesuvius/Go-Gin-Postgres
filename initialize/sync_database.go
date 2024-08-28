package initialize

import "github.com/MountVesuvius/go-gin-postgres-template/models"

// Auto migrate the table
func SyncDatabase() {
    err := DB.AutoMigrate(&models.User{})
    if err != nil {
        panic("Failed to sync the database")
    }
}
