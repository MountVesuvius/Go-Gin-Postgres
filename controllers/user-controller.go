package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MountVesuvius/go-gin-postgres-template/dto"
	"github.com/MountVesuvius/go-gin-postgres-template/initialize"
	"github.com/MountVesuvius/go-gin-postgres-template/models"
	"github.com/MountVesuvius/go-gin-postgres-template/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type (
    UserController interface {
        Signup(context *gin.Context)
        Login (context *gin.Context)
        Validate (context *gin.Context) // temp, remember to delete
    }

    userController struct {
        jwtService services.JWTService
    }
)

func NewUserController(js services.JWTService) UserController {
    return &userController{
        jwtService: js,
    }
}

// temp route to do some quick testing
func (u *userController) Validate (context *gin.Context) {
    tokenString, err := u.jwtService.GenerateAccessToken("username", "role")
    if err != nil {
        fmt.Println("Error creating JWT:", err)
        return
    }

    token, err := u.jwtService.ValidateToken(tokenString)
    if err != nil {
        fmt.Println("Error validating JWT:", err)
        return
    }
    // controller code
    context.JSON(http.StatusOK, gin.H {
        "string": tokenString,
        "parsed": token,
    })
}

/* Example Payload:
{
    "email": "test@email.com",
    "password": "testpassword"
}
*/
func (u *userController) Signup (context *gin.Context) {
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
}

func (u *userController) Login (context *gin.Context) {
    var body dto.Body

    // Ensure payload is correctly structured 
    if context.Bind(&body) != nil {
        context.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to read body",
        })
        return 
    }

    // this should be a service itself
    var user models.User

    // Find the user
    initialize.DB.First(&user, "email = ?", body.Email)

    if user.ID == 0 {
        context.JSON(http.StatusUnauthorized, gin.H {
            "error": "Invalid email or password",
        })
        return
    }

    // validate the password
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
    if err != nil {
        context.JSON(http.StatusUnauthorized, gin.H {
            "error": "Invalid email or password",
        })
        return
    }
    // end of user service

    // Register new jwt claims
    userIdString := strconv.FormatUint(uint64(user.ID), 10)
    accessToken, err := u.jwtService.GenerateAccessToken(userIdString, "user")
    if err != nil {
        fmt.Println("access token error:", err)
    }
    refreshToken, err := u.jwtService.GenerateRefreshToken()
    if err != nil {
        fmt.Println("refresh token error:", err)
    }
    fmt.Println("Hello these are tokens")
    fmt.Println("accesstoken:", accessToken)
    fmt.Println("refreshtoken:", refreshToken)

    context.JSON(http.StatusOK, gin.H {
        "access": accessToken,
        "refresh": refreshToken,
    })
}
