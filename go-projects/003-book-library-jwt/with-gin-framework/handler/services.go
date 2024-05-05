package handler

import (
	"encoding/csv"
	"os"
)

func readBooksFromCSV(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// read the file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// extracting book names
	var books []string
	for i, record := range records {
		// just to skip the header row
		if i == 0 {
			continue
		}
		books = append(books, record[0])
	}

	return books, nil
}

func readAllBooksFromCSV(files ...string) ([]string, error) {
	var allBooks []string

	for _, file := range files {
		books, err := readBooksFromCSV(file)
		if err != nil {
			return nil, err
		}
		allBooks = append(allBooks, books...)
	}

	return allBooks, nil
}
