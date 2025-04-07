package main

import (
	"github.com/vladmorozov2/auction-service/internal/models"
	"fmt"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)
var gormDB *gorm.DB

type ListItem struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Item string `gorm:"not null" json:"item"`
	Done bool   `gorm:"default:false" json:"done"`
}

func main() {
	var err error
	gormDB, err = connectToPostgreSQL()
	var auction models.Auction
  	fmt.Println(auction)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	router := SetupRoutes()
	router.Run(":8081")
}

func connectToPostgreSQL() (*gorm.DB, error) {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		host, username, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&ListItem{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %v", err)
	}

	log.Println("Connected to PostgreSQL")
	return db, nil
}

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/", indexView)
	router.POST("/todo", CreateTodoItem)

	return router
}

func indexView(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
}

func CreateTodoItem(c *gin.Context) {
	var todo ListItem
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := gormDB.Create(&todo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}
