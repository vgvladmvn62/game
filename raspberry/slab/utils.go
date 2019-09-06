package slab

import (
	"log"
	"sync"

	evbus "github.com/asaskevich/EventBus"
)

var testColor = &RGB{237, 17, 230}

// TestLights turns off each slab, then turns them back off after a while.
func TestLights(slabs []*Slab) error {
	var err error
	for _, s := range slabs {
		err = s.On(testColor)
		if err != nil {
			return err
		}
	}

	return nil
}

// Event is returned by object listener.
type Event struct {
	Slab   *Slab
	Sensor bool
}

// ObjectListener emits event when sensor status changes on a slab.
// It can signal detection with lighting slab.
func ObjectListener(s *Slab, bus evbus.Bus) error {
	prev := false
	state := false
	var err error
	for {
		state, err = s.Sensor()
		if err != nil {
			return err
		}

		if state != prev {
			bus.Publish("slab", Event{s, state})
			prev = state
		}
	}
}

// ListenAll starts ObjectListener in parallel on all slabs,
// sending events to event bus.
func ListenAll(slabs []*Slab, bus evbus.Bus, wg *sync.WaitGroup) {
	for _, s := range slabs {
		go func(s *Slab, wg *sync.WaitGroup) {
			id := s.ID()
			_ = ObjectListener(s, bus)
			log.Println("Disconnected ", id, " wg: ", wg)
			if wg != nil {
				wg.Done()
			}
		}(s, wg)
	}
}

// SignalObject light a slab when an object appears.
func SignalObject(bus evbus.Bus, c *RGB) error {
	return bus.Subscribe("slab", func(e Event) {
		if e.Sensor {
			_ = e.Slab.On(c)
		} else {
			_ = e.Slab.Off()
		}
	})
}

// CloseAll closes all slab ports.
func CloseAll(slabs []*Slab) {
	for _, s := range slabs {
		_ = s.Close()
	}
}

// DiscardEvents returns channel that discards events sent to it.
func DiscardEvents() chan Event {
	ch := make(chan Event, 20)
	go func() {
		for range ch {
		}
	}()

	return ch
}
