package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	// "github.com/golang-migrate/migrate/v4"
	// "github.com/golang-migrate/migrate/v4/database/postgres"
	// "github.com/jmoiron/sqlx"
)

var (
	RMS *sqlx.DB
)

// ConnectAndMigrate connects to the database and runs migrations
func ConnectAndMigrate(host, port, user, password, dbname, sslmode string) error {

	// connStr := "postgres://postgres:rx@localhost:5432/todo?sslmode=disable"
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	DB, err := sqlx.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Failed to connect to database")
		return err
	}

	err = DB.Ping()
	if err != nil {
		fmt.Println("Failed to ping database")
		return err
	}

	fmt.Println("Connected to database successfully")

	RMS = DB
	return MigrateUp(DB)

}

// runs migrations
func MigrateUp(db *sqlx.DB) error {

	driver, driErr := postgres.WithInstance(db.DB, &postgres.Config{})

	if driErr != nil {
		fmt.Println("Failed to create migration driver")
		return driErr
	}

	m, migErr := migrate.NewWithDatabaseInstance("file://database/migrations", "postgres", driver)

	if migErr != nil {
		fmt.Println("Failed to create migration instance")
		return migErr
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("Failed to run migrations")
		return err
	}

	fmt.Println("Migrations completed successfully")
	return nil
}

// closes the database
func ShutdownDatabase() error {
	return RMS.Close()
}
