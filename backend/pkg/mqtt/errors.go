package mqtt

// BrokerError informs about possible failures when connecting to the broker.
type BrokerError string

const (
	// BrokerConnectionError when
	BrokerConnectionError BrokerError = "BROKER_CONNECTION_ERROR - Could not connect to the MQTT broker"
)

// Error returns error as a string.
func (e BrokerError) Error() string {
	return string(e)
}
