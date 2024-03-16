package main

import (
<<<<<<< Updated upstream
	"encoding/gob"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
=======
	"time"

	"github.com/gin-contrib/cors"
>>>>>>> Stashed changes
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
<<<<<<< Updated upstream

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
=======
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173/"},
        AllowMethods:     []string{"PUT", "PATCH"},
        AllowHeaders:     []string{"Origin"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        AllowOriginFunc: func(origin string) bool {
          return origin == "http://localhost:5173/"
        },
        MaxAge: 12 * time.Hour,
      }))
    router.POST("/save-sequence", controllers.SaveSequenceHandler)
    router.GET("/get-user-sequences", controllers.GetUserSequencesHandler)
>>>>>>> Stashed changes

    router.Run(":8080")
}
