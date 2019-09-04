package slab

import (
	"log"
	"sync"
	"time"

	evbus "github.com/asaskevich/EventBus"
	"github.com/boltdb/bolt"
)

var (
	green         = &RGB{0, 255, 0}
	orange        = &RGB{252, 169, 3}
	bucketSlabIDs = []byte("slab_ids")
)

type dbR interface {
	View(func(*bolt.Tx) error) error
}

type dbW interface {
	Update(func(*bolt.Tx) error) error
}

// DB is an abstract for accessing database.
type DB interface {
	dbR
	dbW
	Close() error
}

// OpenDB from given path.
func OpenDB(path string) (DB, error) {
	return bolt.Open(path, 0777, nil)
}

// AssignIDs of given number of slabs available from event bus
// in database.
func AssignIDs(bus evbus.Bus, number int, db dbW) {
	assigned := make([]string, 0, 10)
	var wg sync.WaitGroup
	wg.Add(number)
	log.Println("Assigning")
	callback := func(e Event) {
		id := e.Slab.ID()
		log.Println("Processing ", id)
		if slabIDin(id, assigned) {
			return
		}
		log.Println("Assigning: ", e.Slab, " -> ", len(assigned))
		assigned = append(assigned, e.Slab.ID())
		e.Slab.On(orange)
		time.Sleep(700 * time.Millisecond)
		e.Slab.On(green)
		wg.Done()
	}

	err := bus.Subscribe("slab", callback)
	if err != nil {
		log.Fatalln(err)
	}
	wg.Wait()
	bus.Unsubscribe("slab", callback)
	time.Sleep(600 * time.Millisecond)
	pushSlabsToDB(assigned, db)
	log.Println("Assigned all")
}

func pushSlabsToDB(slabs []string, db dbW) {
	db.Update(func(tx *bolt.Tx) error {
		// Error with deleting probably means there was nothing to delete,
		// which is fine
		_ = tx.DeleteBucket(bucketSlabIDs)
		bucket, err := tx.CreateBucket(bucketSlabIDs)
		if err != nil {
			log.Fatalln("Didn't create bucket: ", err)
			return err
		}

		for num, physical := range slabs {
			bucket.Put([]byte{byte(num)}, []byte(physical))
		}

		return nil
	})
}

func readSlabsFromDB(db dbR) []string {
	bindings := make([]string, 0, 50) // FIXME: This limits shelf number
	_ = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketSlabIDs)
		if bucket == nil {
			return nil
		}
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			i := k[0]
			log.Println("i: ", i, " v: ", string(v))
			bindings = append(bindings, string(v))
		}
		return nil
	})

	return bindings
}

func slabIDin(a string, list []string) bool {
	for _, b := range list {
		if a == b {
			return true
		}
	}

	return false
}
