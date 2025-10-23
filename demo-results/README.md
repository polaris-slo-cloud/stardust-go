# Simulation Results

This directory contains the results of various simulations run using the Stardust simulator. The results are organized to provide insights into the performance and behavior of the simulations.

## Overview of Results
Simulation results are stored in different formats to facilitate analysis and visualization. These include:

- **CSV Files**: Contain CPU and memory metrics recorded during the execution.
- **LOG Files**: Contain the logs of the simulation.
- **Binary Files (GOB)**: Precomputed data for replaying simulations.

## How to Use This Directory
1. **Analyze Metrics**:
   - Open CSV files to review CPU usage and memory consumption over time.
   - Use tools like Python or Excel to create visualizations.
   - Use log files to identify some key moments in the simulation.

2. **Replay Simulations**:
   - Use binary files to replay simulations and observe system behavior ([click here to see how to run in precomputed mode](../README.md#-running-the-stardustgo-simulator)).

For example, when we compare mst for newest constellation in [simulated](./mst/simulated/newest.csv) and [precomputed](./mst/precomputed/newest.csv) mode, we can see that the simulation was 2 minutes for simulation mode but only 10 seconds in precomputed mode. 
In simulated mode we can see the link selection process is expensive and single threaded for mst, since most of the time all cores except one are ideling.
The precomputation mode only reads the links from binary file, so no calculations needed, which speeds up the simulation process.
The log files are identical (except log-timestamps) which shows that the simulations are repeatable using the precomputed data.

Let us know if you have any questions or need further assistance interpreting the results!