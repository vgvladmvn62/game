package main

import (
	"log"
	"os"

	"github.com/vrischmann/envconfig"

	"github.com/kyma-incubator/bullseye-showcase/backend/pkg/mqtt"
	"github.com/kyma-incubator/bullseye-showcase/raspberry/slab"
)

// Config stores entire application configuration.
type Config struct {
	Slab slab.Config
	MQTT mqtt.Config
}

func main() {
	var config Config

	err := envconfig.Init(&config)
	if err != nil {
		log.Fatalln(err)
		return
	}

	if len(os.Args) == 2 {
		config.MQTT.MQTT.Broker = os.Args[1]
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

	err = repository.OffAll()
	if err != nil {
		log.Fatalln(err)
		return
	}

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
