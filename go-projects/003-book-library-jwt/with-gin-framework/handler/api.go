package handler

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func baseHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"error":   false,
		"message": "checking server..",
	})
}

func loginHandler(c *gin.Context) {
	var user User

	// parsing user payload
	if err := c.BindJSON(&user); err != nil {
		log.Println("Error parsing user payload: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var foundUser *User
	for _, u := range UserData {
		if u.Username == user.Username && u.Password == user.Password {
			foundUser = &u
			break
		}
	}

	if foundUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// generate token
	token, err := generateToken(foundUser.Username, foundUser.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{"token": token})

	log.Println("Successfully logged in..")
	fmt.Fprintf(c.Writer, "Successfully logged in..")
}

func homeHandler(c *gin.Context) {
	userType, _ := c.Get("user_type")
	log.Println("userType--> ", userType)

	var books []string
	var err error

	// as per user-type, show list of books
	switch strings.ToLower(userType.(string)) {
	case "admin":
		books, err = readAllBooksFromCSV("db/regularUser.csv", "db/adminUser.csv")
	case "regular":
		books, err = readBooksFromCSV("db/regularUser.csv")
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// list of books as API response
	c.String(http.StatusOK, strings.Join(books, "\n"))
}

func addBookHandler(c *gin.Context) {
	userType, _ := c.Get("user_type")
	log.Println("userType--> ", userType)

	if strings.ToLower(userType.(string)) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "User is not authorized for this action!"})
		return
	}

	// parsing payload
	var req struct {
		BookName        string `form:"book_name" binding:"required"`
		Author          string `form:"author" binding:"required"`
		PublicationYear int    `form:"publication_year" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// some common validations for publication year
	currentYear := time.Now().Year()
	if req.PublicationYear < 0 || req.PublicationYear > currentYear {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Publication Year"})
		return
	}

	// now adding the book to the file `regularUser.csv`
	file, err := os.OpenFile("db/regularUser.csv", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error while opening file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write book details to the CSV file
	if err := writer.Write([]string{req.BookName, req.Author, strconv.Itoa(req.PublicationYear)}); err != nil {
		log.Println("Error while adding book:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Book added successfully"})
}

func deleteBookHandler(c *gin.Context) {
	userType, _ := c.Get("user_type")
	log.Println("userType--> ", userType)

	if strings.ToLower(userType.(string)) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "User is not authorized for this action!"})
		return
	}

	// parsing payload
	var req struct {
		BookName string `form:"book_name" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := os.OpenFile("db/regularUser.csv", os.O_RDWR, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer file.Close()

	// read CSV records
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("Error reading file's records:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// remove the book from records
	var updatedRecords [][]string
	for _, record := range records {
		// skip if book name matches
		if strings.EqualFold(record[0], req.BookName) {
			continue
		}
		updatedRecords = append(updatedRecords, record)
	}

	// clear the file content
	file.Truncate(0)
	// reset file offset to beginning i.e. (0,0)
	file.Seek(0, 0)

	// write updated records back to the file
	writer := csv.NewWriter(file)
	if err := writer.WriteAll(updatedRecords); err != nil {
		log.Println("Error writing updated records:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
