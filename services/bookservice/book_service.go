package bookservice

import (
	"context"
	"library_management_system/apperrors"
	"library_management_system/config/dbconfig"
	"library_management_system/db"
	"library_management_system/models"
	"library_management_system/services/userservice"

	"go.mongodb.org/mongo-driver/bson"
)

// GetBookByID retrieves a book by its ID from the database.
func GetBookByID(id string, ctx context.Context) (*models.Book, error) {
	var book models.Book
	filter := bson.M{dbconfig.ID: id}
	err := db.BooksCollection.FindOne(ctx, filter).Decode(&book)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// GetAllBooks retrieves all books from the database.
func GetAllBooks(ctx context.Context) ([]models.Book, error) {
	var books []models.Book
	cursor, err := db.BooksCollection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var book models.Book
		if err := cursor.Decode(&book); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// BorrowBook allows a user to borrow a book, updating both the book and user records.
func BorrowBook(bookId, username string, ctx context.Context) (bool, error) {
	var book models.Book
	filter := bson.M{dbconfig.ID: bookId}
	err := db.BooksCollection.FindOne(ctx, filter).Decode(&book)
	if err != nil {
		return false, &apperrors.BookNotFoundError{BookID: bookId}
	}

	// Find the user and update borrowed books
	user, err := userservice.FindUser(ctx, username)
	if err != nil {
		return false, err
	}

	if user.BorrowedBookIDs == nil {
		user.BorrowedBookIDs = []string{}
	}

	for _, b := range user.BorrowedBookIDs {
		if b == bookId {
			return false, &apperrors.AlreadyHaveBookError{BookTitle: book.Title}
		}
	}

	if book.Amount == 0 {
		return false, &apperrors.AmountIsZeroError{BookTitle: book.Title}
	}

	book.Amount--
	if book.OwnedBy == nil {
		book.OwnedBy = make([]string, 0)
	}
	book.OwnedBy = append(book.OwnedBy, username)
	user.BorrowedBookIDs = append(user.BorrowedBookIDs, book.ID)
	update := bson.M{dbconfig.SetOperator: bson.M{dbconfig.Amount: book.Amount, dbconfig.OwnedBy: book.OwnedBy}}
	_, err = db.BooksCollection.UpdateByID(ctx, bookId, update)
	if err != nil {
		return false, err
	}

	_, err = userservice.UpdateUser(ctx, user)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ReleaseBook(bookId, username string, ctx context.Context) (bool, error) {
	// Check and update the book document
	var book models.Book
	filter := bson.M{dbconfig.ID: bookId}
	err := db.BooksCollection.FindOne(ctx, filter).Decode(&book)
	if err != nil {
		return false, &apperrors.BookNotFoundError{BookID: bookId}
	}

	// Find the user and check borrowed books
	user, err := userservice.FindUser(ctx, username)
	if err != nil {
		return false, err
	}
	var found bool = false
	for i, b := range user.BorrowedBookIDs {
		if b == bookId {
			// Remove the book from borrowed list
			user.BorrowedBookIDs = append(user.BorrowedBookIDs[:i], user.BorrowedBookIDs[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return false, &apperrors.BookNotBorrowedError{BookTitle: book.Title}
	}

	for j, u := range book.OwnedBy {
		if u == user.Username {
			book.OwnedBy = append(book.OwnedBy[:j], book.OwnedBy[j+1:]...)
			break
		}
	}

	// Update the book amount
	book.Amount++
	update := bson.M{dbconfig.SetOperator: bson.M{dbconfig.Amount: book.Amount, dbconfig.OwnedBy: book.OwnedBy}}
	_, err = db.BooksCollection.UpdateByID(ctx, bookId, update)
	if err != nil {
		return false, err
	}

	// Update the user document
	_, err = userservice.UpdateUser(ctx, user)
	if err != nil {
		return false, err
	}

	return true, nil
}

func AddBook(book models.Book, ctx context.Context) (bool, error) {

	filter := bson.M{dbconfig.ID: book.ID}
	var existingBook models.Book
	err := db.BooksCollection.FindOne(ctx, filter).Decode(&existingBook)
	if err == nil {
		return false, &apperrors.BookWithSameIDError{BookID: book.ID}
	}

	_, err = validateBookDataForAddition(book)
	if err != nil {
		return false, err
	}
	_, err = db.BooksCollection.InsertOne(ctx, book)
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteBook(id string, ctx context.Context) (bool, error) {
	filter := bson.M{dbconfig.ID: id}
	var book models.Book
	err := db.BooksCollection.FindOne(ctx, filter).Decode(&book)
	if err != nil {
		return false, &apperrors.BookNotFoundError{BookID: id}
	}
	if len(book.OwnedBy) > 0 {
		return false, &apperrors.DeleteBorrowedBookError{BookTitle: book.Title}
	}
	_, err = db.BooksCollection.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}
	return true, nil
}

func UpdateBook(id string, book models.Book, ctx context.Context) (bool, error) {

	filter := bson.M{dbconfig.ID: id}
	var oldBook models.Book
	err := db.BooksCollection.FindOne(ctx, filter).Decode(&oldBook)
	if err != nil {
		return false, &apperrors.BookNotFoundError{BookID: id}
	}
	_, err = validateBookDataForUpdate(book)
	if err != nil {
		return false, err
	}
	update := bson.M{dbconfig.SetOperator: bson.M{dbconfig.Title: book.Title, dbconfig.Author: book.Author, dbconfig.Amount: book.Amount}}

	_, err = db.BooksCollection.UpdateByID(ctx, id, update)

	if err != nil {
		return false, err
	}

	return true, nil
}

func validateBookDataForAddition(book models.Book) (bool, error) {
	var errorMessages []string = make([]string, 0)

	if book.ID == "" {
		errorMessages = append(errorMessages, "book id is empty")
	}
	if book.Title == "" {
		errorMessages = append(errorMessages, "title is empty")
	}
	if book.Author == "" {
		errorMessages = append(errorMessages, "author is empty")
	}
	if book.Amount <= 0 {
		errorMessages = append(errorMessages, "amount is less than or equal to zero")
	}
	if book.OwnedBy != nil {
		errorMessages = append(errorMessages, "cannot add book with non empty owned_by")
	}

	if len(errorMessages) > 0 {
		return false, &apperrors.BookValidationError{ErrorMessages: errorMessages}
	}
	return true, nil
}

func validateBookDataForUpdate(book models.Book) (bool, error) {
	var errorMessages []string
	if book.ID != "" {
		errorMessages = append(errorMessages, "book id is not empty")
	}
	if book.Title == "" {
		errorMessages = append(errorMessages, "title is empty")
	}
	if book.Author == "" {
		errorMessages = append(errorMessages, "author is empty")
	}
	if book.Amount <= 0 {
		errorMessages = append(errorMessages, "amount is less than or equal to zero")
	}
	if book.OwnedBy != nil {
		errorMessages = append(errorMessages, "cannot edit book owned_by")
	}

	if len(errorMessages) > 0 {
		return false, &apperrors.BookValidationError{ErrorMessages: errorMessages}
	}
	return true, nil
}
