package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/mendoza256/sequencer-2024-backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
    DB *gorm.DB
    db_user string
    db_pass string
    db_host string
    db_name string
)

func getEnv() {
    godotenv.Load(".env")
    db_user = os.Getenv("DB_USER")
    db_pass = os.Getenv("DB_PASS")
    db_host = os.Getenv("DB_HOST")
    db_name = os.Getenv("DB_NAME")
}

func InitDb() {
    getEnv()
    var err error
    dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", db_user, db_pass, db_host, db_name)
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

    if err != nil {
        panic("failed to connect database")
    }

    err = DB.AutoMigrate(&models.User{}, &models.Sequence{})
    if err != nil {
            fmt.Println("AutoMigrate error: ", err)
    }

    fmt.Printf("Successfully connected to the database\n")
}