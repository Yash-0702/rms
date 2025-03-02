package handlers

import (
	"net/http"
	dbhelper "rms/database/dbHelper"
	"rms/models"
	"rms/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func CreateDish(res http.ResponseWriter, req *http.Request) {

	var body models.CreateDish

	// get restaurant id from chi params
	restaurantId := chi.URLParam(req, "restaurantId")

	// validate restaurant id
	if restaurantId == "" {
		utils.ResponseError(res, http.StatusBadRequest, nil, "invalid restaurant id")
		return
	}

	// Decode the request body
	if decodeErr := utils.DecodeJSONBody(req.Body, &body); decodeErr != nil {
		utils.ResponseError(res, http.StatusBadRequest, decodeErr, "invalid request body")
		return
	}

	// Validate the request body
	v := validator.New()
	if validationErr := v.Struct(body); validationErr != nil {
		utils.ResponseError(res, http.StatusBadRequest, validationErr, "validation error")
		return
	}

	// check if dish exist or not
	exist, existErr := dbhelper.IsDishExist(restaurantId, body.Name)
	if existErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, existErr, "failed to check if dish exists")
		return
	}
	if exist {
		utils.ResponseError(res, http.StatusBadRequest, nil, "dish already exist")
		return
	}

	// create dish
	createErr := dbhelper.CreateDish(restaurantId, &body)
	if createErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, createErr, "failed to create dish")
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Dish created successfully",
	})

}

// get all dishes of all restaurants
func GetAllDishesFromAllRestaurants(res http.ResponseWriter, req *http.Request) {

	// get all dishes
	dishes, dishesErr := dbhelper.GetAllDishesFromAllRestaurants()
	if dishesErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, dishesErr, "failed to get dishes")
		return
	}

	// check if dishes not exist return []
	if len(dishes) == 0 {
		utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
			"dishes": []string{},
		})
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"dishes": dishes,
	})
}

// get all dishes of a specific restaurant
func GetAllDishesFromSpecificRestaurant(res http.ResponseWriter, req *http.Request) {

	// get restaurant id from chi params
	restaurantId := chi.URLParam(req, "restaurantId")

	// validate restaurant id
	if restaurantId == "" {
		utils.ResponseError(res, http.StatusBadRequest, nil, "invalid restaurant id")
		return
	}

	// get all dishes
	dishes, dishesErr := dbhelper.GetAllDishesFromSpecificRestaurant(restaurantId)
	if dishesErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, dishesErr, "failed to get dishes")
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"dishes": dishes,
	})
}

// get a specific dish
func GetSpecificDish(res http.ResponseWriter, req *http.Request) {

	// get dish id from chi params
	dishId := chi.URLParam(req, "dishId")
	// get restaurant id from chi params
	restaurantId := chi.URLParam(req, "restaurantId")

	// validate dish id
	if dishId == "" {
		utils.ResponseError(res, http.StatusBadRequest, nil, "invalid dish id")
		return
	}

	// get specific dish
	dish, dishErr := dbhelper.GetSpecificDish(dishId, restaurantId)
	if dishErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, dishErr, "failed to get specific dish")
		return
	}

	// if dish not exist
	if dish == (models.Dish{}) {
		utils.ResponseJSON(res, http.StatusNotFound, map[string]interface{}{
			"message": "Dish not found",
		})
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"dish": dish,
	})
}

func UpdateDish(res http.ResponseWriter, req *http.Request) {

	// get dish id from chi params
	dishId := chi.URLParam(req, "dishId")

	// get restaurant id from chi params
	restaurantId := chi.URLParam(req, "restaurantId")

	// validate dish id
	if dishId == "" || restaurantId == "" {
		utils.ResponseError(res, http.StatusBadRequest, nil, "invalid dish id")
		return
	}

	// Decode the request body
	var body models.UpdateDish
	if decodeErr := utils.DecodeJSONBody(req.Body, &body); decodeErr != nil {
		utils.ResponseError(res, http.StatusBadRequest, decodeErr, "invalid request body")
		return
	}

	// Validate the request body
	v := validator.New()
	if validationErr := v.Struct(body); validationErr != nil {
		utils.ResponseError(res, http.StatusBadRequest, validationErr, "validation error")
		return
	}

	// update dish
	updateErr := dbhelper.UpdateDish(dishId, restaurantId, &body)
	if updateErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, updateErr, "failed to update dish")
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Dish updated successfully",
	})

}

func DeleteDish(res http.ResponseWriter, req *http.Request) {

	// get dish id from chi params
	dishId := chi.URLParam(req, "dishId")
	// get restaurant id from chi params
	restaurantId := chi.URLParam(req, "restaurantId")

	// validate dish id
	if dishId == "" || restaurantId == "" {
		utils.ResponseError(res, http.StatusBadRequest, nil, "invalid dish id")
		return
	}

	// delete dish
	deleteErr := dbhelper.DeleteDish(dishId, restaurantId)
	if deleteErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, deleteErr, "failed to delete dish")
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Dish deleted successfully",
	})
}
