package models

type CreateSubAdmin struct {
	Username string `json:"username" db:"username" validate:"required"`
	Email    string `json:"email" db:"email" validate:"email"`
	Password string `json:"password" db:"password" validate:"gte=6,lte=15"`
}

type SubAdmin struct {
	ID         string `json:"id" db:"id"`
	Name       string `json:"name" db:"username"`
	Email      string `json:"email" db:"email"`
	Role       string `json:"role" db:"role"`
	Created_at string `json:"created_at" db:"created_at"`
}
