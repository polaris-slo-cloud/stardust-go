package simulation

import (
	"log"
	"sync"
	"time"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.SimulationController = (*SimulationIteratorService)(nil)

type SimulationIteratorService struct {
	BaseSimulationService

	simulationStates      []types.SimulationState
	simPlugins            []types.SimulationPlugin
	statePluginRepository types.StatePluginRepository
	running               bool
	currentIx             int
}

func NewSimulationIteratorService(config *configs.SimulationConfig, simulationStates []types.SimulationState, simPlugins []types.SimulationPlugin, statePluginRepository types.StatePluginRepository) *SimulationIteratorService {
	service := &SimulationIteratorService{
		simulationStates:      simulationStates,
		simPlugins:            simPlugins,
		statePluginRepository: statePluginRepository,
		running:               false,
		currentIx:             -1,
	}
	service.BaseSimulationService = NewBaseSimulationService(config, service.runSimulationStep)
	return service
}

func (s *SimulationIteratorService) GetStatePluginRepository() *types.StatePluginRepository {
	return &s.statePluginRepository
}

func (s *SimulationIteratorService) Close() {

}

func (s *SimulationIteratorService) runSimulationStep(nextTime func(time.Time) time.Time) {
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

	s.currentIx++
	s.setSimulationTime(s.simulationStates[s.currentIx].Time)
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
	for _, plugin := range s.statePluginRepository.GetAllPlugins() {
		plugin.PostSimulationStep(s)
	}

	// Execute post-step simulation plugins
	for _, plugin := range s.simPlugins {
		if err := plugin.PostSimulationStep(s); err != nil {
			log.Printf("Plugin %s PostSimulationStep error: %v", plugin.Name(), err)
		}
	}

	time.Sleep(1 * time.Second) // Simulate step duration

	s.running = false
}
