package mqtt

type publisher interface {
	Publish(Command) error
}

// Commander helps in constructing commands for mqtt client.
type Commander struct {
	cli publisher
}

// NewCommander creates commander, that sends commands through given MQTT client.
func NewCommander(cli publisher) *Commander {
	return &Commander{
		cli: cli,
	}
}

// PlatformColor lights up platform.
func (cm *Commander) PlatformColor(platformID byte, rgb RGB) error {
	cmd := Command{
		Platform: platformID,
		Command:  CmdOn,
		RGB:      &rgb,
	}
	return cm.cli.Publish(cmd)
}

// PlatformRotateColor rotates pixel around the platform.
func (cm *Commander) PlatformRotateColor(platformID byte, rgb RGB, delay int) error {
	cmd := Command{
		Platform: platformID,
		Command:  CmdAnimate,
		RGB:      &rgb,
	}
	return cm.cli.Publish(cmd)
}

// PlatformFadePixels turns on every single pixel one by one.
func (cm *Commander) PlatformFadePixels(platformID byte, rgb RGB, delay byte) error {
	cmd := Command{
		Platform: platformID,
		Command:  CmdFade,
		RGB:      &rgb,
	}
	return cm.cli.Publish(cmd)
}

// PlatformSensorReading enables reading callback events.
func (cm *Commander) PlatformSensorReading(platformID byte) error {
	cmd := Command{
		Platform: platformID,
		Command:  CmdSensor,
	}
	return cm.cli.Publish(cmd)
}

// DisableAllLights turns off all lights.
func (cm *Commander) DisableAllLights() error {
	cmd := Command{
		Command: CmdOffAll,
	}
	return cm.cli.Publish(cmd)
}
