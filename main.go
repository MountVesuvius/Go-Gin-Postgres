package main

import (
    "log"

	"github.com/joho/godotenv"
	"github.com/MountVesuvius/go-gin-postgres-template/initialize"
	"github.com/MountVesuvius/go-gin-postgres-template/routes"
	"github.com/gin-gonic/gin"
)

func main() {
    // Load Environment Variables
    err := godotenv.Load()   
    if err != nil {
        log.Fatal("Error loading environment variables: ", err)
    }

    // Setup Database
    initialize.ConnectToDatabase()
    initialize.SyncDatabase()

    // Setup Router
    router := gin.Default()

    routes.Auth(router)

    router.Run()
}
