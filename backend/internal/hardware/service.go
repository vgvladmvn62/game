package hardware

import (
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db/stands"
	"github.com/kyma-incubator/bullseye-showcase/backend/pkg/mqtt"
)

// Service allows to communicate with hardware via easy API
type Service struct {
	client       mqttClient
	standService standService
	standIDs     []int
	commander    *mqtt.Commander
}

type mqttClient interface {
	Publish(mqtt.Command) error
}

type standService interface {
	GetAllActiveStands() ([]stands.StandDTO, error)
}

// NewHardwareService returns new hardware service.
func NewHardwareService(client mqttClient, standService standService) *Service {
	return &Service{
		client:       client,
		standService: standService,
		commander:    mqtt.NewCommander(client),
	}
}

// TurnOffLights sends platform color command to turn off lights.
func (hs *Service) TurnOffLights() error {
	return hs.commander.DisableAllLights()
}

// TurnOnGreenLight sends platform color command to light up platform with green light.
func (hs *Service) TurnOnGreenLight(platformID byte) error {
	return hs.commander.PlatformFadePixels(platformID, mqtt.RGB{50, 255, 50}, 50)
}
