package userservice

import (
	"context"
	"library_management_system/apperrors"
	"library_management_system/config/dbconfig"
	"library_management_system/db"
	"library_management_system/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser creates a new user in the database
func RegisterUser(username, password, role string, ctx context.Context) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	err = checkIfUserExists(ctx, username)

	if err != nil {
		return nil, err
	}

	user := models.User{
		Username: username,
		Password: string(hashedPassword),
		Role:     role,
	}

	result, err := db.UsersCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return &user, nil
}

// AuthenticateUser verifies user credentials
func AuthenticateUser(username, password string, ctx context.Context) (*models.User, error) {
	var user models.User
	filter := bson.M{dbconfig.Username: username}
	err := db.UsersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func checkIfUserExists(ctx context.Context, username string) error {
	_, error := findUserByUsername(ctx, bson.M{dbconfig.Username: username})
	if error == nil {
		return &apperrors.UsernameAlreadyExistsError{Username: username}
	}

	return nil
}

func FindUser(ctx context.Context, username string) (*models.User, error) {
	user, error := findUserByUsername(ctx, bson.M{dbconfig.Username: username})
	if error != nil {
		return nil, error
	}

	return user, nil
}

func UpdateUser(ctx context.Context, user *models.User) (bool, error) {
	filter := bson.M{dbconfig.Username: user.Username}
	update := bson.M{dbconfig.SetOperator: bson.M{dbconfig.BorrowedBookIDs: user.BorrowedBookIDs}}
	_, err := db.UsersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func findUserByUsername(ctx context.Context, filter interface{}) (*models.User, error) {
	var user models.User
	error := db.UsersCollection.FindOne(ctx, filter).Decode(&user)

	return &user, error
}
