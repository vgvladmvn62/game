package server

import "github.com/go-http-utils/logger"

// Config for server.
type Config struct {
	Port int    `envconfig:"default=8080"`
	IP   string `envconfig:"default=0.0.0.0"`

	Logger struct {
		Type logger.Type `envconfig:"default=3"` // TODO: consider changing logger to include request date
	}
}
