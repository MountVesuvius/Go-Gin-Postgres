// Will swap to using the controller package once this works
package routes

import (
	"github.com/MountVesuvius/go-gin-postgres-template/controllers"
	"github.com/MountVesuvius/go-gin-postgres-template/middleware"
	"github.com/MountVesuvius/go-gin-postgres-template/models"
	"github.com/MountVesuvius/go-gin-postgres-template/services"
	"github.com/gin-gonic/gin"
)

// User route controller. All requests will come in through this and be send to the appropriate user controller
func User(router *gin.Engine, controller controllers.UserController, jwtService services.JWTService) {
	routes := router.Group("/api/v1/user/") 

    routes.POST("/register", controller.Register)
    routes.POST("/login", controller.Login)

    routes.GET(
        "",
        middleware.Authenticate(jwtService),
        middleware.RouterGuard(models.UserRoleGeneral, models.UserRoleAdmin),
        controller.GetUserById,
    )
    routes.GET(
        "/admin",
        middleware.Authenticate(jwtService),
        middleware.RouterGuard(models.UserRoleAdmin),
        controller.Admin,
    )

    // update user role - admin only
    // delete user - admin only
    // update user's name - user and admin
        // user via me, admin via id
        // possibly two routes

}
