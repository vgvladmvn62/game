package main

import (
	"log"
	"os"

	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db/products"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/productcache"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/stand"
	"github.com/kyma-incubator/bullseye-showcase/backend/pkg/mqtt"
	"github.com/vrischmann/envconfig"

	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db/attributes"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db/questions"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/db/stands"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/ec"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/hardware"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/matching"
	"github.com/kyma-incubator/bullseye-showcase/backend/internal/server"
)

// Config stores entire application configuration.
type Config struct {
	EC     ec.Config
	Server server.Config
	DB     db.Config
	HW     mqtt.Config
}

func main() {
	log.Println("Starting server...")

	var config Config

	err := envconfig.Init(&config)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Config:", config)

	database, err := db.New(&config.DB)
	if err != nil {
		log.Println(err)
	}

	standsRepo := stands.NewRepository(database)
	attributesRepo := attributes.NewRepository(database)
	questionsRepo := questions.NewRepository(database)
	productsRepo := products.NewRepository(database)

	productService := ec.NewProductService(&config.EC)
	productCacheService := productcache.NewProductCacheService(productService, standsRepo, productsRepo)

	standService := stand.NewStandService(productCacheService, standsRepo)

	client, err := mqtt.FromConfig(&config.HW)
	if err != nil {
		log.Println(err)
	}

	defer func() {
		err = client.Disconnect(config.HW.MQTT.Disconnect.Milliseconds)
		if err != nil {
			log.Println(err)
		}
	}()

	hardwareService := hardware.NewHardwareService(client, standsRepo)

	matcher := matching.NewService(standService, attributesRepo)

	srv := server.NewServer(&config.Server, productCacheService, standsRepo, standService, matcher,
		attributesRepo, hardwareService, questionsRepo)

	err = srv.Start()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
