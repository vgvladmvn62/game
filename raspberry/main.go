package main

import (
	"log"

	"github.com/vrischmann/envconfig"

	"github.wdf.sap.corp/team-wojciech/bullseye-rpi/mqtt"
	"github.wdf.sap.corp/team-wojciech/bullseye-rpi/slab"
)

// Config stores entire application configuration.
type Config struct {
	MQTT mqtt.Config
	Slab slab.Config
}

func main() {
	var config Config

	err := envconfig.Init(&config)
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println("Config: ", config)

	repository, err := slab.NewRepo(config.Slab.Repository.Filter)
	if err != nil {
		log.Fatalln(err)
		return
	}

	defer func() {
		err = repository.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	repository.OffAll()
	err = repository.OpenDB(config.Slab.Repository.Path)
	if err != nil {
		log.Fatalln(err)
		return
	}

	repository.LoadOrAssign()

	cli, err := mqtt.FromConfig(&config.MQTT)
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = cli.Subscribe(repository.Execute)
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println("Waiting for actions...")
	repository.WaitGroup().Wait()
}
