package mongodb

// User is a struct that represents an admin user
type UserStruct struct {
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
	Name     string `bson:"name" json:"name"`
	Email    string `bson:"email" json:"email"`
	Role     string `bson:"role" json:"role"`
	Phone    string `bson:"phone" json:"phone"`
	IsLocked bool   `bson:"islocked" json:"isLocked"`
}
