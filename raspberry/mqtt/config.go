package mqtt

// Config for MQTT supporting envconfig for environmental variable parsing.
type Config struct {
	Broker string `envconfig:"default=tcp://test.mosquitto.org:1883"`

	Topic string `envconfig:"default=bullshelves"`

	KeepAlive struct {
		Seconds int `envconfig:"default=2"`
	}

	Timeout struct {
		Seconds int `envconfig:"default=2"`
	}

	Disconnect struct {
		Milliseconds int `envconfig:"default=250"`
	}
}
