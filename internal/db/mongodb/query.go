package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// CRUD Operation for each collection

// CRUD Operation for Admin Collection
// CreateAdminUser creates a new admin user
func CreateAdminUser(user AdminUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := AdminColl.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
func GetAdminUser() error {
	var user AdminUser
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := AdminColl.FindOne(ctx, bson.D{}).Decode(&user)

	if err != nil {
		return err
	}
	return nil
}

// CRUD Operation for User Collection

// CRUD Operation for Brand Collection
