#!/bin/bash

# Funktion zum Ausführen des Kommandos und Monitoring der CPU
execute() {
  local COMMAND=$1
  local OUTPUT_CSV=$2

  # Starte den Befehl im Hintergrund und speichere die PID
  $COMMAND &
  local PID=$!

  # CSV-Header schreiben
  echo "Timestamp,PID,CPU_Total(%),$(seq 0 $(( $(nproc) - 1 )) | sed 's/^/CPU_Core_/')" > "$OUTPUT_CSV"

  # Funktion, um die CPU-Auslastung zu sammeln und in die CSV zu schreiben
  monitor_cpu() {
    while kill -0 "$PID" 2>/dev/null; do
      # Hole die CPU-Auslastung pro Core für die PID
      CPU_STATS=$(top -b -n 1 -p "$PID" | awk 'NR==8')
      CPU_TOTAL=$(echo "$CPU_STATS" | awk '{print $9}')

      # Hole die CPU-Auslastung pro Core (mit mpstat)
      CORE_STATS=$(mpstat -P ALL 1 1 | awk '/^[0-9]/ {print $3}')

      # Aktuellen Timestamp erstellen
      TIMESTAMP=$(date +"%Y-%m-%d %H:%M:%S")

      # Schreibe die Daten in die CSV-Datei
      echo -n "$TIMESTAMP,$PID,$CPU_TOTAL," >> "$OUTPUT_CSV"
      echo "$CORE_STATS" | paste -sd, >> "$OUTPUT_CSV"

      # Warte 1 Sekunde, bevor die nächsten Daten gesammelt werden
      sleep 1
    done
  }

  # Starte das Monitoring
  monitor_cpu

  # Ausgabe, wenn der Prozess beendet ist
  echo "Prozess mit PID $PID wurde beendet. CPU-Monitoring abgeschlossen für $OUTPUT_CSV."
}

# Simulated MST
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-0250.yaml --islConfig ./resources/configs/islMstConfig.yaml --groundLinkConfig ./resources/configs/groundLinkNearestConfig.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateOutputFile ./precomputed/mst/precomputed-0250.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin > ./results/mst/simulated/0250.log" "./results/mst/simulated/0250.csv"
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-0500.yaml --islConfig ./resources/configs/islMstConfig.yaml --groundLinkConfig ./resources/configs/groundLinkNearestConfig.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateOutputFile ./precomputed/mst/precomputed-0500.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin > ./results/mst/simulated/0500.log" "./results/mst/simulated/0500.csv"
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-1000.yaml --islConfig ./resources/configs/islMstConfig.yaml --groundLinkConfig ./resources/configs/groundLinkNearestConfig.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateOutputFile ./precomputed/mst/precomputed-1000.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin > ./results/mst/simulated/1000.log" "./results/mst/simulated/1000.csv"
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-2000.yaml --islConfig ./resources/configs/islMstConfig.yaml --groundLinkConfig ./resources/configs/groundLinkNearestConfig.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateOutputFile ./precomputed/mst/precomputed-2000.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin > ./results/mst/simulated/2000.log" "./results/mst/simulated/2000.csv"
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-3000.yaml --islConfig ./resources/configs/islMstConfig.yaml --groundLinkConfig ./resources/configs/groundLinkNearestConfig.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateOutputFile ./precomputed/mst/precomputed-3000.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin > ./results/mst/simulated/3000.log" "./results/mst/simulated/3000.csv"
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-newest.yaml --islConfig ./resources/configs/islMstConfig.yaml --groundLinkConfig ./resources/configs/groundLinkNearestConfig.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateOutputFile ./precomputed/mst/precomputed-newest.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin > ./results/mst/simulated/newest.log" "./results/mst/simulated/newest.csv"
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-newest-double.yaml --islConfig ./resources/configs/islMstConfig.yaml --groundLinkConfig ./resources/configs/groundLinkNearestConfig.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateOutputFile ./precomputed/mst/precomputed-newest-double.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin > ./results/mst/simulated/newest-double.log" "./results/mst/simulated/newest-double.csv"
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-newest-triple.yaml --islConfig ./resources/configs/islMstConfig.yaml --groundLinkConfig ./resources/configs/groundLinkNearestConfig.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateOutputFile ./precomputed/mst/precomputed-newest-triple.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin > ./results/mst/simulated/newest-triple.log" "./results/mst/simulated/newest-triple.csv"

# Precomputed MST
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-0250.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateInputFile ./precomputed/mst/precomputed-0250.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin" "./results/mst/precomputed/0250.csv"
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-0500.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateInputFile ./precomputed/mst/precomputed-0500.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin" "./results/mst/precomputed/0500.csv"
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-1000.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateInputFile ./precomputed/mst/precomputed-1000.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin" "./results/mst/precomputed/1000.csv"
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-2000.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateInputFile ./precomputed/mst/precomputed-2000.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin" "./results/mst/precomputed/2000.csv"
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-3000.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateInputFile ./precomputed/mst/precomputed-3000.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin" "./results/mst/precomputed/3000.csv"
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-newest.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateInputFile ./precomputed/mst/precomputed-newest.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin" "./results/mst/precomputed/newest.csv"
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-newest-double.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateInputFile ./precomputed/mst/precomputed-newest-double.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin" "./results/mst/precomputed/newest-double.csv"
execute "./stardust --simulationConfig ./resources/configs/simulationManualConfig-newest-triple.yaml --computingConfig ./resources/configs/computingConfig.yaml --routerConfig ./resources/configs/routerAStarConfig.yaml --simulationStateInputFile ./precomputed/mst/precomputed-newest-triple.gob --simulationPlugins DummyPlugin --statePlugins DummyPlugin" "./results/mst/precomputed/newest-triple.csv"
