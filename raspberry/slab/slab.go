package slab

import (
	"io"
	"sync"
	"time"

	"github.com/jacobsa/go-serial/serial"
)

const (
	cmdID byte = iota
	cmdOn
	cmdOFF
	cmdAnimate
	cmdSensor
	cmdRawSensor
	cmdSetBrightness
	cmdSetThreshold
	cmdSetID
)

// RGB stores color info for slabs.
type RGB struct {
	R, G, B byte
}

// Slab contains single slab info.
type Slab struct {
	lock sync.Mutex
	port io.ReadWriteCloser
	name string
}

// New creates new slab connected to usb port `portPath`.
func New(portPath string) (*Slab, error) {
	options := serial.OpenOptions{
		PortName:        portPath,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	s := Slab{}
	s.name = portPath
	var err error

	s.port, err = serial.Open(options)

	val := make(chan bool)
	go func(port io.Reader) {
		buf := make([]byte, 20)
		s.port.Read(buf)
		val <- true
	}(s.port)

	select {
	case <-val:
		return &s, err
	case <-time.After(500 * time.Millisecond):
		return &s, err
	}

}

// ID is a unique string identifying slab. Currently /dev/... path.
func (s *Slab) ID() string {
	return s.name
}

// SetThreshold sets value below which object is detected as 'on slab'.
func (s *Slab) SetThreshold(threshold byte) error {
	s.lock.Lock()
	_, err := s.port.Write([]byte{cmdSetThreshold, threshold})
	if err != nil {
		return err
	}
	s.lock.Unlock()

	return nil
}

// On turns on slab with given color.
func (s *Slab) On(c *RGB) error {
	s.lock.Lock()
	_, err := s.port.Write([]byte{cmdOn, c.R, c.G, c.B})
	s.lock.Unlock()
	return err
}

// Off turns off the light.
func (s *Slab) Off() error {
	s.lock.Lock()
	_, err := s.port.Write([]byte{cmdOFF})
	s.lock.Unlock()

	return err
}

// FadeIn slowly brightens the slab with given color.
func (s *Slab) FadeIn(delay time.Duration, upto byte, c *RGB) error {
	s.lock.Lock()
	err := s.SetBrightness(0)
	if err != nil {
		return err
	}
	err = s.On(c)
	if err != nil {
		return err
	}

	for b := byte(0); b < upto; b++ {
		err = s.SetBrightness(b)
		if err != nil {
			return err
		}
		err = s.On(c)
		if err != nil {
			return err
		}

		time.Sleep(delay)
	}

	s.lock.Unlock()

	return nil
}

// Sensor asks slab if something is on it.
func (s *Slab) Sensor() (bool, error) {
	s.lock.Lock()
	_, err := s.port.Write([]byte{cmdSensor})
	if err != nil {
		return false, err
	}

	buff := make([]byte, 1)
	_, err = s.port.Read(buff)
	if err != nil {
		return false, err
	}
	s.lock.Unlock()

	return buff[0] == 1, nil
}

// RawSensor asks slab for raw sensor data.
func (s *Slab) RawSensor() (byte, error) {
	s.lock.Lock()
	_, err := s.port.Write([]byte{cmdRawSensor})
	if err != nil {
		return 0, err
	}

	buff := make([]byte, 1)
	_, err = s.port.Read(buff)
	if err != nil {
		return 0, err
	}

	s.lock.Unlock()

	return buff[0], nil
}

// SetBrightness set brightness, which isn't updated until next color set.
func (s *Slab) SetBrightness(value byte) error {
	s.lock.Lock()
	_, err := s.port.Write([]byte{cmdSetBrightness, value})
	s.lock.Unlock()
	return err
}

// Close closes the port used for slab communication.
func (s *Slab) Close() error {
	return s.port.Close()
}
