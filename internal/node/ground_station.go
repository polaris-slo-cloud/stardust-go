package node

import (
	"errors"
	"math"
	"stardustGo/pkg/types"
	"sync"
	"time"
)

// GroundSatelliteLinkProtocol defines the interface for managing ground-to-satellite link behavior
type GroundSatelliteLinkProtocol interface {
	Mount(station *GroundStation)
	UpdateLink() error
	Link() *Link
}

// Link is a placeholder for actual link data structure
type Link struct {
	// Placeholder for the actual link properties
}

// GroundStation represents an Earth-based node that links to satellites
// It updates its position over time and tracks the nearest satellites

type GroundStation struct {
	Node

	Latitude                    float64
	Longitude                   float64
	SimulationStartTime         time.Time
	GroundSatelliteLinkProtocol GroundSatelliteLinkProtocol

	Position types.Vector
	mu       sync.Mutex
}

// NewGroundStation creates and initializes a new ground station with link protocol and position
func NewGroundStation(name string, lon, lat float64, link GroundSatelliteLinkProtocol, simStart time.Time, router types.IRouter, computing types.IComputing) *GroundStation {
	gs := &GroundStation{
		Node: Node{
			Name:      name,
			Router:    router,
			Computing: computing,
		},
		Longitude:                   lon,
		Latitude:                    lat,
		SimulationStartTime:         simStart,
		GroundSatelliteLinkProtocol: link,
	}
	gs.UpdatePositionFromElapsed(0)
	link.Mount(gs)
	return gs
}

func (gs *GroundStation) GetName() string {
	return gs.Name
}

func (gs *GroundStation) PositionVector() types.Vector {
	return gs.Position
}

func (gs *GroundStation) DistanceTo(other types.INode) float64 {
	return gs.Position.Sub(other.PositionVector()).Magnitude()
}

// UpdatePosition sets the current position of the ground station based on simulation time
func (gs *GroundStation) UpdatePosition(simTime time.Time) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	timeElapsed := simTime.Sub(gs.SimulationStartTime).Seconds()
	gs.UpdatePositionFromElapsed(timeElapsed)
	return gs.GroundSatelliteLinkProtocol.UpdateLink()
}

// UpdatePositionFromElapsed calculates Earth-centered coordinates using geodetic formula
func (gs *GroundStation) UpdatePositionFromElapsed(timeElapsed float64) {
	const (
		a             = 6378137.0       // semi-major axis in meters
		b             = 6356752.314245  // semi-minor axis in meters
		e2            = 1 - (b*b)/(a*a) // eccentricity squared
		rotationSpeed = 7.2921150e-5    // Earth's rotation speed rad/s
	)

	latRad := types.DegreesToRadians(gs.Latitude)
	lonRad := types.DegreesToRadians(gs.Longitude)
	alt := 0.0

	N := a / math.Sqrt(1-e2*math.Sin(latRad)*math.Sin(latRad))

	x := (N + alt) * math.Cos(latRad) * math.Cos(lonRad)
	y := (N + alt) * math.Cos(latRad) * math.Sin(lonRad)
	z := ((b * b / (a * a) * N) + alt) * math.Sin(latRad)

	theta := rotationSpeed * timeElapsed
	xRot := x*math.Cos(theta) - y*math.Sin(theta)
	yRot := x*math.Sin(theta) + y*math.Cos(theta)
	zRot := z

	gs.Position = types.Vector{X: xRot, Y: yRot, Z: zRot}
}

// FindNearestSatellite returns the closest satellite in a given list
func (gs *GroundStation) FindNearestSatellite(sats []*Satellite) (*Satellite, error) {
	if len(sats) == 0 {
		return nil, errors.New("satellite list is empty")
	}
	nearest := sats[0]
	minDist := gs.DistanceTo(nearest)
	for _, s := range sats[1:] {
		dist := gs.DistanceTo(s)
		if dist < minDist {
			nearest = s
			minDist = dist
		}
	}
	return nearest, nil
}
