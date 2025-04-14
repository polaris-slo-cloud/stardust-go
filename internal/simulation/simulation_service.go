package simulation

import (
	"log"
	"stardustGo/configs"
	"stardustGo/internal/computing"
	"stardustGo/internal/deployment"
	"stardustGo/internal/node"
	"stardustGo/internal/routing"
	"stardustGo/pkg/types"
	"sync"
	"time"
)

// SimulationService handles simulation lifecycle and state updates
type SimulationService struct {
	config           configs.SimulationConfig
	routerBuilder    *routing.RouterBuilder
	computingBuilder *computing.ComputingBuilder

	all         []types.INode
	satellites  []*node.Satellite
	groundNodes []*node.GroundStation
	simTime     time.Time
	autorun     bool
	maxCores    int
	lock        sync.Mutex

	running      bool
	orchestrator *deployment.DeploymentOrchestrator
}

// NewSimulationService initializes the simulation service
func NewSimulationService(
	config configs.SimulationConfig,
	router *routing.RouterBuilder,
	computing *computing.ComputingBuilder,
) *SimulationService {
	return &SimulationService{
		config:           config,
		routerBuilder:    router,
		computingBuilder: computing,
		all:              []types.INode{},
		satellites:       []*node.Satellite{},
		groundNodes:      []*node.GroundStation{},
		simTime:          config.SimulationStartTime,
		autorun:          config.StepInterval >= 0,
		maxCores:         config.MaxCpuCores,
	}
}

// Inject sets the orchestrator dependency
func (s *SimulationService) Inject(o *deployment.DeploymentOrchestrator) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.orchestrator = o
}

// InjectSatellites adds the loaded satellites to the simulation scope
func (s *SimulationService) InjectSatellites(satellites []*node.Satellite) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.satellites = satellites
	for _, sat := range satellites {
		s.all = append(s.all, sat) // Add satellites as nodes
	}

	log.Printf("Injected %d satellites into simulation", len(satellites))
	return nil
}

// StartAsync initializes simulation components and starts the simulation
func (s *SimulationService) StartAsync() {
	log.Println("Starting simulation...")
	// Additional startup logic if needed
}

// StopAsync terminates the simulation loop
func (s *SimulationService) StopAsync() {
	log.Println("Stopping simulation...")
	s.lock.Lock()
	defer s.lock.Unlock()
	s.autorun = false
	s.running = false
}

// StartAutorunAsync launches a timed simulation loop
func (s *SimulationService) StartAutorunAsync() {
	s.lock.Lock()
	if s.autorun {
		s.lock.Unlock()
		return
	}
	s.autorun = true
	s.lock.Unlock()

	go s.runSimulationStep(func(prev time.Time) time.Time {
		return prev.Add(time.Duration(s.config.StepMultiplier) * time.Second)
	})
}

// StopAutorunAsync disables autorun mode
func (s *SimulationService) StopAutorunAsync() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.autorun = false
}

// StepAsync executes a single step manually (e.g. UI trigger)
func (s *SimulationService) StepAsync(seconds float64) {
	s.lock.Lock()
	if s.autorun {
		s.lock.Unlock()
		return
	}
	s.lock.Unlock()

	s.runSimulationStep(func(prev time.Time) time.Time {
		return prev.Add(time.Duration(seconds * float64(time.Second)))
	})
}

// runSimulationStep is the core loop to simulate node and orchestrator logic
func (s *SimulationService) runSimulationStep(nextTime func(time.Time) time.Time) {
	s.running = true

	for s.autorun || !s.running {
		s.simTime = nextTime(s.simTime)
		log.Printf("Simulation time is %s", s.simTime.Format(time.RFC3339))

		// Update positions of all nodes (satellites and ground stations)
		var wg sync.WaitGroup
		for _, n := range s.all {
			wg.Add(1)
			go func(n types.INode) {
				defer wg.Done()
				n.UpdatePosition(s.simTime) // Update each node's position
			}(n)
		}
		wg.Wait()

		// ISL updates (Inter-Satellite Links)
		for _, sat := range s.satellites {
			go sat.ISLProtocol.UpdateLinks()
		}

		// Routing and computation (if enabled)
		if s.config.UsePreRouteCalc {
			for _, sat := range s.satellites {
				go sat.Router.CalculateRoutingTableAsync()
			}
		}

		// Check if the orchestrator needs to reschedule
		if s.orchestrator != nil {
			log.Println("Checking orchestrator for reschedule...")
			// s.orchestrator.CheckReschedule()
		}

		if !s.autorun {
			break
		}

		time.Sleep(1 * time.Second) // Simulate step duration
	}

	s.running = false
}
