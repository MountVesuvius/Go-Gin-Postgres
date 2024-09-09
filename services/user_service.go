package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/MountVesuvius/go-gin-postgres-template/initialize"
	"github.com/MountVesuvius/go-gin-postgres-template/models"
)

type (
    UserService interface {
        Register(string, string, string) (models.User, error)
        Login(string, string) (models.User, error)
        GetUserById(string) (models.User, error)
    }

    userService struct {}
)

func NewUserService() UserService {
    return &userService {}
}

// Register first checks the user doesn't already exist, then hashes
// the passed password with Bcrypt, then fills the User model which is
// passed to the Database.
func (u *userService) Register(password string, email string, role string) (models.User, error) {
    var user models.User

    // 1. Check the user doesn't already exist
    initialize.DB.First(&user, "email = ?", email)
    if user.ID != 0 {
        return user, errors.New("User Already Exists")
    }

    // 2. Hash user password
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 15) // Turn this into .env var??
    if err != nil {
        return models.User{}, err
    }

    // 3. Add user to database
    user = models.User{ Email: email, Password: string(hash), Role: role }
    result := initialize.DB.Create(&user)
    if result.Error != nil {
        return models.User{}, err
    }

    return user, nil
}

// Login first searches for the specified user via the passed email, then compares
// the passed hash with the stored hash in the User table.
func (u *userService) Login(password string, email string) (models.User, error) {
    var user models.User

    // 1. Find the user
    initialize.DB.First(&user, "email = ?", email)
    if user.ID == 0 {
        return user, errors.New("Failed to find user")
    }

    // 2. Validate the password
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return user, err
    }

    return user, nil
}

// GetUserById searches for the first user by the passed userId.
func (u *userService) GetUserById(userId string) (models.User, error) {
    var user models.User

    if err := initialize.DB.First(&user, "id = ?", userId).Error; err != nil {
        return user, errors.New("User cannot be found")
    }

    return user, nil
}
