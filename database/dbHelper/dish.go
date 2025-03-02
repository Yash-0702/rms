package dbhelper

import (
	"rms/database"
	"rms/models"
)

func CreateDish(restaurantId string, dish *models.CreateDish) error {
	args := []interface{}{dish.Name, dish.Price, restaurantId}

	query := "INSERT INTO dishes(name, price, restaurant_id) VALUES ($1, $2, $3)"
	_, err := database.RMS.Exec(query, args...)
	return err
}

func IsDishExist(restaurantId string, name string) (bool, error) {
	query := "SELECT count(id) as is_exist FROM dishes WHERE name = $1 AND restaurant_id = $2 AND archived_at IS NULL"
	var exist bool
	err := database.RMS.Get(&exist, query, name, restaurantId)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func GetAllDishesFromAllRestaurants() ([]models.Dish, error) {
	var dishes []models.Dish
	query := "SELECT id, name, price, restaurant_id, created_at FROM dishes WHERE archived_at IS NULL"
	err := database.RMS.Select(&dishes, query)
	return dishes, err
}

func GetAllDishesFromSpecificRestaurant(restaurantId string) ([]models.Dish, error) {
	var dishes []models.Dish
	query := "SELECT id, name, price, restaurant_id, created_at FROM dishes WHERE restaurant_id = $1 AND archived_at IS NULL"
	err := database.RMS.Select(&dishes, query, restaurantId)
	return dishes, err
}

func GetSpecificDish(DishId, restaurantId string) (models.Dish, error) {
	var dish models.Dish
	query := "SELECT id, name, price, restaurant_id, created_at FROM dishes WHERE id = $1 AND restaurant_id = $2  AND archived_at IS NULL"
	err := database.RMS.Get(&dish, query, DishId, restaurantId)
	return dish, err
}

func UpdateDish(DishId, RestaurantId string, dish *models.UpdateDish) error {
	query := "UPDATE dishes SET name = $1, price = $2 WHERE id = $3 AND restaurant_id = $4 AND archived_at IS NULL"
	_, err := database.RMS.Exec(query, dish.Name, dish.Price, DishId, RestaurantId)
	return err
}

func DeleteDish(DishId, RestaurantId string) error {
	query := "UPDATE dishes SET archived_at = now() WHERE id = $1 AND restaurant_id = $2 AND archived_at IS NULL"
	_, err := database.RMS.Exec(query, DishId, RestaurantId)
	return err
}
