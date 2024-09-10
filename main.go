package main

import (
	"log"

	"github.com/MountVesuvius/go-gin-postgres-template/controllers"
	"github.com/MountVesuvius/go-gin-postgres-template/initialize"
	"github.com/MountVesuvius/go-gin-postgres-template/routes"
	"github.com/MountVesuvius/go-gin-postgres-template/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

    // Setup Services
    jwtService := services.NewJWTService()
    userService := services.NewUserService()

    // Setup Controllers
    userController := controllers.NewUserController(jwtService, userService)

    // Setup Router
    router := gin.Default()
    routes.User(router, userController, jwtService)
    routes.Auth(router, jwtService)

    routerErr := router.Run()
    if routerErr != nil {
        log.Fatal("Router has not started", err)
    }
}
