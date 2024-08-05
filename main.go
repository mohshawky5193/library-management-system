package main

import (
	"library_management_system/db"
	"library_management_system/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()

	router := mux.NewRouter()
	router.Use(handlers.ErrorHandler)

	// Public routes
	router.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", handlers.GenerateJWT).Methods("POST")

	// Create a subrouter for all /books/* routes
	booksRouter := router.PathPrefix("/books").Subrouter()
	booksRouter.Use(handlers.AuthMiddleware)

	booksRouter.HandleFunc("", handlers.GetBooks).Methods("GET")
	booksRouter.HandleFunc("/{id}", handlers.GetBookByID).Methods("GET")
	booksRouter.HandleFunc("/{id}/borrow", handlers.BorrowBook).Methods("PATCH")
	booksRouter.HandleFunc("/{id}/release", handlers.ReleaseBook).Methods("PATCH")

	adminBooksRouter := booksRouter.PathPrefix("").Subrouter()
	adminBooksRouter.Use(handlers.RoleMiddleware("admin"))
	adminBooksRouter.HandleFunc("", handlers.AddBook).Methods("POST")
	adminBooksRouter.HandleFunc("/{id}", handlers.DeleteBook).Methods("DELETE")
	adminBooksRouter.HandleFunc("/{id}", handlers.UpdateBook).Methods("PUT")

	adminUsersRouter := router.PathPrefix("/users").Subrouter()
	adminUsersRouter.Use(handlers.AuthMiddleware)
	adminUsersRouter.Use(handlers.RoleMiddleware("admin"))
	adminUsersRouter.HandleFunc("", handlers.GetUsers).Methods("GET")
	adminUsersRouter.HandleFunc("/{username}", handlers.GetUserByUsername).Methods("GET")
	http.ListenAndServe(":8000", router)
}
