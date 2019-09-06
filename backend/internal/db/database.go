package db

import (
	"database/sql"
	"fmt"

	// PostgreSQL driver
	_ "github.com/lib/pq"
)

// Database holds connection with extern database.
type Database struct {
	*sql.DB
}

// New creates table in new database and returns pointer to said database
func New(config *Config) (*Database, error) {
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Name)

	database, err := sql.Open(config.Driver, dbInfo)

	return &Database{database}, err
}
