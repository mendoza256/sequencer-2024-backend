package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mendoza256/sequencer-2024-backend/controllers"
	"github.com/mendoza256/sequencer-2024-backend/db"
)


func main() {
    db.InitDb()

    router := gin.Default()
    router.POST("/save-sequence", controllers.SaveSequenceHandler)
    router.GET("/get-user-sequences", controllers.GetUserSequencesHandler)
    router.POST("/auth/signup", controllers.SignUpHandler)
    router.POST("/auth/login", controllers.LoginHandler)

    router.Run(":8080")
}
