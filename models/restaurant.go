package models

type CreateRestaurant struct {
	Name      string  `json:"name" validate:"required"`
	Address   string  `json:"address" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

// for get all restaurants or get a specific restaurant
type Restaurant struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Address   string `json:"address" db:"address"`
	Latitude  string `json:"latitude" db:"latitude"`
	Longitude string `json:"longitude" db:"longitude"`
	CreatedAt string `json:"created_at" db:"created_at"`
	CreatedBy string `json:"created_by" db:"created_by"`
}
