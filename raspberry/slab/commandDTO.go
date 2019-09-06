package slab

// List of commands that can be sent.
const (
	mqttOn byte = iota
	mqttOFF
	mqttAnimate
	mqttSensor
	mqttFade
	mqttOffAll
)

// CommandDTO describes data that has to be provided for hardware.
type CommandDTO struct {
	ID      byte `json:"id"`
	Command byte `json:"command,omitempty"`
	RGB     *RGB `json:"rgb,omitempty"`
	Delay   byte `json:"delay,omitempty"`
}
