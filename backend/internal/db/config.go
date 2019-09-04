package db

// Config for database.
type Config struct {
	Host     string `envconfig:"default=localhost"`
	Port     int    `envconfig:"default=5432"`
	User     string `envconfig:"default=postgres"`
	Password string `envconfig:"default=password"`
	Name     string `envconfig:"default=postgres"`
	Driver   string `envconfig:"default=postgres"`
}
