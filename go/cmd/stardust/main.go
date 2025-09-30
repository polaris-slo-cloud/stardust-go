package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

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
	simulationConfigString := flag.String("simulationConfig", "", "Path to the simulation config file")
	islConfigString := flag.String("islConfig", "", "Path to inter satellite link config file")
	groundLinkConfigString := flag.String("groundLinkConfig", "", "Path to ground link config file")
	computingConfigString := flag.String("computingConfig", "", "Path to computing config file")
	routerConfigString := flag.String("routerConfig", "", "Path to router config file")
	simulationStateOutputFile := flag.String("simulationStateOutputFile", "", "Path to output the simulation state (optional)")
	simulationStateInputFile := flag.String("simulationStateInputFile", "", "Path to input the simulation state (optional)")
	simulationPluginString := flag.String("simulationPlugins", "", "Plugin names (optional, comma-separated list)")
	statePluginString := flag.String("statePlugins", "", "Plugin names (optional, comma-seperated list)")

	flag.Parse()

	simulationPluginList := strings.Split(*simulationPluginString, ",")
	statePluginList := strings.Split(*statePluginString, ",")

	// Step 1: Load configuration
	simulationConfig, err := configs.LoadConfigFromFile[configs.SimulationConfig](*simulationConfigString)
	if err != nil {
		log.Fatalf("Failed to load simulation configuration: %v", err)
	}

	computingConfig, err := configs.LoadConfigFromFile[[]configs.ComputingConfig](*computingConfigString)
	if err != nil {
		log.Fatalf("Failed to load simulation configuration: %v", err)
	}

	routerConfig, err := configs.LoadConfigFromFile[configs.RouterConfig](*routerConfigString)
	if err != nil {
		log.Fatalf("Failed to load simulation configuration: %v", err)
	}

	var simService types.SimulationController
	if *simulationStateInputFile != "" {
		simService = startSimulationIteration(*simulationConfig, *computingConfig, *routerConfig, *simulationStateInputFile, simulationPluginList)
	} else {
		simService = startSimulation(*simulationConfig, *islConfigString, *groundLinkConfigString, *computingConfig, *routerConfig, simulationStateOutputFile, simulationPluginList, statePluginList)
	}

	myCode(simService, *simulationConfig)
}

func startSimulationIteration(simulationConfig configs.SimulationConfig, computingConfig []configs.ComputingConfig, routerConfig configs.RouterConfig, simulationStateInputFile string, simulationPluginList []string) types.SimulationController {
	// Step 2: Build computing builder with configured strategies
	var computingBuilder computing.ComputingBuilder = computing.NewComputingBuilder(computingConfig)

	// Step 3: Build router builder
	routerBuilder := routing.NewRouterBuilder(routerConfig)

	// Step 4.1: Initialize plugin builder
	simPluginBuilder := simplugin.NewPluginBuilder()
	simPlugins, err := simPluginBuilder.BuildPlugins(simulationPluginList)
	if err != nil {
		log.Fatalf("Failed to build simualtion plugins: %v", err)
		return nil
	}

	// Step 5: State Plugin Builder
	statePluginBuilder := stateplugin.NewStatePluginPrecompBuilder(simulationStateInputFile)

	// Step 6: Inject orchestrator (if used)
	orchestrator := deployment.NewDeploymentOrchestrator()

	simStateDeserializer := simulation.NewSimulationStateDeserializer(&simulationConfig, simulationStateInputFile, computingBuilder, routerBuilder, orchestrator, simPlugins, statePluginBuilder)
	return simStateDeserializer.LoadIterator()
}

func startSimulation(simulationConfig configs.SimulationConfig, islConfigString string, groundLinkConfigString string, computingConfig []configs.ComputingConfig, routerConfig configs.RouterConfig, simulationStateOutputFile *string, simulationPluginList []string, statePluginList []string) types.SimulationController {
	islConfig, err := configs.LoadConfigFromFile[configs.InterSatelliteLinkConfig](islConfigString)
	if err != nil {
		log.Fatalf("Failed to load isl configuration: %v", err)
	}

	groundLinkConfig, err := configs.LoadConfigFromFile[configs.GroundLinkConfig](groundLinkConfigString)
	if err != nil {
		log.Fatalf("Failed to load isl configuration: %v", err)
	}

	// Step 2: Build computing builder with configured strategies
	computingBuilder := computing.NewComputingBuilder(computingConfig)

	// Step 3: Build router builder
	routerBuilder := routing.NewRouterBuilder(routerConfig)

	// Step 4.1: Initialize plugin builder
	simPluginBuilder := simplugin.NewPluginBuilder()
	simPlugins, err := simPluginBuilder.BuildPlugins(simulationPluginList)
	if err != nil {
		log.Fatalf("Failed to build simualtion plugins: %v", err)
		return nil
	}

	// Step 4.2: Initialize state plugin builder
	statePluginBuilder := stateplugin.NewStatePluginBuilder()
	statePlugins, err := statePluginBuilder.BuildPlugins(statePluginList)
	if err != nil {
		log.Fatalf("Failed to build state plugins: %v", err)
		return nil
	}

	// Step 5.1: Initialize the satellite builder
	satBuilder := satellite.NewSatelliteBuilder(routerBuilder, computingBuilder, *islConfig)
	tleLoader := satellite.NewTleLoader(*islConfig, satBuilder)

	// Step 4.2: Initialize the ground station loader
	groundStationBuilder := ground.NewGroundStationBuilder(simulationConfig.SimulationStartTime, routerBuilder, computingBuilder, *groundLinkConfig)
	ymlLoader := ground.NewGroundStationYmlLoader(*groundLinkConfig, groundStationBuilder)

	// Step 4.3: Initialize constellation loader and register TLE loader
	constellationLoader := satellite.NewSatelliteConstellationLoader()
	constellationLoader.RegisterDataSourceLoader("tle", tleLoader)

	// Step 5: Initialize simulation service
	simService := simulation.NewSimulationService(&simulationConfig, routerBuilder, computingBuilder, simPlugins, types.NewStatePluginRepository(statePlugins), simulationStateOutputFile)

	// Step 6: Inject orchestrator (if used)
	orchestrator := deployment.NewDeploymentOrchestrator()
	simService.Inject(orchestrator)

	// Step 8: Load satellites using the loader service
	loaderService := satellite.NewSatelliteLoaderService(*islConfig, satBuilder, constellationLoader, simService, fmt.Sprintf("./resources/%s/%s", simulationConfig.SatelliteDataSourceType, simulationConfig.SatelliteDataSource), simulationConfig.SatelliteDataSourceType)
	if err := loaderService.Start(); err != nil {
		log.Fatalf("Failed to load satellites: %v", err)
	}

	// Step 9: Load ground stations using the ground station loader service
	groundLoaderService := ground.NewGroundStationLoaderService(simService, groundStationBuilder, ymlLoader, fmt.Sprintf("./resources/%s/%s", simulationConfig.GroundStationDataSourceType, simulationConfig.GroundStationDataSource), simulationConfig.GroundStationDataSourceType)
	if err := groundLoaderService.Start(); err != nil {
		log.Fatalf("Failed to load ground stations: %v", err)
	}

	return simService
}

func myCode(simulationController types.SimulationController, simulationConfig configs.SimulationConfig) {
	defer simulationController.Close()

	// Start the simulation loop or run individual code
	if simulationConfig.StepInterval >= 0 {
		done := simulationController.StartAutorun()
		<-done // blocks main goroutine until simulation stops
	} else {
		log.Println("Simulation loaded. Not autorunning as StepInterval < 0.")
		for range simulationConfig.StepCount {
			simulationController.StepBySeconds(60 * 10) // Example: step by 60 seconds
			var sats = simulationController.GetGroundStations()
			var ground1 = sats[0]
			var ground2 = sats[80]
			var l1 = ground1.GetLinkNodeProtocol().Established()[0]
			var l2 = ground2.GetLinkNodeProtocol().Established()[0]
			var uplinkSat1 = l1.GetOther(ground1)
			var uplinkSat2 = l2.GetOther(ground2)
			var route, err = ground1.GetRouter().RouteToNode(ground2, nil)
			var interSatelliteRoute, _ = uplinkSat1.GetRouter().RouteToNode(uplinkSat2, nil)
			if err != nil {
				log.Println("Routing error:", err)
			} else {
				log.Println("Route from", ground1.GetName(), "to", ground2.GetName(), "in", route.Latency(), "ms")
				log.Println("Uplink latency", l1.Latency()+l2.Latency(), "ms")
				log.Println("Latency between uplink nodes:", interSatelliteRoute.Latency(), "ms")
				log.Println(ground1.GetName(), "->", uplinkSat1.GetName(), "->", uplinkSat2.GetName(), "->", ground2.GetName())
				log.Println(l1.Distance(), "->", uplinkSat1.DistanceTo(uplinkSat2), "->", l2.Distance())
				log.Println(l1.Latency(), "->", interSatelliteRoute.Latency(), "->", l2.Latency())
				log.Println(uplinkSat1.DistanceTo(uplinkSat2)/1000, "km apart")
				log.Println(uplinkSat1.GetPosition(), uplinkSat2.GetPosition())
			}
			log.Println(len(sats), "satellites in simulation.")
			log.Println("Simulation stepped by 60 seconds.")

			var statePlugin = types.GetStatePlugin[stateplugin.SunStatePlugin](simulationController.GetStatePluginRepository())
			log.Println("Sunlight exposure of", uplinkSat1.GetName(), "is", statePlugin.GetSunlightExposure(uplinkSat1))
		}
	}
}
