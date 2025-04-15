# StardustGo
A simulator for the 3D Continuum
## Overview
StardustGo is a modular, extensible simulation framework for modeling and analyzing space-ground computing constellations. It includes abstractions for routing, inter-satellite link (ISL) protocols, satellite dynamics, and orchestrated deployments.

---

## 🚀 Running the StardustGo Simulator

From the project root, run the simulator:

```bash
go run ./cmd/stardust
```

## ⚙️ Configuration

Edit configs/appsettings.json to control the simulation. Important configuration options include:

    SimulationConfiguration.SatelliteDataSource: File path to a TLE file (e.g. resources/tle/starlink_500.tle)

    SimulationConfiguration.SimulationStartTime: Timestamp for simulation start (ISO8601)

    InterSatelliteLinkConfig.Protocol: Choose from:

        mst, mst_loop, mst_smart_loop

        pst, pst_loop, pst_smart_loop

    RouterConfig.Protocol: Choose from:

        a-star, dijkstra

## 🧠 Writing Your Own Simulation Logic

You can plug in your own service logic by using the ISimulationController interface. Here's an example of a custom service:

```go
type YourService struct {
	simulation types.ISimulationController
}

func (s *YourService) Start() error {
	nodes := s.simulation.GetAllNodes()
	satellites := s.simulation.GetSatellites()

	for step := 0; step < 100; step++ {
		// Simulate a 60-second step
		s.simulation.Step(60 * time.Second)

		// Insert your own logic here
	}

	return nil
}
```
## 🧱 Project Structure
```aiignore
├── cmd/stardust/           # Main entry point
├── configs/                # Configuration files
├── internal/
│   ├── satellite/          # Node and constellation modeling
│   ├── routing/            # Routing protocols
│   ├── computing/          # Compute strategies
│   ├── deployment/         # Orchestration strategies
│   └── simulation/         # Simulation engine
├── pkg/types/              # Interfaces and shared types
├── resources/              # TLE files and data
└── go.mod                  # Module definition

```

## 📦 Build & Run
To build and run the simulator, use the following commands:

```bash
# Build the simulator
go build -o bin/stardust ./cmd/stardust
./bin/stardust
```

