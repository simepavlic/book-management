package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Book struct {
	ID            uint   `gorm:"primaryKey"`
	Title         string `json:"title" gorm:"not null;unique"`
	Author        string `json:"author" gorm:"not null"`
	PublishedDate string `json:"published_date" gorm:"not null"`
	Edition       string `json:"edition"`
	Description   string `json:"description"`
	Genre         string `json:"genre"`
}

type Collection struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `json:"name" gorm:"not null"`
	Books []Book `json:"books" gorm:"many2many:book_collection;"`
}

var db *gorm.DB

func main() {
	// Connect to the PostgreSQL database
	db, err := gorm.Open(postgres.Open("host=localhost port=5432 user=postgres password=root dbname=postgres sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db = db.Exec("CREATE DATABASE books_db;")
	if db.Error != nil {
		fmt.Println("Unable to create DB books_db, attempting to connect assuming it exists...")
		db, err = gorm.Open(postgres.Open("host=localhost port=5432 user=postgres password=root dbname=books_db sslmode=disable"), &gorm.Config{})
		if err != nil {
			fmt.Println("Unable to connect to books_db")
			log.Fatal(err)
		}
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the database schema
	db.AutoMigrate(&Book{}, &Collection{})

	// Create a Gin router
	r := gin.Default()

	// Define the routes
	r.GET("/books", listBooks)
	r.POST("/books", createBook)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func listBooks(c *gin.Context) {
	var books []Book

	// Query the books from the database
	if err := db.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the books as JSON
	c.JSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var book Book

	// Bind the request body to a Book struct
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the book in the database
	if err := db.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the created book as JSON
	c.JSON(http.StatusCreated, book)
}
