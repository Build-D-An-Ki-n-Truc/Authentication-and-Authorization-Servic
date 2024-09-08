package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// CRUD Operation for each collection

// CRUD Operation for User Collection
// CreateUser creates a new  user
func CreateUser(user UserStruct) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if the username already exists
	var existingUser UserStruct
	err := UserColl.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		return errors.New("username already exists")
	}

	// Insert the new user
	_, err = UserColl.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// ReadUser reads an user by username
func ReadUser(username string) (UserStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user UserStruct
	err := UserColl.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// ReadAllUsers reads all users
func ReadAllUsers() ([]UserStruct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := UserColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []UserStruct
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser updates an user by ID
func UpdateUser(username string, user UserStruct) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := UserColl.ReplaceOne(ctx, bson.M{"username": username}, user)
	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes an user by ID
func DeleteUser(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := UserColl.DeleteOne(ctx, bson.M{"username": username})
	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}
	if err != nil {
		return err
	}
	return nil
}
