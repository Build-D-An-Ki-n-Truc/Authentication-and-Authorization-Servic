package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

// AdminUser is a struct that represents an admin user
type AdminUser struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
}

// NormalUser is a struct that represents a normal user
type NormalUser struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Phone    string             `bson:"phone" json:"phone"`
}

// BrandUser is a struct that represents a brand user
type BrandUser struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
	Name     string             `bson:"name" json:"name"`
	Field    string             `bson:"field" json:"field"`
}

//
