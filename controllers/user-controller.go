package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MountVesuvius/go-gin-postgres-template/dto"
	"github.com/MountVesuvius/go-gin-postgres-template/helpers"
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

    genericResponse := helpers.BuildFailedResponse("Unexpected error occured. Please try again later", nil, nil)

    // Ensure payload is correctly structured 
    bindErr := context.Bind(&body)
    if bindErr != nil {
        response := helpers.BuildFailedResponse("Failed to read body", bindErr, nil)
        context.JSON(http.StatusBadRequest, response)
        return 
    }

    // Will come back to this and change it up
    // https://labs.clio.com/bcrypt-cost-factor-4ca0a9b03966
    // https://stackoverflow.com/questions/4443476/optimal-bcrypt-work-factor/61304956#61304956
    hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
    if err != nil {
        fmt.Println("bcrypt error when hashing password", err)
        context.JSON(http.StatusInternalServerError, genericResponse)
        return
    }

    // Add user to database
    user := models.User{ Email: body.Email, Password: string(hash) }
    result := initialize.DB.Create(&user)
    if result.Error != nil {
        fmt.Println("error creating the user", result.Error)
        context.JSON(http.StatusInternalServerError, genericResponse)
        return
    }

    response := helpers.BuildSuccessfulResponse("Sucessfully created user", user)
    context.JSON(http.StatusOK, response)
}

func (u *userController) Login (context *gin.Context) {
    var body dto.Body

    // Ensure payload is correctly structured 
    bindErr := context.Bind(&body)
    if bindErr != nil {
        response := helpers.BuildFailedResponse("Failed to read body", bindErr, nil)
        context.JSON(http.StatusBadRequest, response)
        return 
    }

    // this should be a service itself
    // ----------
    var user models.User

    // Find the user
    initialize.DB.First(&user, "email = ?", body.Email)

    if user.ID == 0 {
        response := helpers.BuildFailedResponse("Invalid Email or Password", nil, nil)
        context.JSON(http.StatusUnauthorized, response)
        return
    }

    // Validate the password
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
    if err != nil {
        response := helpers.BuildFailedResponse("Invalid Email or Password", nil, nil)
        context.JSON(http.StatusUnauthorized, response)
        return
    }
    // end of user service
    // ----------

    // Register new jwt claims
    userIdString := strconv.FormatUint(uint64(user.ID), 10)
    accessToken, err := u.jwtService.GenerateAccessToken(userIdString, "user")
    if err != nil {
        fmt.Println("Access token error:", err)
    }
    refreshToken, err := u.jwtService.GenerateRefreshToken(userIdString, "user")
    if err != nil {
        fmt.Println("refresh token error:", err)
    }

    // Custom response to save on space
    context.JSON(http.StatusOK, gin.H {
        "access": accessToken,
        "refresh": refreshToken,
    })
}
