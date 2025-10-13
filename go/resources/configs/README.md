# StardustGo Configuration Guide

This directory contains configuration files for the StardustGo Simulator. Each file defines specific aspects of the simulation, such as satellite behavior, routing, computing resources, and more.

## Simulation Config
Defines the core simulation parameters.

| Field                         | Type        | Description                                                                         |
|-------------------------------|-------------|-------------------------------------------------------------------------------------|
| `StepInterval`                | `int`       | Time interval (in seconds) between simulation steps. Use `-1` to indicate autorun.  |
| `StepMultiplier`              | `int`       | Multiplier for simulation speed (only used in autorun, e.g. set to `2` the simulation runs in double the speed) |
| `StepCount`                   | `int`       | Total number of steps to simulate (only used in autorun).                           |
| `SatelliteDataSource`         | `string`    | Path to the satellite data source file.                                             |
| `SatelliteDataSourceType`     | `string`    | Type of satellite data source (currently only `tle` supported).                     |
| `GroundStationDataSource`     | `string`    | Path to the ground station data source file.                                        |
| `GroundStationDataSourceType` | `string`    | Type of ground station data source (currently `yml` and `json` supported).          |
| `SimulationStartTime`         | `time.Time` | Start time of the simulation (ISO 8601 format).                                     |

**Example for autorun:** (`simulationAutorunConfig.yaml`)
```yaml
StepInterval: 1
StepMultiplier: 10
StepCount: 10
SatelliteDataSource: starlink_newest.tle
SatelliteDataSourceType: tle
GroundStationDataSource: ground_stations.yml
GroundStationDataSourceType: yml
SimulationStartTime: "2025-10-01T00:00:00Z"
```

**Example for manual:** (`simulationManualConfig.yaml`)
```yaml
StepInterval: -1
SatelliteDataSource: starlink_newest.tle
SatelliteDataSourceType: tle
GroundStationDataSource: ground_stations.yml
GroundStationDataSourceType: yml
SimulationStartTime: "2025-10-01T00:00:00Z"
```

## Inter-Satellite Link Config
Configures the inter-satellite communication link selection algorithm

| Field                     | Type      | Description                                                                             |
|---------------------------|-----------|-----------------------------------------------------------------------------------------|
| `Protocol`                | `string`  | Name of the link selection protocol (e.g., `mst`, `nearest`)                            |
| `Neighbours`              | `int`     | Numbers of links a satellite should establish (might gets ignored by some protocols).   |

**Example:** (`islMstConfig.yaml`)
```yaml
Neighbours: 4
Protocol: mst
```

## Ground Link Config
Configures communication links between ground stations and satellites

| Field                     | Type      | Description                                                                             |
|---------------------------|-----------|-----------------------------------------------------------------------------------------|
| `Protocol`                | `string`  | Name of the link selection protocol (currently only `nearest` supported)                |


**Example:** (`groundLinkNearestConfig.yaml`)
```yaml
Protocol: nearest
```

## Router Config
Defines the routing strategy for the simulation

| Field                     | Type      | Description                                                               |
|---------------------------|-----------|---------------------------------------------------------------------------|
| `Protocol`                | `string`  | Name of the routing protocol (e.g., `a-star`, `dijkstra`)                 |


**Example:** (`routerAStarConfig.yaml`)
```yaml
Protocol: a-star
```

## Computing  Config
Specifies computing resources for satellites or ground stations per computing type



| Field                     | Type      | Description                                                   |
|---------------------------|-----------|---------------------------------------------------------------|
| `Cores`                   | `int`     | Number of CPU cores.                                          |
| `Memory`                  | `int`     | Memory capacity (in MB).                                      |
| `Type`                    | `string`  | Type of computing resource (`None`, `Edge` or `Cloud`).               |

**Example:** (`computingConfig.yaml`)
```yaml
- Cores: 0
  Memory: 0
  Type: None
- Cores: 512
  Memory: 4096
  Type: Edge
- Cores: 1024
  Memory: 32768
  Type: Cloud
```

## File Formats

Configuration files can be in YAML (.yaml or .yml) or JSON (.json) format.
Use the appropriate file extension for your chosen format.