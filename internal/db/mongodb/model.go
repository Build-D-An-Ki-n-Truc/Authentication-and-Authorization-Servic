package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

// User is a struct that represents an admin user
type UserStruct struct {
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Role     string             `bson:"role" json:"role"`
	Phone    string             `bson:"phone" json:"phone"`
	IsLocked bool               `bson:"islocked" json:"isLocked"`
	Turn     int                `bson:"turn" json:"turn"`
	BrandID  primitive.ObjectID `bson:"brand_id,omitempty" json:"brand_id,omitempty"`
}

type GPS struct {
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
}
type BrandStruct struct {
	Brand_Name     string               `bson:"brand_name" json:"brand_name"`
	Industry       string               `bson:"industry" json:"industry"`
	Address        string               `bson:"address" json:"address"`
	Status         bool                 `bson:"status" json:"status"`
	GPS_Coordinate GPS                  `bson:"gps_coordinate" json:"gps_coordinate"`
	Events_IDs     []primitive.ObjectID `bson:"events_ids" json:"events_ids"`
	CreatedAt      primitive.DateTime   `bson:"created_at" json:"created_at"`
	UpdatedAt      primitive.DateTime   `bson:"updated_at" json:"updated_at"`
	V              int                  `bson:"__v" json:"__v"`
	Is_deleted     bool                 `bson:"is_deleted" json:"is_deleted"`
}
