// Will swap to using the controller package once this works
package routes

import (
	"github.com/MountVesuvius/go-gin-postgres-template/controllers"
	"github.com/gin-gonic/gin"
)

func User(router *gin.Engine, controller controllers.UserController) {
	routes := router.Group("/api/auth") 

    routes.POST("/signup", controller.Signup)
    routes.POST("/login", controller.Login)

    // temp route remember to delete
    routes.GET("/validate", controller.Validate)
}

