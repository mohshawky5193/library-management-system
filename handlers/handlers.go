package handlers

import (
	"encoding/json"
	"library_management_system/apperrors"
	"library_management_system/config/jsonconfig"
	"library_management_system/models"
	"library_management_system/services/bookservice"
	"library_management_system/services/userservice"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

const IDPathVariable = "id"
const UserRole = "user"
const UsernamePathVariable = "username"

var jwtKey = []byte(os.Getenv("JWT_KEY")) // Store this securely
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		panic(&apperrors.CredentialsDecodingError{})
	}

	user, err := userservice.RegisterUser(creds.Username, creds.Password, UserRole, r.Context())
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func GenerateJWT(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		panic(&apperrors.CredentialsDecodingError{})
	}

	user, err := userservice.AuthenticateUser(creds.Username, creds.Password, r.Context())
	if err != nil {
		panic(&apperrors.UnauthenticatedUserError{})
	}
	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		jsonconfig.UsernameClaimKey:   user.Username,
		jsonconfig.RoleClaimKey:       user.Role,
		jsonconfig.ExpirationClaimKey: time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}

	w.Header().Set(jsonconfig.ContentType, jsonconfig.ApplicationJson)
	json.NewEncoder(w).Encode(map[string]string{jsonconfig.TokenKey: tokenString})
}

func GetBookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars[IDPathVariable]

	book, err := bookservice.GetBookByID(id, r.Context())
	if err != nil {
		panic(&apperrors.BookNotFoundError{BookID: id})
	}

	json.NewEncoder(w).Encode(book)
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := bookservice.GetAllBooks(r.Context())
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(books)
}

func BorrowBook(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(jsonconfig.UsernameContextKey).(string)
	vars := mux.Vars(r)
	id := vars[IDPathVariable]

	success, err := bookservice.BorrowBook(id, username, r.Context())

	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(success)
}

func ReleaseBook(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(jsonconfig.UsernameContextKey).(string)
	vars := mux.Vars(r)
	id := vars[IDPathVariable]

	success, err := bookservice.ReleaseBook(id, username, r.Context())

	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(success)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		panic(err)
	}

	success, err := bookservice.AddBook(book, r.Context())

	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(success)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars[IDPathVariable]

	success, err := bookservice.DeleteBook(id, r.Context())

	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(success)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars[IDPathVariable]
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		panic(err)
	}

	success, err := bookservice.UpdateBook(id, book, r.Context())

	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(success)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := userservice.GetAllUsers(r.Context())
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(users)
}

func GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars[UsernamePathVariable]
	user, err := userservice.FindUser(r.Context(), username)
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(user)
}
