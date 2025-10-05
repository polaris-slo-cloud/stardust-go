# StardustGo
*A scalable and extensible simulator for the 3D Continuum*
## Overview
StardustGo is a modular, extensible simulation framework for modeling and analyzing space-ground computing constellations. It includes abstractions for routing, inter-satellite link (ISL) protocols, satellite dynamics, and orchestrated deployments.

Stardust is an open-source, scalable simulator designed to:  

- Simulate mega-constellations of up to 20.6k satellites on a single machine
- Support dynamic routing protocols for experimentation
- Provide SimPlugin and StatePlugin extensibility to integrate and test your code directly  
- Provide a precomputed simulation mode, where the simulation state for every step is precomputed before the simulation starts
- Cover the entire 3D Continuum: Edge, Cloud, and Space

---

## Prerequisites

Before running the StardustGo Simulator, ensure you have Go installed on your system ([official Go installation guide](https://go.dev/doc/install)).

Then you can clone the repository and build the project:

```bash
# Clone repo
git clone https://github.com/polaris-slo-cloud/stardust-go.git
cd stardust-go/go

# Build
go build -o . ./...
```

---

## 🚀 Running the StardustGo Simulator

From the project root, you can run the simulator in simulation mode with the following command:

```bash
go run ./cmd/stardust \
  --simulationConfig <path-to-simulation-config> \
  --islConfig <path-to-isl-config> \
  --groundLinkConfig <path-to-ground-link-config> \
  --computingConfig <path-to-computing-config> \
  --routerConfig <path-to-router-config> \
  [--simulationStateOutputFile <output-file-path>] \
  [--simulationPlugins <comma-separated-plugin-names>] \
  [--statePlugins <comma-separated-plugin-names>]
```

From the project root, you can run the simulator in precomputed mode with the following command:

```bash
go run ./cmd/stardust \
  --simulationConfig <path-to-simulation-config> \
  --computingConfig <path-to-computing-config> \
  --routerConfig <path-to-router-config> \
  [--simulationStateInputFile <output-file-path>] \
  [--simulationPlugins <comma-separated-plugin-names>]
```

### Run a Sample Simulation

```bash
# Run simulation with 500 satellites and 85 ground stations and save precomputed data in file
./stardust \
  --simulationConfig ./resources/configs/simulationAutorunConfig.yaml \
  --islConfig ./resources/configs/islMstConfig.yaml \
  --groundLinkConfig ./resources/configs/groundLinkNearestConfig.yaml \
  --computingConfig ./resources/configs/computingConfig.yaml \
  --routerConfig ./resources/configs/routerAStarConfig.yaml \
  --simulationStateOutputFile precomputed_data.gob
```


```bash
# Run simulation with 500 satellites and 85 ground stations in precomputed mode using previous saved simulation state data
./stardust \
  --simulationConfig ./resources/configs/simulationAutorunConfig.yaml \
  --computingConfig ./resources/configs/computingConfig.yaml \
  --routerConfig ./resources/configs/routerAStarConfig.yaml \
  --simulationStateInputFile precomputed_data.gob
```

---

## ⚙️ Configuration

Edit the simulation configuration files in the `./resources/configs/` directory. For detailed documentation on available config types, fields, and examples, see the [Configuration Guide](./go/resources/configs/README.md).

---

## 🧠 Writing Your Own Simulation Logic

You can plug in your own service logic by using the SimulationController in main entrypoint or by implementing SimPlugin or StatePlugin.

```go
// TODO main, SimPlugin, StatePlugin
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
