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

// Register user
func CreateUser(res http.ResponseWriter, req *http.Request) {

	var body models.RegisterRequest

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

	// Check if the user already exists
	exist, existErr := dbhelper.IsUserExist(body.Email)
	if existErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, existErr, "failed to check if user exists")
		return
	}

	if exist {
		utils.ResponseError(res, http.StatusBadRequest, nil, "user already exist")
		return
	}

	// Hash the password
	hashedPassword, hashErr := utils.HashPassword(body.Password)
	if hashErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, hashErr, "failed to hash password")
		return
	}
	body.Password = hashedPassword

	// Create the user
	saveErr := dbhelper.CreateUser(&body)
	if saveErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, saveErr, "failed to create user")
		return
	}

	// Return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"message": "User created successfully",
	})

}

// login user
func LoginUser(res http.ResponseWriter, req *http.Request) {

	var body models.LoginRequest

	// Decode the request body
	if decodeErr := utils.DecodeJSONBody(req.Body, &body); decodeErr != nil {
		utils.ResponseError(res, http.StatusBadRequest, decodeErr, "invalid request body")
		return
	}

	// Validate the request body
	v := validator.New()
	if valdationErr := v.Struct(body); valdationErr != nil {
		utils.ResponseError(res, http.StatusBadRequest, valdationErr, "validation error")
		return
	}

	// check if user exist or not
	exist, existErr := dbhelper.IsUserActive(body.Email)
	if existErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, existErr, "failed to check if user exists")
		return
	}
	if !exist {
		utils.ResponseError(res, http.StatusBadRequest, nil, "user not found")
		return
	}

	// check if password is correct or not
	isCorrect, isCorrectErr := dbhelper.IsPasswordCorrect(&body)
	if isCorrectErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, isCorrectErr, "failed to check if password is correct")
		return
	}
	if !isCorrect {
		utils.ResponseError(res, http.StatusBadRequest, nil, "password is incorrect")
		return
	}

	// get user id
	userId, userIdErr := dbhelper.GetUserId(body.Email)
	if userIdErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, userIdErr, "failed to get user id")
		return
	}

	// create session for user
	sessionId, sessionErr := dbhelper.CreateSession(userId)
	if sessionErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, sessionErr, "failed to create session")
		return
	}

	// Get user role
	role, roleErr := dbhelper.GetRole(userId)
	if roleErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, roleErr, "failed to get user role")
		return
	}

	// generate token
	token, tokenErr := utils.GenerateJWT(body.Email, userId, sessionId, role)
	if tokenErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, tokenErr, "failed to generate token")
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"message": "User login successfully",
		"token":   "Bearer " + token,
	})

}

func GetAllUsersByAdminAndSubAdmin(res http.ResponseWriter, req *http.Request) {

	// get all users
	users, usersErr := dbhelper.GetAllUsersByAdmin()
	if usersErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, usersErr, "failed to get users")
		return
	}

	// check if users not exist return []
	if len(users) == 0 {
		utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
			"users": []string{},
		})
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"users": users,
	})
}

func GetUser(res http.ResponseWriter, req *http.Request) {

	// get context
	UserCtx := middlewares.UserContext(req)
	userId := UserCtx.UserId

	// get user info from db
	userInfo, userInfoErr := dbhelper.GetUserInfo(userId)
	if userInfoErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, userInfoErr, "failed to get user info")
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"userInfo": userInfo,
	})

}

func AddAddress(res http.ResponseWriter, req *http.Request) {

	// get user id from context
	userCtx := middlewares.UserContext(req)
	userId := userCtx.UserId

	// Decode the request body
	var body models.AddAddress
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

	// add address
	addressErr := dbhelper.AddAddress(userId, &body)
	if addressErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, addressErr, "failed to add address")
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Address added successfully",
	})
}

func GetAllAddress(res http.ResponseWriter, req *http.Request) {

	// get user id from context
	userCtx := middlewares.UserContext(req)
	userId := userCtx.UserId

	// get address
	address, addressErr := dbhelper.GetAllAddress(userId)
	if addressErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, addressErr, "failed to get address")
		return
	}

	// check if address not exist return []
	if len(address) == 0 {
		utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
			"address": []string{},
		})
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"address": address,
	})
}

func GetSpecificAddress(res http.ResponseWriter, req *http.Request) {

	// get user id from context
	userCtx := middlewares.UserContext(req)
	userId := userCtx.UserId

	// get address id
	addressId := chi.URLParam(req, "addressId")

	// get address
	address, addressErr := dbhelper.GetSpecificAddress(userId, addressId)
	if addressErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, addressErr, "failed to get address")
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"address": address,
	})
}

func UpdateAddress(res http.ResponseWriter, req *http.Request) {

	// get user id from context
	userCtx := middlewares.UserContext(req)
	userId := userCtx.UserId

	// get address id
	addressId := chi.URLParam(req, "addressId")

	// Decode the request body

	var body models.UpdateAddress
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

	// update address
	addressErr := dbhelper.UpdateAddress(userId, addressId, &body)
	if addressErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, addressErr, "failed to update address")
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Address updated successfully",
	})
}

func DeleteAddress(res http.ResponseWriter, req *http.Request) {

	// get user id from context
	userCtx := middlewares.UserContext(req)
	userId := userCtx.UserId

	// get address id
	addressId := chi.URLParam(req, "addressId")

	// delete address
	addressErr := dbhelper.DeleteAddress(userId, addressId)
	if addressErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, addressErr, "failed to delete address")
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Address deleted successfully",
	})
}

func CalculateDistance(res http.ResponseWriter, req *http.Request) {

	// Decode the request body
	var body models.DistanceRequest
	if decodeErr := utils.DecodeJSONBody(req.Body, &body); decodeErr != nil {
		utils.ResponseError(res, http.StatusBadRequest, decodeErr, "invalid request body")
		return
	}

	// Validate the request body
	v := validator.New()
	if validationErr := v.Struct(&body); validationErr != nil {
		utils.ResponseError(res, http.StatusBadRequest, validationErr, "validation error")
		return
	}

	// get user id from context
	userCtx := middlewares.UserContext(req)
	userId := userCtx.UserId

	// fmt.Println(userId)

	// get user address
	UserAddressCoordinates, UserAddressErr := dbhelper.GetUserCoordinates(userId, body.UserAddressId)
	if UserAddressErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, UserAddressErr, "failed to get user address")
		return
	}

	// get restaurant address
	RestaurantAddressCoordinates, RestaurantAddressErr := dbhelper.GetRestaurantCoordinates(body.RestaurantAddressId)
	if RestaurantAddressErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, RestaurantAddressErr, "failed to get restaurant address")
		return
	}

	// calculate distance
	distance, distanceErr := dbhelper.CalculateDistance(&UserAddressCoordinates, &RestaurantAddressCoordinates)
	if distanceErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, distanceErr, "failed to calculate distance")
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, struct {
		Distance float64 `json:"distance"`
	}{distance})
	// utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
	// 	"distance": distance,
	// })
}

func LogoutUser(res http.ResponseWriter, req *http.Request) {

	// get user id from context
	userCtx := middlewares.UserContext(req)
	session_id := userCtx.SessionId

	// delete session
	sessionErr := dbhelper.DeleteSession(session_id)
	if sessionErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, sessionErr, "failed to delete session")
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"message": "User logout successfully",
	})
}

func DeactivateUser(res http.ResponseWriter, req *http.Request) {

	// get user id from context
	userCtx := middlewares.UserContext(req)
	userId := userCtx.UserId

	// delete user
	DelUserErr := dbhelper.DeleteUser(userId)
	if DelUserErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, DelUserErr, "failed to delete user")
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"message": "User deleted successfully",
	})

}

// func ActivateUser(res http.ResponseWriter, req *http.Request) {
// }
