package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Name string
    Id int
    Email string
    Sequences []Sequence
  }

type Sequence struct {
    gorm.Model
    Name string
    Id int
    Notation [8][16]string
    User User
}


func initDb() {
    dsn := "root:couk-iw-slih-GHAI-maff@tcp(127.0.0.1:3306)/sequencer-2024?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

    if err != nil {
      panic("failed to connect database")
    }

    fmt.Printf("Successfully connected to the database\n")
    fmt.Printf("Database memory address: %p", db)
}

func testHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "Hello, World!",
    })
}

// save a sequence to the database
func saveSequenceHandler(c *gin.Context) {
    var sequence Sequence
    err := c.BindJSON(&sequence)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    fmt.Println(sequence)
    c.JSON(http.StatusOK, gin.H{
        "message": "Sequence saved successfully",
    })
}




func main() {
    initDb()
    router := gin.Default()
    router.GET("/sequence/test/", testHandler)
    router.POST(("/sequence/save/"), saveSequenceHandler)

    router.Run(":8080")
}
