package main

import (
	"encoding/gob"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mendoza256/sequencer-2024-backend/controllers"
	"github.com/mendoza256/sequencer-2024-backend/db"
	"github.com/mendoza256/sequencer-2024-backend/middleware"
)


func main() {
    db.InitDb()
    if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

    router := gin.Default()

    // To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})
    store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

    auth, err := controllers.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}
    
	router.GET("/login", controllers.LoginHandler(auth))
	router.GET("/callback", controllers.LoginHandler(auth))
	router.GET("/logout", controllers.LogoutHandler)
    router.GET("/user", middleware.IsAuthenticated, controllers.UserHandler)

    router.Run(":8080")
}
