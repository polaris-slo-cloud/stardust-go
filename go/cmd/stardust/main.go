package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/computing"
	"github.com/keniack/stardustGo/internal/deployment"
	"github.com/keniack/stardustGo/internal/ground"
	"github.com/keniack/stardustGo/internal/routing"
	"github.com/keniack/stardustGo/internal/satellite"
	"github.com/keniack/stardustGo/internal/simplugin"
	"github.com/keniack/stardustGo/internal/simulation"
	"github.com/keniack/stardustGo/internal/stateplugin"
	"github.com/keniack/stardustGo/pkg/types"
)

func main() {

	configFile := flag.String("configFile", "", "Path to the configuration file")
	simulationStateOutputFile := flag.String("simulationStateOutputFile", "", "Path to output the simulation state (optional)")
	simulationStateInputFile := flag.String("simulationStateInputFile", "", "Path to input the simulation state (optional)")

	flag.Parse()

	if *configFile == "" {
		log.Fatal("--configFile missing")
	}

	// Step 1: Load application configuration (from configs/appsettings.json)
	cfg, err := configs.LoadConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	var simService types.SimulationController
	if *simulationStateInputFile != "" {
		simService = startSimulationIteration(cfg, *simulationStateInputFile)
	} else {
		simService = startSimulation(cfg, simulationStateOutputFile)
	}
	defer simService.Close()

	// Start the simulation loop or run individual code
	if cfg.Simulation.StepInterval >= 0 {
		done := simService.StartAutorun()
		<-done // blocks main goroutine until simulation stops
	} else {
		log.Println("Simulation loaded. Not autorunning as StepInterval < 0.")
		for range 2 {
			simService.StepBySeconds(60 * 10) // Example: step by 60 seconds
			var sats = simService.GetGroundStations()
			var ground1 = sats[0]
			var ground2 = sats[80]
			var l1 = ground1.GetLinkNodeProtocol().Established()[0]
			var l2 = ground2.GetLinkNodeProtocol().Established()[0]
			var uplinkSat1 = l1.GetOther(ground1)
			var uplinkSat2 = l2.GetOther(ground2)
			var route, err = ground1.GetRouter().RouteAsyncToNode(ground2, nil)
			var x, _ = uplinkSat1.GetRouter().RouteAsyncToNode(uplinkSat2, nil)
			if err != nil {
				log.Println("Routing error:", err)
			} else {
				log.Println(ground1.DistanceTo(ground2)/1000, "km apart")
				log.Println("Route from", ground1.GetName(), "to", ground2.GetName(), "in", route.Latency(), "ms")
				log.Println("Uplink latency", l1.Latency()+l2.Latency(), "ms")
				log.Println("Latency between uplink nodes:", x.Latency(), "ms")
				log.Println(uplinkSat1.GetName(), "->", l2.GetOther(ground2).GetName())
				log.Println(uplinkSat1.DistanceTo(uplinkSat2)/1000, "km apart")
				log.Println(uplinkSat1.GetPosition(), uplinkSat2.GetPosition())
			}
			log.Println(len(sats), "satellites in simulation.")
			log.Println("Simulation stepped by 60 seconds.")

			// var statePlugin = types.GetStatePlugin[*stateplugin.DummySunStatePlugin](simService.GetStatePluginRepository())
			// log.Println("Sunlight exposure of", uplinkSat1.GetName(), "is", statePlugin.GetSunlightExposure(uplinkSat1))
		}
	}
}

func startSimulationIteration(cfg *configs.Config, simulationStateInputFile string) types.SimulationController {
	// Step 2: Build computing builder with configured strategies
	var computingBuilder computing.ComputingBuilder = computing.NewComputingBuilder(cfg.Computing)

	// Step 3: Build router builder
	routerBuilder := routing.NewRouterBuilder(cfg.Router)

	// Step 4.1: Initialize plugin builder
	simPluginBuilder := simplugin.NewPluginBuilder()
	simPlugins, err := simPluginBuilder.BuildPlugins(cfg.Simulation.Plugins)
	if err != nil {
		log.Fatalf("Failed to build simualtion plugins: %v", err)
		return nil
	}

	// Step 6: Inject orchestrator (if used)
	orchestrator := deployment.NewDeploymentOrchestrator()

	simStateDeserializer := simulation.NewSimulationStateDeserializer(&cfg.Simulation, simulationStateInputFile, computingBuilder, routerBuilder, orchestrator, simPlugins)
	return simStateDeserializer.LoadIterator()
}

func startSimulation(cfg *configs.Config, simulationStateOutputFile *string) types.SimulationController {
	// Step 2: Build computing builder with configured strategies
	computingBuilder := computing.NewComputingBuilder(cfg.Computing)

	// Step 3: Build router builder
	routerBuilder := routing.NewRouterBuilder(cfg.Router)

	// Step 4.1: Initialize plugin builder
	simPluginBuilder := simplugin.NewPluginBuilder()
	simPlugins, err := simPluginBuilder.BuildPlugins(cfg.Simulation.Plugins)
	if err != nil {
		log.Fatalf("Failed to build simualtion plugins: %v", err)
		return nil
	}

	// Step 4.2: Initialize state plugin builder
	statePluginBuilder := stateplugin.NewStatePluginBuilder()
	statePlugins, err := statePluginBuilder.BuildPlugins(cfg.Simulation.Plugins)
	if err != nil {
		log.Fatalf("Failed to build state plugins: %v", err)
		return nil
	}

	// Step 5.1: Initialize the satellite builder
	satBuilder := satellite.NewSatelliteBuilder(routerBuilder, computingBuilder, cfg.ISL)
	tleLoader := satellite.NewTleLoader(cfg.ISL, satBuilder)

	// Step 4.2: Initialize the ground station loader
	groundStationBuilder := ground.NewGroundStationBuilder(cfg.Simulation.SimulationStartTime, routerBuilder, computingBuilder, cfg.Ground)
	ymlLoader := ground.NewGroundStationYmlLoader(cfg.Ground, groundStationBuilder)

	// Step 4.3: Initialize constellation loader and register TLE loader
	constellationLoader := satellite.NewSatelliteConstellationLoader()
	constellationLoader.RegisterDataSourceLoader("tle", tleLoader)

	// Step 5: Initialize simulation service
	simService := simulation.NewSimulationService(&cfg.Simulation, routerBuilder, computingBuilder, simPlugins, types.NewStatePluginRepository(statePlugins), simulationStateOutputFile)

	// Step 6: Inject orchestrator (if used)
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

	return simService
}
