package models

const (
	RoleAdmin    = "admin"
	RoleSubAdmin = "sub-admin"
	RoleUser     = "user"
)

type RegisterRequest struct {
	Username string `json:"username" db:"username" validate:"required"`
	Email    string `json:"email" db:"email" validate:"email"`
	Password string `json:"password" db:"password" validate:"gte=6,lte=15"`
	// Role     string `json:"role" db:"role" validate:"required,oneof=admin sub-admin user" `
}

type LoginRequest struct {
	Email    string `json:"email" db:"email" validate:"email"`
	Password string `json:"password" db:"password" validate:"gte=6,lte=15"`
}

type UserInfo struct {
	ID         string `json:"user_id" db:"id"`
	Name       string `json:"name" db:"username"`
	Email      string `json:"email" db:"email"`
	Role       string `json:"role" db:"role"`
	Created_at string `json:"created_at" db:"created_at"`
}

type UserCtx struct {
	UserId    string `json:"user_id"`
	Email     string `json:"email"`
	SessionId string `json:"session_id"`
	Role      string `json:"role"`
}

type AddAddress struct {
	Address   string  `json:"address" db:"address" validate:"required"`
	Latitude  float64 `json:"latitude" db:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" db:"longitude" validate:"required"`
}

type GetAddress struct {
	ID        string `json:"id" db:"id"`
	Address   string `json:"address" db:"address"`
	Latitude  string `json:"latitude" db:"latitude"`
	Longitude string `json:"longitude" db:"longitude"`
	UserId    string `json:"user_id" db:"user_id"`
}

type UpdateAddress struct {
	Address   string  `json:"address" db:"address" validate:"required"`
	Latitude  float64 `json:"latitude" db:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" db:"longitude" validate:"required"`
}

type DistanceRequest struct {
	UserAddressId       string `json:"user_address_id" db:"id" validate:"required"`
	RestaurantAddressId string `json:"restaurant_address_id" db:"id" validate:"required"`
}

type UserAddressCoordinates struct {
	Latitude  float64 `json:"latitude" db:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" db:"longitude" validate:"required"`
}
