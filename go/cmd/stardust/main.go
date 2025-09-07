package main

import (
	"fmt"
	"log"
	"os"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/computing"
	"github.com/keniack/stardustGo/internal/deployment"
	"github.com/keniack/stardustGo/internal/ground"
	"github.com/keniack/stardustGo/internal/plugin"
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
	computingBuilder := computing.NewComputingBuilder(cfg.Computing)

	// Step 3: Build router builder
	routerBuilder := routing.NewRouterBuilder(cfg.Router)

	// Step 4: Initialize plugin builder
	pluginBuilder := plugin.NewPluginBuilder()
	plugins, err := pluginBuilder.BuildPlugins(cfg.Simulation.Plugins)
	if err != nil {
		log.Fatalf("Failed to build plugins: %v", err)
		return
	}

	// Step 5.1: Initialize the satellite builder
	satBuilder := satellite.NewSatelliteBuilder(routerBuilder, computingBuilder, cfg.ISL)
	tleLoader := satellite.NewTleLoader(cfg.ISL, satBuilder)

	// Step 5.2: Initialize constellation loader and register TLE loader
	constellationLoader := satellite.NewSatelliteConstellationLoader()
	constellationLoader.RegisterDataSourceLoader("tle", tleLoader)

	// Step 5.3: Initialize the ground station loader
	groundStationBuilder := ground.NewGroundStationBuilder(cfg.Simulation.SimulationStartTime, routerBuilder, computingBuilder)
	ymlLoader := ground.NewGroundStationYmlLoader(cfg.Ground, groundStationBuilder)

	// Step 6: Initialize simulation service
	simService := simulation.NewSimulationService(cfg.Simulation, routerBuilder, computingBuilder, plugins)

	// Step 7: Inject orchestrator (if used)
	orchestrator := deployment.NewDeploymentOrchestrator()
	simService.Inject(orchestrator)

	// Step 8: Load satellites using the loader service
	loaderService := satellite.NewSatelliteLoaderService(cfg.ISL, satBuilder, constellationLoader, simService, fmt.Sprintf("./resources/%s/%s", cfg.Simulation.SatelliteDataSourceType, cfg.Simulation.SatelliteDataSource), cfg.Simulation.SatelliteDataSourceType)
	if err := loaderService.Start(); err != nil {
		log.Fatalf("Failed to load satellites: %v", err)
	}

	// Step 9: Load ground stations using the ground station loader service
	groundLoaderService := ground.NewGroundStationLoaderService(simService, groundStationBuilder, ymlLoader, fmt.Sprintf("./resources/%s/%s", cfg.Simulation.GroundStationDataSourceType, cfg.Simulation.GroundStationDataSource), cfg.Simulation.GroundStationDataSourceType)
	if err := groundLoaderService.Start(); err != nil {
		log.Fatalf("Failed to load ground stations: %v", err)
	}

	// Step 10: Start the simulation loop or run individual code
	if cfg.Simulation.StepInterval >= 0 {
		done := simService.StartAutorun()
		<-done // blocks main goroutine until simulation stops
	} else {
		log.Println("Simulation loaded. Not autorunning as StepInterval < 0.")
		simService.StepBySeconds(60) // Example: step by 60 seconds
		var sats = simService.GetGroundStations()
		log.Println(len(sats), "satellites in simulation.")
		log.Println("Simulation stepped by 60 seconds.")
	}
}
