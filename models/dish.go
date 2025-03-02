package models

type CreateDish struct {
	Name  string `json:"name" db:"name" validate:"required"`
	Price int    `json:"price" db:"price" validate:"required"`
	// RestaurantId string `json:"restaurant_id" db:"restaurant_id"`
}

// get all dishes or get a specific dish
type Dish struct {
	ID           string `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Price        int    `json:"price" db:"price"`
	RestaurantId string `json:"restaurant_id" db:"restaurant_id"`
	CreatedAt    string `json:"created_at" db:"created_at"`
}

type UpdateDish struct {
	Name  string `json:"name" db:"name" validate:"required"`
	Price int    `json:"price" db:"price" validate:"required"`
}
