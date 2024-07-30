package mongodb

// AdminUser is a struct that represents an admin user
type AdminUser struct {
	ID       string `json:"_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

// NormalUser is a struct that represents a normal user
type NormalUser struct {
	ID       string `json:"_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// BrandUser is a struct that represents a brand user
type BrandUser struct {
	ID        string  `json:"_id"`
	Username  string  `json:"username"`
	Password  string  `json:"password"`
	Name      string  `json:"name"`
	Field     string  `json:"field"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    string  `json:"status"`
}
