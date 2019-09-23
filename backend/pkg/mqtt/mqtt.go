package mqtt

import (
	"encoding/json"
	"log"
	"time"

	evbus "github.com/asaskevich/EventBus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// List of commands that can be sent.
const (
	// CmdOn represents command for lighting the stands
	CmdOn byte = iota
	// CmdOff represents command for disabling lights on stands
	CmdOff
	// CmdAnimate represents command for lighting stand led by led
	CmdAnimate
	// CmdSensor represents command for getting sensor data
	CmdSensor
	// CmdFade represents command for slowly lighting the stand
	CmdFade
	// CmdOffAll represents command for disabling lights on all stands
	CmdOffAll
)

const topic = "bullshelves"

// Config for MQTT supporting envconfig
// for environmental variable parsing.
// If initialized correctly (by envconfig), contains sane defaults.
type Config struct {
	MQTT struct {
		Broker string `envconfig:"default=tcp://test.mosquitto.org:1883"`

		KeepAlive struct {
			Seconds time.Duration `envconfig:"default=2s"`
		}

		Timeout struct {
			Seconds time.Duration `envconfig:"default=2s"`
		}

		Disconnect struct {
			Milliseconds time.Duration `envconfig:"default=250ms"`
		}
	}
}

// KeepAlive is a functional option for creating MQTT
// clients.
func KeepAlive(sec time.Duration) func(*Config) {
	return func(conf *Config) {
		conf.MQTT.KeepAlive.Seconds = sec
	}
}

// Timeout is a functional option for creating MQTT
// clients.
func Timeout(sec time.Duration) func(*Config) {
	return func(conf *Config) {
		conf.MQTT.Timeout.Seconds = sec
	}
}

// DisconnectTimeout is a functional option for creating MQTT
// clients.
func DisconnectTimeout(sec time.Duration) func(*Config) {
	return func(conf *Config) {
		conf.MQTT.Disconnect.Milliseconds = sec
	}
}

// Command describes command RPI should execute on slabs.
type Command struct {
	Platform byte
	Command  byte
	RGB      *RGB
}

// RGB describes color with 3 bytes.
// Each color (r, g, b) fits between 0 and 255.
type RGB struct {
	R, G, B byte
}

// MQTT is an abstraction over MQTT client
// with API specific for Bullseye slab control.
type MQTT struct {
	client mqtt.Client
	topic  string
	evbus.Bus
}

func defaultHandler(bus evbus.Bus) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("TOPIC: %s\n", msg.Topic())
		log.Printf("MSG: %s\n", msg.Payload())
		cmd := new(Command)

		err := json.Unmarshal(msg.Payload(), cmd)
		if err != nil {
			log.Fatalln(err)
			return
		}

		bus.Publish(topic, cmd)
	}
}

// New creates and connects MQTT client (by calling FromConfig).
func New(broker string, options ...func(*Config)) (*MQTT, error) {
	conf := new(Config)
	for _, opt := range options {
		opt(conf)
	}

	conf.MQTT.Broker = broker

	return FromConfig(conf)
}

// FromConfig initializes MQTT client based on config,
// as opposed to New, which initializes based on functional options.
func FromConfig(config *Config) (*MQTT, error) {
	cli := MQTT{
		topic: topic,
		Bus:   evbus.New(),
	}

	opts := mqtt.NewClientOptions().AddBroker(config.MQTT.Broker)
	opts.SetKeepAlive(config.MQTT.KeepAlive.Seconds * time.Second)
	opts.SetDefaultPublishHandler(defaultHandler(cli.Bus))
	opts.SetPingTimeout(config.MQTT.Timeout.Seconds * time.Second)

	cli.client = mqtt.NewClient(opts)

	if token := cli.client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &cli, cli.subscribe(cli.topic)
}

// Subscribe mqtt topic. From now on messages sent on this topic
// will be broadcast through event bus.
func (m *MQTT) subscribe(topic string) error {
	if token := m.client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (m *MQTT) unsubscribe(topic string) error {
	if token := m.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// Subscribe on event bus using passed Command.
func (m *MQTT) Subscribe(fn func(*Command)) error {
	return m.Bus.Subscribe(topic, fn)
}

// Publish message to the broker.
func (m *MQTT) Publish(cmd Command) error {
	msg, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	log.Println("Sending ", m.topic, ": ", string(msg))
	token := m.client.Publish(m.topic, 0, false, msg)
	token.Wait()

	return nil
}

// Disconnect from the server.
func (m *MQTT) Disconnect(milliseconds time.Duration) error {
	m.client.Disconnect(uint(milliseconds))
	return nil
}
