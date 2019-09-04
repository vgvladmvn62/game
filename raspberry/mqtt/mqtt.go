package mqtt

import (
	"encoding/json"
	"log"
	"time"

	evbus "github.com/asaskevich/EventBus"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	CmdOn byte = iota
	CmdOff
	CmdAnimate
	CmdSensor
	CmdFade
	CmdOffAll
)

// KeepAlive is a functional option for creating MQTT
// clients.
func KeepAlive(sec time.Duration) func(*Config) {
	return func(conf *Config) {
		conf.KeepAlive.Seconds = int(sec)
	}
}

// Timeout is a functional option for creating MQTT
// clients.
func Timeout(sec time.Duration) func(*Config) {
	return func(conf *Config) {
		conf.Timeout.Seconds = int(sec)
	}
}

// DisconnectTimeout is a functional option for creating MQTT
// clients.
func DisconnectTimeout(sec time.Duration) func(*Config) {
	return func(conf *Config) {
		conf.Disconnect.Milliseconds = int(sec)
	}
}

// Command describes command that RPI should execute on slabs.
type Command struct {
	Platform byte
	Command  byte
	RGB      *RGB
}

// RGB describes color with 3 bytes.
// Each color (R, G, B) fits between 0 and 255.
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

func defaultHandler(bus evbus.Bus, topic string) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("TOPIC: %s\n", msg.Topic())
		log.Printf("MSG: %s\n", msg.Payload())
		cmd := new(Command)
		err := json.Unmarshal(msg.Payload(), cmd)
		if err != nil {
			log.Println(err)
			return
		}

		bus.Publish(topic, cmd)
	}
}

// New creates an instance of MQTT client.
func New(broker string, options ...func(*Config)) (*MQTT, error) {
	conf := new(Config)
	for _, opt := range options {
		opt(conf)
	}

	conf.Broker = broker

	return FromConfig(conf)
}

// FromConfig initializes MQTT client based on config,
// as opposed to New, which initializes based on functional options
func FromConfig(config *Config) (*MQTT, error) {
	cli := MQTT{
		topic: config.Topic,
		Bus:   evbus.New(),
	}

	opts := mqtt.NewClientOptions().AddBroker(config.Broker)
	opts.SetKeepAlive(time.Duration(config.KeepAlive.Seconds) * time.Second)
	opts.SetDefaultPublishHandler(defaultHandler(cli.Bus, cli.topic))
	opts.SetPingTimeout(time.Duration(config.Timeout.Seconds) * time.Second)

	cli.client = mqtt.NewClient(opts)

	if token := cli.client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &cli, cli.subscribe(cli.topic)
}

// Subscribe MQTT topic. Messages that are being sent on this topic
// are broadcast through event bus.
func (m *MQTT) subscribe(topic string) error {
	if token := m.client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// Subscribe to Bullseye backend events. Messages sent on it's topic
// are broadcast through event bus.
func (m *MQTT) Subscribe(fn func(*Command)) error {
	return m.Bus.Subscribe(m.topic, fn)
}

func (m *MQTT) unsubscribe(topic string) error {
	if token := m.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
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
func (m *MQTT) Disconnect(milliseconds int) error {
	m.client.Disconnect(uint(milliseconds))
	return nil
}
