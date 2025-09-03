package simulation

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/computing"
	"github.com/keniack/stardustGo/internal/deployment"
	"github.com/keniack/stardustGo/internal/node"
	"github.com/keniack/stardustGo/internal/routing"
	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.SimulationController = (*SimulationService)(nil)

// SimulationService handles simulation lifecycle and state updates
type SimulationService struct {
	config           configs.SimulationConfig
	routerBuilder    *routing.RouterBuilder
	computingBuilder *computing.DefaultComputingBuilder

	all         []types.Node
	satellites  []*node.Satellite
	groundNodes []*node.GroundStation
	simTime     time.Time
	maxCores    int
	lock        sync.Mutex

	autorun      bool
	running      bool
	orchestrator *deployment.DeploymentOrchestrator
}

// NewSimulationService initializes the simulation service
func NewSimulationService(
	config configs.SimulationConfig,
	router *routing.RouterBuilder,
	computing *computing.DefaultComputingBuilder,
) *SimulationService {
	return &SimulationService{
		config:           config,
		routerBuilder:    router,
		computingBuilder: computing,
		all:              []types.Node{},
		satellites:       []*node.Satellite{},
		groundNodes:      []*node.GroundStation{},
		simTime:          config.SimulationStartTime,
		maxCores:         config.MaxCpuCores,
	}
}

// Inject sets the orchestrator dependency
func (s *SimulationService) Inject(o *deployment.DeploymentOrchestrator) {
	s.orchestrator = o
}

// InjectSatellites adds the loaded satellites to the simulation scope
func (s *SimulationService) InjectSatellites(satellites []types.Node) error {
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

// InjectGroundStations adds the loaded ground stations to the simulation scope
func (s *SimulationService) InjectGroundStations(groundStations []types.Node) error {
	s.groundNodes = make([]*node.GroundStation, 0, len(groundStations))
	for _, n := range groundStations {
		gs, ok := n.(*node.GroundStation)
		if !ok {
			return fmt.Errorf("InjectGroundStations: expected *node.GroundStation but got %T", n)
		}
		s.groundNodes = append(s.groundNodes, gs)
		s.all = append(s.all, gs) // Add ground station as generic nodes
	}

	log.Printf("Injected %d ground stations into simulation", len(s.satellites))
	return nil
}

// StartAutorun begins the simulation loop in autorun mode
func (s *SimulationService) StartAutorun() <-chan struct{} {
	s.lock.Lock()
	if s.autorun {
		s.lock.Unlock()
		done := make(chan struct{})
		close(done)
		return done // autorun already active
	}
	s.autorun = true
	s.lock.Unlock()

	done := make(chan struct{})
	go func() {
		// While autorun is enabled, run simulation steps at configured intervals
		for {
			if !s.autorun {
				break
			}

			s.runSimulationStep(func(prev time.Time) time.Time {
				return prev.Add(time.Duration(s.config.StepMultiplier) * time.Second)
			})

			time.Sleep(time.Duration(s.config.StepInterval) * time.Millisecond)
		}
		close(done) // closed when simulation loop exits
	}()

	return done
}

// StopAutorun disables autorun mode
func (s *SimulationService) StopAutorun() {
	s.autorun = false
}

// StepBySeconds executes a single step manually (e.g. UI trigger)
func (s *SimulationService) StepBySeconds(seconds float64) {
	s.runSimulationStep(func(prev time.Time) time.Time {
		return prev.Add(time.Duration(seconds * float64(time.Second)))
	})
}

// StepByTime executes a single step manually (e.g. UI trigger)
func (s *SimulationService) StepByTime(newTime time.Time) {
	s.runSimulationStep(func(prev time.Time) time.Time {
		return newTime
	})
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

	s.simTime = nextTime(s.simTime)
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

	time.Sleep(1 * time.Second) // Simulate step duration

	s.running = false
}

func (s *SimulationService) GetAllNodes() []types.Node {
	return s.all
}

func (s *SimulationService) GetSatellites() []types.Node {
	nodes := make([]types.Node, len(s.satellites))
	for i, sat := range s.satellites {
		nodes[i] = sat
	}
	return nodes
}

func (s *SimulationService) GetGroundStations() []types.Node {
	nodes := make([]types.Node, len(s.groundNodes))
	for i, gs := range s.groundNodes {
		nodes[i] = gs
	}
	return nodes
}
