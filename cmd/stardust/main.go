package main

import (
	"log"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/deployment"
	"github.com/keniack/stardustGo/internal/routing"
	"github.com/keniack/stardustGo/internal/satellite"
	"github.com/keniack/stardustGo/internal/simulation"
)

func main() {
	// Step 1: Load application configuration (from configs/appsettings.json)
	cfg, err := configs.LoadConfig("configs/appsettings.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Step 2: Build computing builder with configured strategies
	computingBuilder := simulation.NewComputingBuilder(cfg.Computing)

	// Step 3: Build router builder
	routerBuilder := routing.NewRouterBuilder(cfg.Router)

	// Step 4: Initialize the satellite builder and constellation loader
	satBuilder := satellite.NewSatelliteBuilder(routerBuilder, computingBuilder)
	constellationLoader := satellite.NewSatelliteConstellationLoader()
	tleLoader := satellite.NewTleLoader(cfg.ISL, satBuilder)
	constellationLoader.RegisterDataSourceLoader("tle", tleLoader)

	// Step 5: Initialize simulation service
	simService := simulation.NewSimulationService(cfg.Simulation, constellationLoader, routerBuilder, computingBuilder)

	// Step 6: Inject orchestrator (if used)
	orchestrator := deployment.NewDefaultDeploymentOrchestrator()
	simService.Inject(orchestrator)

	// Step 7: Load satellites using the loader service
	loaderService := satellite.NewLoaderService(cfg.ISL, satBuilder, constellationLoader, simService, "resources/tle/starlink_500.tle", "tle")
	if err := loaderService.Start(); err != nil {
		log.Fatalf("Failed to load satellites: %v", err)
	}

	// Step 8: Optionally start the simulation loop
	if cfg.Simulation.StepInterval >= 0 {
		simService.StartAutorunAsync()
	}

	// Block forever (simulate .NET HostedService lifetime)
	select {}
}
