package mongodb

import (
	"context"
	"time"
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
