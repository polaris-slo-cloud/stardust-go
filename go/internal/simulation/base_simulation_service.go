package simulation

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/deployment"
	"github.com/keniack/stardustGo/pkg/types"
)

type BaseSimulationService struct {
	config *configs.SimulationConfig

	all         []types.Node
	satellites  []types.Satellite
	groundNodes []types.GroundStation

	autorun      bool
	simTime      time.Time
	orchestrator *deployment.DeploymentOrchestrator

	lock              sync.Mutex
	runSimulationStep func(func(time.Time) time.Time)
}

func NewBaseSimulationService(config *configs.SimulationConfig, runSimulationStep func(func(time.Time) time.Time)) BaseSimulationService {
	return BaseSimulationService{
		config:            config,
		all:               []types.Node{},
		satellites:        []types.Satellite{},
		groundNodes:       []types.GroundStation{},
		simTime:           config.SimulationStartTime,
		runSimulationStep: runSimulationStep,
	}
}

// Inject sets the orchestrator dependency
func (s *BaseSimulationService) Inject(o *deployment.DeploymentOrchestrator) {
	s.orchestrator = o
}

// InjectSatellites adds the loaded satellites to the simulation scope
func (s *BaseSimulationService) InjectSatellites(satellites []types.Node) error {
	s.satellites = make([]types.Satellite, len(satellites))
	for i, n := range satellites {
		sat, ok := n.(types.Satellite)
		if !ok {
			return fmt.Errorf("InjectSatellites: expected *node.Satellite but got %T", n)
		}
		s.satellites[i] = sat
		s.all = append(s.all, sat) // Add satellites as generic nodes
	}

	log.Printf("Injected %d satellites into simulation", len(s.satellites))
	return nil
}

// InjectGroundStations adds the loaded ground stations to the simulation scope
func (s *BaseSimulationService) InjectGroundStations(groundStations []types.Node) error {
	s.groundNodes = make([]types.GroundStation, len(groundStations))
	for i, n := range groundStations {
		gs, ok := n.(types.GroundStation)
		if !ok {
			return fmt.Errorf("InjectGroundStations: expected *node.GroundStation but got %T", n)
		}
		s.groundNodes[i] = gs
		s.all = append(s.all, gs) // Add ground station as generic nodes
	}

	log.Printf("Injected %d ground stations into simulation", len(s.groundNodes))
	return nil
}

// StartAutorun begins the simulation loop in autorun mode
func (s *BaseSimulationService) StartAutorun() <-chan struct{} {
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

func (s *BaseSimulationService) GetAllNodes() []types.Node {
	return s.all
}

func (s *BaseSimulationService) GetSatellites() []types.Satellite {
	// TODO find out why copying is needed here to avoid
	sats := make([]types.Satellite, len(s.satellites))
	copy(sats, s.satellites)
	return sats
}

func (s *BaseSimulationService) GetGroundStations() []types.GroundStation {
	return s.groundNodes
}

func (s *BaseSimulationService) GetSimulationTime() time.Time {
	return s.simTime
}

// StopAutorun disables autorun mode
func (s *BaseSimulationService) StopAutorun() {
	s.autorun = false
}

// StepBySeconds executes a single step manually (e.g. UI trigger)
func (s *BaseSimulationService) StepBySeconds(seconds float64) {
	s.runSimulationStep(func(prev time.Time) time.Time {
		return prev.Add(time.Duration(seconds * float64(time.Second)))
	})
}

// StepByTime executes a single step manually (e.g. UI trigger)
func (s *BaseSimulationService) StepByTime(newTime time.Time) {
	s.runSimulationStep(func(prev time.Time) time.Time {
		return newTime
	})
}
