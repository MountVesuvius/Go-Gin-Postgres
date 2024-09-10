package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MountVesuvius/go-gin-postgres-template/dto"
	"github.com/MountVesuvius/go-gin-postgres-template/helpers"
	"github.com/MountVesuvius/go-gin-postgres-template/models"
	"github.com/MountVesuvius/go-gin-postgres-template/services"
)

type (
    UserController interface {
        Register(context *gin.Context)
        Login (context *gin.Context)
        GetUserById(context *gin.Context)
        Admin(context *gin.Context)
    }

    userController struct {
        jwtService services.JWTService
        userService services.UserService
    }
)

func NewUserController(js services.JWTService, us services.UserService) UserController {
    return &userController{
        jwtService: js,
        userService: us,
    }
}

// Register first checks that the payload matches the expected input, then attempts to
// register a new user via the UserService, failing if the user already exists.
func (u *userController) Register(context *gin.Context) {
    var payload dto.AuthenticateUser

    // Ensure payload is correctly structured 
    bindErr := context.Bind(&payload)
    if bindErr != nil {
        response := helpers.BuildFailedResponse("Failed to read body", bindErr, nil)
        context.JSON(http.StatusBadRequest, response)
        return 
    }

    // Register a new user
    _, err := u.userService.Register(payload.Password, payload.Email, models.UserRoleGeneral)
    if err != nil {
        // 99.999% of the time this will fail as a result of the user already existing.
        // There is a tiny chance that bcrypt fails...
        genericResponse := helpers.BuildFailedResponse("User already exists", nil, nil)
        context.JSON(http.StatusInternalServerError, genericResponse)
        return
    }

    response := helpers.BuildSuccessfulResponse("Sucessfully created user", nil)
    context.JSON(http.StatusOK, response)
}

func (u *userController) Login(context *gin.Context) {
    var payload dto.AuthenticateUser

    // 1. Ensure payload is correctly structured 
    bindErr := context.Bind(&payload)
    if bindErr != nil {
        response := helpers.BuildFailedResponse("Failed to read body", bindErr, nil)
        context.JSON(http.StatusBadRequest, response)
        return 
    }

    // 2. User Login
    user, err := u.userService.Login(payload.Password, payload.Email)
    if err != nil {
        response := helpers.BuildFailedResponse("Wrong Username or Password", err, nil)
        context.JSON(http.StatusBadRequest, response)
        return
    }

    // 3. Register new JWT claims
    userIdString := strconv.FormatUint(uint64(user.ID), 10)
    accessToken, err := u.jwtService.GenerateAccessToken(userIdString, user.Role)
    if err != nil {
        log.Println("JWT ERROR OCCURED:", err)
        response := helpers.BuildFailedResponse("Internal Server Error", nil, nil)
        context.JSON(http.StatusInternalServerError, response)
        return
    }
    refreshToken, err := u.jwtService.GenerateRefreshToken(userIdString, user.Role)
    if err != nil {
        log.Println("JWT ERROR OCCURED:", err)
        response := helpers.BuildFailedResponse("Internal Server Error", nil, nil)
        context.JSON(http.StatusInternalServerError, response)
        return
    }

    // Custom response to save on space
    context.JSON(http.StatusOK, gin.H {
        "access": accessToken,
        "refresh": refreshToken,
    })
}

func (u *userController) GetUserById(context *gin.Context) {
    id := context.Query("id")

    user, err := u.userService.GetUserById(id)
    if err != nil {
        response := helpers.BuildFailedResponse("User could not be found", err, nil)
        context.JSON(http.StatusBadRequest, response)
        return 
    }

    // Trying not to give out passwords for free
    displayUser := dto.DisplayUser{
        Role: user.Role,
        Email: user.Email,
        Name: user.Name,
    }

    response := helpers.BuildSuccessfulResponse("User was found", displayUser)
    context.JSON(http.StatusOK, response)
}

func (u *userController) Admin(context *gin.Context) {
    response := helpers.BuildSuccessfulResponse("You are an admin", "smile")
    context.JSON(http.StatusOK, response)
}
