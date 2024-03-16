package models

import (
	"gorm.io/gorm"
)

type Role string

const (
    SuperAdmin Role = "superadmin"
    Admin      Role = "admin"
    Standard   Role = "standard"
)

type User struct {
    gorm.Model
    ID uint `gorm:"primaryKey"`
    Name string
    Email string
    Sequences []Sequence
    Role Role
    Password string
  }

type Sequence struct {
    gorm.Model
    SeqID uint `gorm:"primaryKey"`
    Name string
    Notation string
    UserID uint `gorm:"index"`
}
