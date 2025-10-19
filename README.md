# StardustGo
*A scalable and extensible simulator for the 3D Continuum*

## Overview
StardustGo is a modular, extensible simulation framework for modeling and analyzing space-ground computing constellations. It includes abstractions for routing, inter-satellite link (ISL) protocols, satellite dynamics, and orchestrated deployments.

Stardust is an open-source, scalable simulator designed to:  

- Run resouce efficient on a single machine
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

You can plug in your own service logic by using the SimulationController in main entrypoint or by implementing [SimPlugin](./go/internal/simplugin/dummy_plugin.go) or [StatePlugin](./go/internal/stateplugin/dummy_sun_state_plugin.go).

## 🏗️ Architecture & Extensibility

StardustGo provides a flexible plugin architecture that allows developers to extend the simulator with custom components. The system is built around well-defined interfaces that enable seamless integration of new functionality.

### Node Type Architecture

The simulator supports extensible node types through interfaces defined in `./go/pkg/types`. You can implement custom component behaviors for different environments by matching interface implementations:

- Node
- Link
- LinkNodeProtocol
- Router
- SimulationPlugin
- StatePlugin

For example, to add a new node type:
1. Implement the `Node` interface from `./go/pkg/types/`
2. Define the node's computational and networking capabilities
3. Register the node type with the simulation framework

### Plugin Architecture

StardustGo supports two primary plugin types:

#### Simulation Plugins
Located in `./go/internal/simplugins/`, these plugins extend simulation behavior for scenario-specific simulation logic.

#### State Plugins  
Located in `./go/internal/stateplugins/`, these plugins manage simulation state:
- Custom state persistence mechanisms
- Scenario-specific simulation state logic (i.e. energy consumption and usage)

## 🧱 Project Structure
```aiignore
├── cmd/stardust/           # Main entry point
├── configs/                # Configuration files
├── internal/
│   ├── computing/          # Compute strategies
│   ├── deployment/         # Orchestration strategies
│   ├── ground/             # Utils to load ground stations
│   ├── links/              # Links and link protocols
│   ├── node/               # Node and ground station modeling
│   ├── routing/            # Routing protocols
│   ├── satellite/          # Utils to load satellite constellations
│   ├── simulation/         # Simulation engine
│   ├── simplugins/         # Simulation plugins
│   └── stateplugins/       # State plugins
├── pkg/types/              # Interfaces and shared types
├── resources/
│   ├── configs/            # configurations
│   └── tle/                # TLE datasets
└── go.mod                  # Module definition
```
