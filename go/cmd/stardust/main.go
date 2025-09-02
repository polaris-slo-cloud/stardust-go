package main

import (
	"fmt"
	"log"
	"os"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/computing"
	"github.com/keniack/stardustGo/internal/deployment"
	"github.com/keniack/stardustGo/internal/routing"
	"github.com/keniack/stardustGo/internal/satellite"
	"github.com/keniack/stardustGo/internal/simulation"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <configFile>", os.Args[0])
	}
	configFile := os.Args[1]

	// Step 1: Load application configuration (from configs/appsettings.json)
	cfg, err := configs.LoadConfigFromFile(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Step 2: Build computing builder with configured strategies
	computingBuilder := computing.NewComputingBuilder(cfg.Computing[0])

	// Step 3: Build router builder
	routerBuilder := routing.NewRouterBuilder(cfg.Router)

	// Step 4: Initialize the satellite builder
	satBuilder := satellite.NewSatelliteBuilder(routerBuilder, computingBuilder, cfg.ISL)
	tleLoader := satellite.NewTleLoader(cfg.ISL, satBuilder)

	// Step 4.1: Initialize constellation loader and register TLE loader
	constellationLoader := satellite.NewSatelliteConstellationLoader()
	constellationLoader.RegisterDataSourceLoader("tle", tleLoader)

	// Step 5: Initialize simulation service
	simService := simulation.NewSimulationService(cfg.Simulation, routerBuilder, computingBuilder)

	// Step 6: Inject orchestrator (if used)
	orchestrator := deployment.NewDeploymentOrchestrator()
	simService.Inject(orchestrator)

	// Step 7: Load satellites using the loader service
	loaderService := satellite.NewLoaderService(cfg.ISL, satBuilder, constellationLoader, simService, fmt.Sprintf("./resources/%s/%s", cfg.Simulation.SatelliteDataSourceType, cfg.Simulation.SatelliteDataSource), cfg.Simulation.SatelliteDataSourceType)
	if err := loaderService.Start(); err != nil {
		log.Fatalf("Failed to load satellites: %v", err)
	}

	// Step 8: Start the simulation loop or run individual code
	if cfg.Simulation.StepInterval >= 0 {
		done := simService.StartAutorun()
		<-done // blocks main goroutine until simulation stops
	} else {
		log.Println("Simulation loaded. Not autorunning as StepInterval < 0.")
		simService.StepBySeconds(60) // Example: step by 60 seconds
		log.Println("Simulation stepped by 60 seconds.")
	}
}
