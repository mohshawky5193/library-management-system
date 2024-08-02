package apperrors

import (
	"fmt"
	"strings"
)

type UsernameAlreadyExistsError struct {
	Username string
}

type UnauthorizedUserError struct {
}

type UnauthenticatedUserError struct {
}

type AmountIsZeroError struct {
	BookTitle string
}

type BookNotBorrowedError struct {
	BookTitle string
}

type AlreadyHaveBookError struct {
	BookTitle string
}

type BookNotFoundError struct {
	BookID string
}

type CredentialsDecodingError struct {
}

type MalFormedTokenError struct {
}

type InvalidTokenError struct {
}

type RoleNotMatchingError struct {
	Role         string
	RequiredRole string
}

type DeleteBorrowedBookError struct {
	BookTitle string
}

type BookValidationError struct {
	ErrorMessages []string
}

type BookWithSameIDError struct {
	BookID string
}

func (e *UsernameAlreadyExistsError) Error() string {
	return fmt.Sprintf("%s already exists", e.Username)
}

func (e *UnauthorizedUserError) Error() string {
	return "User not authorized"
}

func (e *AmountIsZeroError) Error() string {
	return fmt.Sprintf("You can't borrow %s as the amount is zero", e.BookTitle)
}

func (e *BookNotBorrowedError) Error() string {
	return fmt.Sprintf("You can't release  %s because you don't have the book", e.BookTitle)
}

func (e *AlreadyHaveBookError) Error() string {
	return fmt.Sprintf("You can't borrow  %s because you already have the book", e.BookTitle)
}

func (e *BookNotFoundError) Error() string {
	return fmt.Sprintf("No Book With Id %s", e.BookID)
}

func (e *CredentialsDecodingError) Error() string {
	return "Error decoding credentials"
}

func (e *MalFormedTokenError) Error() string {
	return "Missing or Malformed token"
}

func (e *InvalidTokenError) Error() string {
	return "Invalid token"
}

func (e *RoleNotMatchingError) Error() string {
	return fmt.Sprintf("your role %s is not matching %s", e.Role, e.RequiredRole)
}

func (e *DeleteBorrowedBookError) Error() string {
	return fmt.Sprintf("cannot delete %s because it is still borrowed", e.BookTitle)
}

func (e *UnauthenticatedUserError) Error() string {
	return "invalid username/password"
}

func (e *BookValidationError) Error() string {
	return strings.Join(e.ErrorMessages, ",")
}

func (e *BookWithSameIDError) Error() string {
	return fmt.Sprintf("book with id %s exists", e.BookID)
}
