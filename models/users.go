package models

import (
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Email string `gorm:"unique"` // will add an ID, Created, Updated, Deleted column
    Password string
}
