// Will swap to using the controller package once this works
package routes

import (
	"net/http"

	// "github.com/MountVesuvius/go-gin-postgres-template/controllers"
	"github.com/MountVesuvius/go-gin-postgres-template/dto"
	"github.com/MountVesuvius/go-gin-postgres-template/initialize"
	"github.com/MountVesuvius/go-gin-postgres-template/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// security related errors that could be exploited will be returned as an internal server error
// security by obscurity
func Auth(router *gin.Engine) {
	routes := router.Group("/api/auth") 

    /* Example Payload:
    {
        "email": "test@email.com",
        "password": "testpassword"
    }
    */
    routes.POST("/signup", func(context *gin.Context) {
        var body dto.Body

        // Ensure payload is correctly structured 
        if context.Bind(&body) != nil {
            context.JSON(http.StatusBadRequest, gin.H{
                "error": "Failed to read body",
            })
            return 
        }

        // Will come back to this and change it up
        // https://labs.clio.com/bcrypt-cost-factor-4ca0a9b03966
        // https://stackoverflow.com/questions/4443476/optimal-bcrypt-work-factor/61304956#61304956
        hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
        if err != nil {
            context.JSON(http.StatusInternalServerError, gin.H {
                "error": "Unexpected error occured. Please try again later",
            })
            return
        }

        // Add user to database
        user := models.User{Email: body.Email, Password: string(hash)}
        result := initialize.DB.Create(&user)
        if result.Error != nil {
            context.JSON(http.StatusInternalServerError, gin.H {
                "error": "Unexpected error occured. Please try again later",
            })
            return
        }

        context.JSON(http.StatusOK, gin.H {
            "message": user,
        })

    })

    routes.POST("/login", func(context *gin.Context) {
        var body dto.Body

        // Ensure payload is correctly structured 
        if context.Bind(&body) != nil {
            context.JSON(http.StatusBadRequest, gin.H{
                "error": "Failed to read body",
            })
            return 
        }

        var user models.User

        // don't really like this solution 
        initialize.DB.First(&user, "email = ?", body.Email)

        if user.ID == 0 {
            context.JSON(http.StatusBadRequest, gin.H {
                "error": "Invalid email or password",
            })
            return
        }

        err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
        if err != nil {
            context.JSON(http.StatusBadRequest, gin.H {
                "error": "Invalid email or password",
            })
            return
        }

        // wouldn't do this normally, it's just to see if it's working
        context.JSON(http.StatusOK, gin.H {
            "message": user,
            "status": "logged in",
        })
    })

    routes.GET("/validate", func(context *gin.Context) {
        // controller code
        context.JSON(http.StatusNotImplemented, gin.H {
            "message": "not yet implemented validate",
        })
    })

}
