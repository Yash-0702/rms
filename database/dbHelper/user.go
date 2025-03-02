package dbhelper

import (
	"fmt"
	"rms/database"
	"rms/models"
	"rms/utils"
)

func CreateUser(usr *models.RegisterRequest) error {
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3)"
	_, err := database.RMS.Exec(query, usr.Username, usr.Email, usr.Password)

	return err
}

func IsUserExist(email string) (bool, error) {
	query := "SELECT count(id) as is_exist FROM users WHERE email = $1"
	var exist bool
	err := database.RMS.Get(&exist, query, email)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func IsUserActive(email string) (bool, error) {
	query := "SELECT count(id) as is_exist FROM users WHERE email = $1 AND archived_at IS NULL"
	var active bool
	err := database.RMS.Get(&active, query, email)
	if err != nil {
		return false, err
	}
	return active, nil
}

func CreateSession(userId string) (string, error) {
	var sessionID string
	query := `INSERT INTO user_session(user_id)
    			VALUES ($1) RETURNING id `

	Err := database.RMS.Get(&sessionID, query, userId)
	if Err != nil {
		return "", Err
	}
	return sessionID, nil
}

func DeleteSession(session_id string) error {
	query := "UPDATE user_session SET archived_at = now() WHERE id = $1 AND archived_at IS NULL"
	_, err := database.RMS.Exec(query, session_id)
	return err
}

func GetUserByEmail(email string) (string, error) {
	query := "SELECT id FROM users WHERE email = $1 AND archived_at IS NULL"
	var id string
	err := database.RMS.Get(&id, query, email)
	if err != nil {
		return "", err
	}
	return id, nil
}

func GetUserId(email string) (string, error) {
	query := "SELECT id FROM users WHERE email = $1 AND archived_at IS NULL"
	var id string
	err := database.RMS.Get(&id, query, email)
	if err != nil {
		return "", err
	}
	return id, nil
}

func IsPasswordCorrect(usr *models.LoginRequest) (bool, error) {
	query := "SELECT password FROM users where email = $1 AND archived_at IS NULL"
	var hashedPassword string
	err := database.RMS.Get(&hashedPassword, query, usr.Email)
	if err != nil {
		return false, err
	}

	// check if password is incorrect or not
	IsPasswordCorrect := utils.UnhashPassword(hashedPassword, usr.Password)
	if IsPasswordCorrect != nil {
		return false, nil
	}
	return true, nil
}

func DeleteUser(userId string) error {
	query := "UPDATE users SET archived_at = now() WHERE id = $1"
	_, err := database.RMS.Exec(query, userId)
	return err
}

func GetUserInfo(userId string) (models.UserInfo, error) {
	var userInfo models.UserInfo
	query := "SELECT  id, username, email, role, created_at FROM users WHERE id = $1"
	err := database.RMS.Get(&userInfo, query, userId)
	return userInfo, err
}

func GetRole(userId string) (string, error) {
	var role string
	query := "SELECT role FROM users WHERE id = $1 AND archived_at IS NULL"
	err := database.RMS.Get(&role, query, userId)
	return role, err
}

func IsSessionExist(sessionId string) (bool, error) {
	query := "SELECT count(id) as is_exist FROM user_session WHERE id = $1 AND archived_at IS NULL"
	var exist bool
	err := database.RMS.Get(&exist, query, sessionId)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func GetAllUsersByAdmin() ([]models.UserInfo, error) {
	var users []models.UserInfo
	query := "SELECT id, username, email,role, created_at FROM users WHERE role = $1 AND archived_at IS NULL"
	err := database.RMS.Select(&users, query, models.RoleUser)
	return users, err
}

func AddAddress(userId string, address *models.AddAddress) error {
	args := []interface{}{userId, address.Address, address.Latitude, address.Longitude}
	query := "INSERT INTO address (user_id, address, latitude, longitude) VALUES ($1, $2, $3, $4)"
	_, err := database.RMS.Exec(query, args...)
	return err
}

func GetAllAddress(userId string) ([]models.GetAddress, error) {
	var address []models.GetAddress
	query := "SELECT id, address, latitude, longitude , user_id FROM address WHERE user_id = $1 AND archived_at IS NULL"
	err := database.RMS.Select(&address, query, userId)
	return address, err
}

func GetUserCoordinates(userId string, UserAddressId string) (models.UserAddressCoordinates, error) {
	var coordinates models.UserAddressCoordinates
	fmt.Println(UserAddressId)
	query := "SELECT latitude, longitude FROM address WHERE user_id = $1 AND id = $2 AND archived_at IS NULL"
	err := database.RMS.Get(&coordinates, query, userId, UserAddressId)
	return coordinates, err
}

func CalculateDistance(userCoordinates *models.UserAddressCoordinates, RestaurtantCoordinates *models.RestaurantCoordinates) (float64, error) {

	args := []interface{}{userCoordinates.Latitude, userCoordinates.Longitude, RestaurtantCoordinates.Latitude, RestaurtantCoordinates.Longitude}
	query := "SELECT earth_distance(ll_to_earth($1, $2), ll_to_earth($3, $4))/1000 as distance_Km"
	var distance float64
	err := database.RMS.Get(&distance, query, args...)
	return distance, err

}

func GetSpecificAddress(userId, addressId string) (models.GetAddress, error) {
	var address models.GetAddress
	query := "SELECT id, address, latitude, longitude , user_id FROM address WHERE user_id = $1 AND id = $2 AND archived_at IS NULL"
	err := database.RMS.Get(&address, query, userId, addressId)
	return address, err
}

func UpdateAddress(userId, addressId string, address *models.UpdateAddress) error {
	args := []interface{}{address.Address, address.Latitude, address.Longitude, userId, addressId}
	query := "UPDATE address SET address = $1, latitude = $2, longitude = $3 WHERE user_id = $4 AND id = $5 AND archived_at IS NULL"
	_, err := database.RMS.Exec(query, args...)
	return err
}

func DeleteAddress(userId, addressId string) error {
	query := "UPDATE address SET archived_at = now() WHERE user_id = $1 AND id = $2 AND archived_at IS NULL"
	_, err := database.RMS.Exec(query, userId, addressId)
	return err
}
