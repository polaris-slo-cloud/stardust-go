package simulation

import (
	"log"
	"sync"
	"time"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/computing"
	"github.com/keniack/stardustGo/internal/routing"
	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.SimulationController = (*SimulationService)(nil)

// SimulationService handles simulation lifecycle and state updates
type SimulationService struct {
	BaseSimulationService

	routerBuilder    *routing.RouterBuilder
	computingBuilder *computing.DefaultComputingBuilder

	simplugins      []types.SimulationPlugin
	statePluginRepo *types.StatePluginRepository
	maxCores        int
	running         bool

	simulationStateSerializer *SimulationStateSerializer
}

// NewSimulationService initializes the simulation service
func NewSimulationService(
	config *configs.SimulationConfig,
	router *routing.RouterBuilder,
	computing *computing.DefaultComputingBuilder,
	plugins []types.SimulationPlugin,
	statePluginRepo *types.StatePluginRepository,
	simualtionStateOutputFile *string,
) *SimulationService {
	simService := &SimulationService{
		routerBuilder:    router,
		computingBuilder: computing,
		maxCores:         config.MaxCpuCores,
		simplugins:       plugins,
		statePluginRepo:  statePluginRepo,
	}
	simService.BaseSimulationService = NewBaseSimulationService(config, simService.runSimulationStep)

	if *simualtionStateOutputFile != "" {
		simService.simulationStateSerializer = NewSimulationStateSerializer(*simualtionStateOutputFile)
		log.Printf("Simulation state will be serialized to %s", *simualtionStateOutputFile)
	}

	return simService
}

func (s *SimulationService) GetStatePluginRepository() *types.StatePluginRepository {
	return s.statePluginRepo
}

func (s *SimulationService) Close() {
	if s.simulationStateSerializer != nil {
		s.simulationStateSerializer.Save(s)
	}
}

// runSimulationStep is the core loop to simulate node and orchestrator logic
func (s *SimulationService) runSimulationStep(nextTime func(time.Time) time.Time) {
	if s.running {
		return
	}
	s.lock.Lock()
	if s.running {
		s.lock.Unlock()
		return
	}
	s.running = true
	s.lock.Unlock()

	s.setSimulationTime(nextTime(s.GetSimulationTime()))
	log.Printf("Simulation time is %s", s.simTime.Format(time.RFC3339))

	// Update positions of all nodes (satellites and ground stations)
	var wg sync.WaitGroup
	for _, n := range s.all {
		wg.Add(1)
		go func(n types.Node) {
			defer wg.Done()
			n.UpdatePosition(s.simTime) // Update each node's position
		}(n)
	}
	wg.Wait()

	// Link updates (ISL and ground links)
	for _, node := range s.all {
		wg.Add(1)
		go func(n types.Node) {
			defer wg.Done()
			node.GetLinkNodeProtocol().UpdateLinks()
		}(node)
	}
	wg.Wait()

	// Routing and computation (if enabled)
	if s.config.UsePreRouteCalc {
		for _, node := range s.all {
			wg.Add(1)
			go func(n types.Node) {
				defer wg.Done()
				n.GetRouter().CalculateRoutingTableAsync()
			}(node)
		}
		wg.Wait()
	}

	// Check if the orchestrator needs to reschedule
	if s.orchestrator != nil {
		log.Println("Checking orchestrator for reschedule...")
		// s.orchestrator.CheckReschedule()
	}

	// Execute post-step state plugins
	for _, plugin := range s.statePluginRepo.GetAllPlugins() {
		plugin.PostSimulationStep(s)
	}

	// Execute post-step simulation plugins
	for _, plugin := range s.simplugins {
		if err := plugin.PostSimulationStep(s); err != nil {
			log.Printf("Plugin %s PostSimulationStep error: %v", plugin.Name(), err)
		}
	}

	if s.simulationStateSerializer != nil {
		s.simulationStateSerializer.AddState(s)
	}

	time.Sleep(1 * time.Second) // Simulate step duration

	s.running = false
}
