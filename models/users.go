package models

import (
	"gorm.io/gorm"
)

const (
    UserRoleAdmin = "Admin"
    UserRoleGeneral = "User"
    UserRoleReadOnly = "ReadOnly"
)

type User struct {
    gorm.Model
    Email string `gorm:"unique"` // will add an ID, Created, Updated, Deleted column
    Password string
    Role string
}
