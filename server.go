package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofor-little/env"
	"gorm.io/driver/mysql"
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
  }

type Sequence struct {
    gorm.Model
    SeqID uint `gorm:"primaryKey"`
    Name string
    Notation string
    UserID uint `gorm:"index"`
}

var (
    db *gorm.DB
    db_user string
    db_pass string
    db_host string
    db_name string
)

func getEnv() {
    if err := env.Load("./env"); err != nil {
		panic(err)
	}
    db_user = env.Get("DB_USER", "root")
    db_pass = env.Get("DB_PASS", "password")
    db_host = env.Get("DB_HOST", "localhost")
    db_name = env.Get("DB_NAME", "test")
}

func initDb() {
    var err error
    dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", db_user, db_pass, db_host, db_name)
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

    if err != nil {
        panic("failed to connect database")
    }

    err = db.AutoMigrate(&User{}, &Sequence{})
    if err != nil {
            fmt.Println("AutoMigrate error: ", err)
    }

    fmt.Printf("Successfully connected to the database\n")
}

func saveSequenceHandler(c *gin.Context) {
    var sequence Sequence
    err := c.BindJSON(&sequence)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    if db != nil {
        result := db.Create(&sequence)
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": result.Error.Error(),
            })
            return
        }
    } else {
        fmt.Println("db is nil")
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Sequence saved successfully",
    })
}

func main() {
    initDb()
    getEnv()
    // insertSuperAdmin()

    router := gin.Default()
    router.POST(("/sequence/save/"), saveSequenceHandler)
    // router.POST(("/user/register"), registerUserHandler)
    // router.POST(("/user/login"), loginUserHandler)

    router.Run(":8080")
}
