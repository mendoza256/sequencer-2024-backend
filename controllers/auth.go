package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mendoza256/sequencer-2024-backend/db"
	"github.com/mendoza256/sequencer-2024-backend/lib"
	"github.com/mendoza256/sequencer-2024-backend/models"
)


func LoginHandler(c *gin.Context) {
    var user models.User
    err := c.BindJSON(&user)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    var dbUser models.User
    if db.DB != nil {
        result := db.DB.Where("email = ?", user.Email).First(&dbUser)
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": result.Error.Error(),
            })
            return
        }
    } else {
        fmt.Println("db is nil")
    }

    if lib.CheckPasswordHash(user.Password, dbUser.Password) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Login successful",
        })
    } else {
        c.JSON(http.StatusUnauthorized, gin.H{
            "message": "Invalid credentials",
        })
    }
}

func SignUpHandler(c *gin.Context) {
    var user models.User
    err := c.BindJSON(&user)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    hashedPassword, err := lib.HashPassword(user.Password)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    user.Password = hashedPassword

    if db.DB != nil {
        result := db.DB.Create(&user)
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