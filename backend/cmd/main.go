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
	router.GET("/todo/:item", CreateTodoItem)

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

type ListItem struct {
	Id   string `json:"id"`
	Item string `json:"item"`
	Done bool   `json:"done"`
}

func CreateTodoItem(c *gin.Context) {
	item := c.Param("item")

	// Validate item
	if len(item) == 0 {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "please enter an item"})
		return
	}

	// Create todo item struct
	var TodoItem ListItem
	TodoItem.Item = item
	TodoItem.Done = false

	// Insert item into DB and return the generated ID
	err := db.QueryRow("INSERT INTO list(item, done) VALUES($1, $2) RETURNING id;", TodoItem.Item, TodoItem.Done).Scan(&TodoItem.Id)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error with DB"})
		return
	}

	// Log message
	log.Println("created todo item with ID", TodoItem.Id)

	// Return success response with item ID
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	c.JSON(http.StatusCreated, gin.H{"item": TodoItem})
}
