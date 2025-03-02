package handlers

import (
	"net/http"
	dbhelper "rms/database/dbHelper"
	"rms/models"
	"rms/utils"

	"github.com/go-playground/validator/v10"
)

func CreateSubAdmin(res http.ResponseWriter, req *http.Request) {

	var body models.CreateSubAdmin

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
		utils.ResponseError(res, http.StatusBadRequest, nil, "sub-admin already exist")
		return
	}

	// Hash the password
	hashedPassword, hashErr := utils.HashPassword(body.Password)
	if hashErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, hashErr, "failed to hash password")
		return
	}
	body.Password = hashedPassword

	// Create the subadmin user
	saveErr := dbhelper.CreateSubAdmin(&body)
	if saveErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, saveErr, "failed to create user")
		return
	}

	// Return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"message": "Sub-admin created successfully",
	})

}

func GetAllSubAdmins(res http.ResponseWriter, req *http.Request) {

	// get all sub-admins
	subAdmins, subAdminsErr := dbhelper.GetAllSubAdmins()
	if subAdminsErr != nil {
		utils.ResponseError(res, http.StatusInternalServerError, subAdminsErr, "failed to get sub-admins")
		return
	}

	// return the response
	utils.ResponseJSON(res, http.StatusOK, map[string]interface{}{
		"subAdmins": subAdmins,
	})

}
