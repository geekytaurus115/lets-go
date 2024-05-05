package handler

import "net/http"

func Routes() {

	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/home", authMiddleware(homeHandler))
	http.HandleFunc("/addBook", authMiddleware(addBookHandler))
	http.HandleFunc("/deleteBook", authMiddleware(deleteBookHandler))
}
