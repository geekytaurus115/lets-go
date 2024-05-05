package handler

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func baseHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Welcome to ScaleX eBook Library")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	// parsing user payload
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
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
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// generate token
	token, err := generateToken(foundUser.Username, foundUser.UserType)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})

	w.WriteHeader(http.StatusOK)
	log.Println("Successfully logged in..")
	fmt.Fprintf(w, "Successfully logged in..")

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	userType := r.Context().Value("user_type").(string)
	log.Println("userType--> ", userType)

	var books []string
	var err error

	// as per user-type show list of books
	switch strings.ToLower(userType) {
	case "admin":
		books, err = readAllBooksFromCSV("db/regularUser.csv", "db/adminUser.csv")
	case "regular":
		books, err = readBooksFromCSV("db/regularUser.csv")
	default:
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// list of books as API response
	for _, book := range books {
		w.Write([]byte(book + "\n"))
	}
}

func addBookHandler(w http.ResponseWriter, r *http.Request) {
	userType := r.Context().Value("user_type").(string)
	log.Println("userType--> ", userType)

	if strings.ToLower(userType) != "admin" {
		http.Error(w, "User is not authorized for this action!", http.StatusForbidden)
		return
	}

	// parsing payload
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// extract book details from form data
	bookName := r.Form.Get("book_name")
	author := r.Form.Get("author")
	publicationYearStr := r.Form.Get("publication_year")

	log.Println("\nbookName: ", bookName,
		"\nauthor: ", author,
		"\npublicationYear: ", publicationYearStr)

	if bookName == "" || author == "" || publicationYearStr == "" {
		http.Error(w, "Book Name, Author, and Publication Year are required", http.StatusBadRequest)
		return
	}

	// some common validations for publication year
	publicationYear, err := strconv.Atoi(publicationYearStr)
	if err != nil || publicationYear < 0 {
		http.Error(w, "Invalid Publication Year", http.StatusBadRequest)
		return
	}

	currentYear := time.Now().Year()
	if publicationYear > currentYear {
		http.Error(w, "Publication Year cannot be in the future", http.StatusBadRequest)
		return
	}

	// now adding the book to the file `regularUser.csv`
	file, err := os.OpenFile("db/regularUser.csv", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("error while opening file --> ", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("error while getting file info --> ", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// need to handle this case, as file has already some data
	// to ensure newly append data from next
	if fileInfo.Size() > 0 {
		_, err = file.WriteString("\n")
		if err != nil {
			log.Println("error while adding newline --> ", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{bookName, author, publicationYearStr})
	if err != nil {
		log.Println("error while adding book --> ", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Book added successfully")
}

func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	userType := r.Context().Value("user_type").(string)
	log.Println("userType--> ", userType)

	if strings.ToLower(userType) != "admin" {
		http.Error(w, "User is not authorized for this action!", http.StatusForbidden)
		return
	}

	// parsing payload
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// extract params from form data
	bookName := r.Form.Get("book_name")

	if bookName == "" {
		http.Error(w, "Book Name is required", http.StatusBadRequest)
		return
	}

	file, err := os.OpenFile("db/regularUser.csv", os.O_RDWR, 0644)
	if err != nil {
		log.Println("Error opening file", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// read CSV records
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("Error reading file's records", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// remove the book from record
	var updatedRecords [][]string
	for _, record := range records {
		// skip if book name matches
		if strings.EqualFold(record[0], bookName) {
			continue
		}
		updatedRecords = append(updatedRecords, record)
	}

	// clear the file content
	file.Truncate(0)
	// reset file offset
	//(0, 0) sets the file pointer's offset to the beginning of the file
	file.Seek(0, 0)

	// write updated records back to the file
	writer := csv.NewWriter(file)
	err = writer.WriteAll(updatedRecords)
	if err != nil {
		log.Println("Error writing updated records", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Book deleted successfully")
}
