package routes

import (
	"net/http"

	"github.com/MountVesuvius/go-gin-postgres-template/helpers"
	"github.com/MountVesuvius/go-gin-postgres-template/services"
	"github.com/gin-gonic/gin"
)

type Body struct {
    Token string
}

func Auth(router *gin.Engine, jwtService services.JWTService) {
    routes := router.Group("/api/v1/auth")

    // should send in the refresh token that was generated on login to refresh
    routes.POST("/refresh", func (context *gin.Context) {
        var body Body
        // Get the Post body
        bindErr := context.Bind(&body)
        if bindErr != nil {
            response := helpers.BuildFailedResponse("bind failed", bindErr, nil)
            context.JSON(http.StatusBadRequest, response)
            return
        }

        // Refresh the access token
        newAccessToken, err := jwtService.RefreshToken(body.Token)
        if err != nil {
            response := helpers.BuildFailedResponse("Failed to refresh access token", body, err)
            context.JSON(http.StatusBadRequest, response)
            return
        }

        // No need for a message to be returned here. Just send the token down
        context.JSON(http.StatusOK, gin.H {
            "accessToken": newAccessToken,
        })
    })
}
