package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofor-little/env"
	"golang.org/x/crypto/bcrypt"
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
    Password string
  }

type Sequence struct {
    gorm.Model
    SeqID uint `gorm:"primaryKey"`
    Name string
    Notation string
    UserID uint `gorm:"index"`
}

var (
    DB *gorm.DB
    db_user string
    db_pass string
    db_host string
    db_name string
)

func getEnv() {
    if err := env.Load(".env"); err != nil {
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
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

    if err != nil {
        panic("failed to connect database")
    }

    err = DB.AutoMigrate(&User{}, &Sequence{})
    if err != nil {
            fmt.Println("AutoMigrate error: ", err)
    }

    fmt.Printf("Successfully connected to the database\n")
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func SaveSequenceHandler(c *gin.Context) {
    var sequence Sequence
    err := c.BindJSON(&sequence)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    if DB != nil {
        result := DB.Create(&sequence)
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


func GetUserSequencesHandler(c *gin.Context) {
	var sequences []Sequence
	var user User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if DB != nil {
		result := DB.Where("user_id = ?", user.ID).Find(&sequences)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
			})
			return
		}
	} else {
		fmt.Println("db is nil")
	}

	c.JSON(http.StatusOK, sequences)
}

func SignUpHandler(c *gin.Context) {
    var user User
    err := c.BindJSON(&user)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    hashedPassword, err := HashPassword(user.Password)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    user.Password = hashedPassword

    if DB != nil {
        result := DB.Create(&user)
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
        "message": "User saved successfully",
    })
}

func loginHandler(c *gin.Context) {
    var user User
    err := c.BindJSON(&user)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    var dbUser User
    if DB != nil {
        result := DB.Where("email = ?", user.Email).First(&dbUser)
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": result.Error.Error(),
            })
            return
        }
    } else {
        fmt.Println("db is nil")
    }

    if CheckPasswordHash(user.Password, dbUser.Password) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Login successful",
        })
    } else {
        c.JSON(http.StatusUnauthorized, gin.H{
            "message": "Invalid credentials",
        })
    }
}


func main() {
    getEnv()
    initDb()

    router := gin.Default()
    router.POST("/save-sequence", SaveSequenceHandler)
    router.GET("/get-user-sequences", GetUserSequencesHandler)
    router.POST("/auth/signup", SignUpHandler)
    router.POST("/auth/login", loginHandler)

    router.Run(":8080")
}
