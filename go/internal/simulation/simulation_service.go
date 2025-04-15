package simulation

import (
	"fmt"
	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/computing"
	"github.com/keniack/stardustGo/internal/deployment"
	"github.com/keniack/stardustGo/internal/node"
	"github.com/keniack/stardustGo/internal/routing"
	"github.com/keniack/stardustGo/pkg/types"
	"log"
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
func (s *SimulationService) InjectSatellites(satellites []types.INode) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.satellites = make([]*node.Satellite, 0, len(satellites))
	for _, n := range satellites {
		sat, ok := n.(*node.Satellite)
		if !ok {
			return fmt.Errorf("InjectSatellites: expected *node.Satellite but got %T", n)
		}
		s.satellites = append(s.satellites, sat)
		s.all = append(s.all, sat) // Add satellites as generic nodes
	}

	log.Printf("Injected %d satellites into simulation", len(s.satellites))
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

func (s *SimulationService) StartAutorunAsync() <-chan struct{} {
	done := make(chan struct{})

	go func() {
		s.runSimulationStep(func(prev time.Time) time.Time {
			return prev.Add(time.Duration(s.config.StepMultiplier) * time.Second)
		})
		close(done) // closed when simulation loop exits
	}()

	return done
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
