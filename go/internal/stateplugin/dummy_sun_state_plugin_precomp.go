package stateplugin

import (
	"encoding/gob"
	"log"
	"os"
	"reflect"

	"github.com/keniack/stardustGo/pkg/helper"
	"github.com/keniack/stardustGo/pkg/types"
)

var _ SunStatePlugin = (*DummySunStatePluginPrecomp)(nil)

type DummySunStatePluginPrecomp struct {
	states    []map[string]float64
	currentIx int
}

func NewDummySunStatePrecompPlugin(origFile string) *DummySunStatePluginPrecomp {
	filename := helper.ExtendFilename(origFile, ".dummySimPlugin")
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	var states []map[string]float64
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&states); err != nil {
		log.Fatalf("failed to decode: %v", err)
	}

	return &DummySunStatePluginPrecomp{
		states:    states,
		currentIx: -1,
	}
}

func (d *DummySunStatePluginPrecomp) GetSunlightExposure(node types.Node) float64 {
	return d.states[d.currentIx][node.GetName()]
}

func (p *DummySunStatePluginPrecomp) GetName() string {
	return "DummyPlugin"
}

func (d *DummySunStatePluginPrecomp) GetType() reflect.Type {
	var dummy SunStatePlugin
	return reflect.TypeOf(dummy)
}

func (p *DummySunStatePluginPrecomp) PostSimulationStep(simulationController types.SimulationController) {
	p.currentIx++
}
func (p *DummySunStatePluginPrecomp) AddState(simulationController types.SimulationController) {
	// no-op for precomputed sim plugins
}

func (p *DummySunStatePluginPrecomp) Save(origFile string) {
	// no-op for precomputed sim plugins
}
