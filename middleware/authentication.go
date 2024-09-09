package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/MountVesuvius/go-gin-postgres-template/helpers"
	"github.com/MountVesuvius/go-gin-postgres-template/services"
	"github.com/gin-gonic/gin"
)

func Authenticate(jwtService services.JWTService) gin.HandlerFunc {
    return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
        // Using AbortWithStatusJSON so that all further handlers won't run, killing the request chain instantly
        // Authentic requests will come in with headers attempting to authenticate
		if authHeader == "" {
            response := helpers.BuildFailedResponse("Token is missing", nil, nil)
            context.AbortWithStatusJSON(http.StatusUnauthorized, response)
            return
		}
        // Bearer token missing
		if !strings.Contains(authHeader, "Bearer ") {
			response := helpers.BuildFailedResponse("Bearer token is missing", nil, nil)
			context.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
        // Strip out just the JWT token & check that it's actually a jwt
		authHeader = strings.Replace(authHeader, "Bearer ", "", -1)
        token, err := jwtService.ValidateToken(authHeader)
        if err != nil {
            response := helpers.BuildFailedResponse("Invalid Token presented", err, authHeader)
            context.AbortWithStatusJSON(http.StatusUnauthorized, response)
        }
        // Is the token valid?
        if !token.Valid {
            response := helpers.BuildFailedResponse("Token is no longer valid", nil, token)
            context.AbortWithStatusJSON(http.StatusUnauthorized, response)
        }

        // These bottom two are somewhat unlikely....
        claims, err := jwtService.GetTokenClaims(token)
        if err != nil {
            response := helpers.BuildFailedResponse("Claims missing from Token", nil, token)
            context.AbortWithStatusJSON(http.StatusUnauthorized, response)
        }

        userRole, ok := claims["role"].(string)
        if !ok {
            response := helpers.BuildFailedResponse("User Role missing from Token", nil, token)
            context.AbortWithStatusJSON(http.StatusUnauthorized, response)
        }

        log.Println(userRole)

        // Router Guard needs to know the role of the user
        context.Set("userRole", userRole)
        context.Next()
    }
}
