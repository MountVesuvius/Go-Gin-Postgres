package middleware

import (
	"net/http"

	"github.com/MountVesuvius/go-gin-postgres-template/helpers"
	"github.com/gin-gonic/gin"
)

func RouterGuard(requiredRoles ...string) gin.HandlerFunc {
    return func (context *gin.Context) {
        role, ok := context.Get("userRole")
        if !ok {
            response := helpers.BuildFailedResponse("Failed to access user role", nil, nil)
            context.AbortWithStatusJSON(http.StatusInternalServerError, response)
        }

        userRole, ok := role.(string)
        if !ok {
            response := helpers.BuildFailedResponse("Failed to access user role", nil, nil)
            context.AbortWithStatusJSON(http.StatusInternalServerError, response)
        }

        for _, reqRole := range requiredRoles {
            if userRole == reqRole {
                context.Next()
                return
            }
        }

        response := helpers.BuildFailedResponse("User is unauthorised to access this route", nil, userRole)
        context.AbortWithStatusJSON(http.StatusUnauthorized, response)
    }
}
