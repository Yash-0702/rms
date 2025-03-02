package dbhelper

import (
	"rms/database"
	"rms/models"
)

func CreateSubAdmin(subAdmin *models.CreateSubAdmin) error {
	args := []interface{}{subAdmin.Username, subAdmin.Email, subAdmin.Password, models.RoleSubAdmin}
	query := "INSERT INTO users (username, email, password, role) VALUES ($1, $2, $3, $4)"
	_, err := database.RMS.Exec(query, args...)
	return err
}

func GetAllSubAdmins() ([]models.SubAdmin, error) {
	var subAdmins []models.SubAdmin
	query := "SELECT id, username, email, role , created_at FROM users WHERE role = $1 AND archived_at IS NULL"
	err := database.RMS.Select(&subAdmins, query, models.RoleSubAdmin)
	return subAdmins, err
}
