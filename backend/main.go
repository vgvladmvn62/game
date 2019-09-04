package main

import (
	"fmt"
	"os"

	"github.com/vrischmann/envconfig"
	"github.wdf.sap.corp/Magikarpie/bullseye/internal/db/products"
	"github.wdf.sap.corp/Magikarpie/bullseye/internal/productcache"
	"github.wdf.sap.corp/Magikarpie/bullseye/internal/stand"
	"github.wdf.sap.corp/Magikarpie/bullseye/pkg/mqtt"

	"github.wdf.sap.corp/Magikarpie/bullseye/internal/db"
	"github.wdf.sap.corp/Magikarpie/bullseye/internal/db/attributes"
	"github.wdf.sap.corp/Magikarpie/bullseye/internal/db/questions"
	"github.wdf.sap.corp/Magikarpie/bullseye/internal/db/stands"
	"github.wdf.sap.corp/Magikarpie/bullseye/internal/ec"
	"github.wdf.sap.corp/Magikarpie/bullseye/internal/hardware"
	"github.wdf.sap.corp/Magikarpie/bullseye/internal/matching"
	"github.wdf.sap.corp/Magikarpie/bullseye/internal/server"
)

// Config stores entire application configuration.
type Config struct {
	EC     ec.Config
	Server server.Config
	DB     db.Config
	HW     mqtt.Config
}

func main() {
	var config Config

	err := envconfig.Init(&config)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(config)

	database, err := db.New(&config.DB)

	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
	}

	defer func() {
		err = client.Disconnect(config.HW.MQTT.Disconnect.Milliseconds)
		if err != nil {
			fmt.Println(err)
		}
	}()

	hardwareService := hardware.NewHardwareService(client, standsRepo)

	matcher := matching.NewService(standService, attributesRepo)

	srv := server.NewServer(&config.Server, productCacheService, standsRepo, standService, matcher,
		attributesRepo, hardwareService, questionsRepo)

	err = srv.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
