package slab

import (
	"io/ioutil"
	"log"
	"strings"
	"sync"

	evbus "github.com/asaskevich/EventBus"
	"github.com/boltdb/bolt"
	"github.com/kyma-incubator/bullseye-showcase/backend/pkg/mqtt"
)

// Repo stores slabs' data. It allows controlling
// multiple slabs at once.
type Repo struct {
	slabs    []*Slab
	bindings []string
	db       DB
	wg       *sync.WaitGroup

	evbus.Bus
}

// Close closes all managed slabs and database.
func (r *Repo) Close() error {
	for _, s := range r.slabs {

		err := s.Close()
		if err != nil {
			return err
		}
	}

	if r.db != nil {
		err := r.db.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// NewRepo creates Repo of all slabs,
// iterating through /dev files and connecting to those, containing `filter`.
func NewRepo(filter string) (r Repo, err error) {
	files, err := ioutil.ReadDir("/dev")
	if err != nil {
		return
	}

	r.slabs = make([]*Slab, 0, 20)
	r.Bus = evbus.New()
	r.wg = new(sync.WaitGroup)

	for _, f := range files {
		name := f.Name()
		if strings.Contains(name, filter) {
			path := "/dev/" + name

			s, err := New(path)
			if err != nil {
				continue
			}

			raw, _ := s.RawSensor()

			err = s.SetThreshold(raw / 2)
			if err != nil {
				continue
			}

			r.slabs = append(r.slabs, s)
			log.Print(path, " -> ", s.ID())
		}
	}

	log.Println("Creating repo with ", len(r.slabs), " slabs")

	r.wg.Add(len(r.slabs))
	ListenAll(r.slabs, r, r.wg)

	return
}

// WaitGroup waits for a collection of goroutines to finish.
// Proxy method.
func (r *Repo) WaitGroup() *sync.WaitGroup {
	return r.wg
}

// OpenDB from given path.
func (r *Repo) OpenDB(path string) error {
	db, err := bolt.Open(path, 0777, nil)
	if err != nil {
		return err
	}

	r.db = db

	log.Println("Opened db: ", db)

	return nil
}

// OffAll disables all lights.
func (r *Repo) OffAll() error {
	for _, s := range r.slabs {
		if err := s.Off(); err != nil {
			return err
		}
	}

	return nil
}

// Slab finds slab by ID assigned from database
func (r *Repo) Slab(id byte) *Slab {
	realID := r.bindings[id]

	for _, s := range r.slabs {
		if realID == s.ID() {
			return s
		}
	}

	return nil
}

// AssignIDs to all slabs.
func (r *Repo) AssignIDs() {
	AssignIDs(r, len(r.slabs), r.db)
	err := r.OffAll()
	if err != nil {
		log.Println(err)
		return
	}

	r.LoadIDs()
}

// LoadIDs of all slabs.
func (r *Repo) LoadIDs() {
	log.Println("Loading IDs")
	r.bindings = readSlabsFromDB(r.db)
}

// LoadOrAssign IDs of all slabs.
func (r *Repo) LoadOrAssign() {
	r.LoadIDs()
	if len(r.bindings) < len(r.slabs) {
		r.AssignIDs()
	}
}

func (r *Repo) idOf(s *Slab) (int, bool) {
	id := s.ID()
	for i, b := range r.bindings {
		if b == id {
			return i, true
		}
	}
	return 0, false
}

// SignalObject with given color represented as RGB.
func (r *Repo) SignalObject(color *RGB) {
	err := r.Subscribe("slab", func(e Event) {
		id, _ := r.idOf(e.Slab)
		log.Println("Slab ", id, " hasObject? ", e.Sensor)
		if e.Sensor {
			_ = e.Slab.On(color)
		} else {
			_ = e.Slab.Off()
		}
	})

	if err != nil {
		return
	}
}

func importCmd(in *mqtt.Command) CommandDTO {
	var rgb *RGB
	if in.RGB != nil {
		rgb = &RGB{in.RGB.R, in.RGB.R, in.RGB.B}
	}
	return CommandDTO{
		ID:      in.Platform,
		Command: in.Command,
		RGB:     rgb,
		Delay:   0,
	}
}

// Execute given MQTT command.
func (r *Repo) Execute(cmdAlien *mqtt.Command) {
	cmd := importCmd(cmdAlien)

	switch cmd.Command {
	case mqttOffAll:
		log.Println("Turning off all")
		err := r.OffAll()
		if err != nil {
			return
		}
	case mqttFade:
		slab := r.Slab(cmd.ID)
		log.Println("Fading in: ", slab.ID())
		err := slab.SetBrightness(50)
		if err != nil {
			return
		}

		err = slab.On(cmd.RGB)
		if err != nil {
			return
		}
	}
}
