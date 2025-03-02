package dbhelper

import (
	"rms/database"
	"rms/models"
)

// func IsRestaurantExist(name string) (bool, error) {
// 	query := "SELECT count(id) as is_exist FROM restaurants WHERE name = $1"
// 	var exist bool
// 	err := database.RMS.Get(&exist, query, name)
// 	if err != nil {
// 		return false, err
// 	}
// 	return exist, nil
// }

func CreateRestaurant(restaurant *models.CreateRestaurant, userId string) error {
	query := "INSERT INTO restaurants(name, address, latitude, longitude, created_by) VALUES ($1, $2, $3, $4, $5)"
	_, err := database.RMS.Exec(query, restaurant.Name, restaurant.Address, restaurant.Latitude, restaurant.Longitude, userId)
	return err
}

func GetAllRestaurants() ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	query := "SELECT id, name, address, latitude, longitude, created_at , created_by FROM restaurants WHERE archived_at IS NULL"
	err := database.RMS.Select(&restaurants, query)
	return restaurants, err
}

func GetSpecificRestaurant(id string) (models.Restaurant, error) {
	var restaurant models.Restaurant
	query := "SELECT id, name, address, latitude, longitude, created_at , created_by FROM restaurants WHERE id = $1 AND archived_at IS NULL"
	err := database.RMS.Get(&restaurant, query, id)
	return restaurant, err
}

func GetRestaurantCoordinates(restaurantId string) (models.RestaurantCoordinates, error) {
	var coordinates models.RestaurantCoordinates
	query := "SELECT latitude, longitude FROM restaurants WHERE id = $1 AND archived_at IS NULL"
	err := database.RMS.Get(&coordinates, query, restaurantId)
	return coordinates, err
}
