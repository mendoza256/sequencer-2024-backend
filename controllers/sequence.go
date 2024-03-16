package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mendoza256/sequencer-2024-backend/db"
	"github.com/mendoza256/sequencer-2024-backend/models"
)


func GetUserSequencesHandler(c *gin.Context) {
	var sequences []models.Sequence
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if db.DB != nil {
		result := db.DB.Where("user_id = ?", user.ID).Find(&sequences)
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


func SaveSequenceHandler(c *gin.Context) {
    var sequence models.Sequence

	// check if user has more than 10 sequences
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if db.DB != nil {
		var sequences []models.Sequence
		result := db.DB.Where("user_id = ?", user.ID).Find(&sequences)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
			})
			return
		}
		if len(sequences) >= 10 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User has reached the maximum number of sequences",
			})
			return
		}
	} else {
		fmt.Println("db is nil")
	}


	err = c.BindJSON(&sequence)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    if db.DB != nil {
        result := db.DB.Create(&sequence)
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