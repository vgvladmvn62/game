package ec

// Config for EC.
type Config struct {
	Host struct {
		API string `envconfig:"default=https://localhost:9002/rest/v2"`
		Static string `envconfig:"default=https://localhost:9002"`
	}

	Products struct {
		Site string `envconfig:"default=electronics"`
	}

	HTTP struct {
		Client struct {
			Timeout struct {
				Seconds int `envconfig:"default=20"`
			}
		}
	}
}
