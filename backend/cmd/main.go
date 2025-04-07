package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

var db *sql.DB // Додано глобальну змінну

func indexView(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
}

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.GET("/", indexView)

	// Додайте ваші маршрути для роботи з БД тут

	return router
}

func SetupPostgres() {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")

	// Перевірка змінних оточення
	if username == "" || password == "" || host == "" || dbname == "" {
		log.Fatal("Environment variables not set properly")
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, username, password, dbname)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Connected to PostgreSQL")
}

func main() {
	SetupPostgres()
	defer db.Close() // Закриття з'єднання при завершенні

	router := SetupRoutes()
	router.Run(":8081")
}


