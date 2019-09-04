package slab

// Config for slabs supporting envconfig for environmental variable parsing.
type Config struct {
	Repository struct {
		Filter string `envconfig:"default=ACM"`

		Path string `envconfig:"default=./slabs.db"`
	}
}
