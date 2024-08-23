// Will swap to using the controller package once this works
package routes

import (
	"github.com/MountVesuvius/go-gin-postgres-template/controllers"
	"github.com/MountVesuvius/go-gin-postgres-template/middleware"
	"github.com/MountVesuvius/go-gin-postgres-template/services"
	"github.com/gin-gonic/gin"
)

// User route controller. All requests will come in through this and be send to the appropriate user controller
func User(router *gin.Engine, controller controllers.UserController, jwtService services.JWTService) {
	routes := router.Group("/api/user") 

    routes.POST("/signup", controller.dignup)
    routes.POST("/login", controller.Login)

    // temp route remember to delete
    routes.GET("/validate", middleware.Authenticate(jwtService), controller.Validate)
}

