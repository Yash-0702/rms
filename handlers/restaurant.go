package handlers

import (
	"net/http"
	dbhelper "rms/database/dbHelper"
	"rms/middlewares"
	"rms/models"
	"rms/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// CreateRestaurant
func CreateRestaurant(res http.ResponseWriter, req *http.Request) {
	var body models.CreateRestaurant

	//get context
	UserCtx := middlewares.UserContext(req)
	userId := UserCtx.UserId

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

	// create restaurant
	createErr := dbhelper.CreateRestaurant(&body, userId)
	if createErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, createErr, "failed to create restaurant")
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Restaurant created successfully",
	})
}

// GetAllRestaurants
func GetAllRestaurants(res http.ResponseWriter, req *http.Request) {

	// get all restaurants
	restaurants, restaurantsErr := dbhelper.GetAllRestaurants()
	if restaurantsErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, restaurantsErr, "failed to get restaurants")
		return
	}

	// check if restaurants not exist return []
	if len(restaurants) == 0 {
		utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
			"restaurants": []string{},
		})
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"restaurants": restaurants,
	})
}

// Get specific Restaurant

func GetSpecificRestaurant(res http.ResponseWriter, req *http.Request) {

	// get restaurant id from chi params
	restaurantId := chi.URLParam(req, "restaurantId")

	// validate restaurant id
	if restaurantId == "" {
		utils.ResponseError(res, http.StatusBadRequest, nil, "invalid restaurant id")
		return
	}

	// get specific restaurant
	restaurant, restaurantErr := dbhelper.GetSpecificRestaurant(restaurantId)
	if restaurantErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, restaurantErr, "failed to get specific restaurant")
		return
	}

	// if restaurant not exist
	if restaurant == (models.Restaurant{}) {
		utils.ResponseJSON(res, http.StatusNotFound, map[string]interface{}{
			"message": "Restaurant not found",
		})
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"restaurant": restaurant,
	})

}
