package initialize

import "github.com/MountVesuvius/go-gin-postgres-template/models"

// Auto migrate the table
func SyncDatabase() {
    DB.AutoMigrate(&models.User{})
}
