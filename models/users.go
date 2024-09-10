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
    Name string
    Email string `gorm:"unique"` // will add an ID, Created, Updated, Deleted column
    Password string
    Role string
}

// wip
type DisplayUser struct {
    gorm.Model
    Name string
    Email string `gorm:"unique"` // will add an ID, Created, Updated, Deleted column
    Role string
}
