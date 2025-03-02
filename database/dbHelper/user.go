package dbhelper

import (
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
