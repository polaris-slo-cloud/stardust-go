package types

import (
	"errors"
	"math"
	"sync"
	"time"
)

var _ Node = (*GroundStation)(nil)

// GroundStation represents an Earth-based node that links to satellites
// It updates its position over time and tracks the nearest satellites
type GroundStation struct {
	BaseNode

	Latitude                    float64
	Longitude                   float64
	SimulationStartTime         time.Time
	GroundSatelliteLinkProtocol GroundSatelliteLinkProtocol

	Position Vector
	mu       sync.Mutex
}

// NewGroundStation creates and initializes a new ground station with link protocol and position
func NewGroundStation(name string, lat float64, lon float64, protocol GroundSatelliteLinkProtocol, simStart time.Time, router Router, computing Computing) *GroundStation {
	gs := &GroundStation{
		BaseNode: BaseNode{
			Name:      name,
			Router:    router,
			Computing: computing,
		},
		Latitude:                    lat,
		Longitude:                   lon,
		SimulationStartTime:         simStart,
		GroundSatelliteLinkProtocol: protocol,
	}
	protocol.Mount(gs)
	router.Mount(gs)
	gs.UpdatePositionFromElapsed(0)
	return gs
}

func (gs *GroundStation) GetName() string {
	return gs.Name
}

func (gs *GroundStation) PositionVector() Vector {
	return gs.Position
}

func (gs *GroundStation) DistanceTo(other Node) float64 {
	return gs.Position.Subtract(other.PositionVector()).Magnitude()
}

// UpdatePosition sets the current position of the ground station based on simulation time
func (gs *GroundStation) UpdatePosition(simTime time.Time) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	timeElapsed := simTime.Sub(gs.SimulationStartTime).Seconds()
	gs.UpdatePositionFromElapsed(timeElapsed)
}

// UpdatePositionFromElapsed calculates Earth-centered coordinates using geodetic formula
func (gs *GroundStation) UpdatePositionFromElapsed(timeElapsed float64) {
	const (
		a             = 6378137.0       // semi-major axis in meters
		b             = 6356752.314245  // semi-minor axis in meters
		e2            = 1 - (b*b)/(a*a) // eccentricity squared
		rotationSpeed = 7.2921150e-5    // Earth's rotation speed rad/s
	)

	latRad := DegreesToRadians(gs.Latitude)
	lonRad := DegreesToRadians(gs.Longitude)
	alt := 0.0

	N := a / math.Sqrt(1-e2*math.Sin(latRad)*math.Sin(latRad))

	x := (N + alt) * math.Cos(latRad) * math.Cos(lonRad)
	y := (N + alt) * math.Cos(latRad) * math.Sin(lonRad)
	z := ((b * b / (a * a) * N) + alt) * math.Sin(latRad)

	theta := rotationSpeed * timeElapsed
	xRot := x*math.Cos(theta) - y*math.Sin(theta)
	yRot := x*math.Sin(theta) + y*math.Cos(theta)
	zRot := z

	gs.Position = Vector{X: xRot, Y: yRot, Z: zRot}
}

func (gs *GroundStation) GetLinkNodeProtocol() LinkNodeProtocol {
	return gs.GroundSatelliteLinkProtocol
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

func (gs *GroundStation) GetLinks() []Link {
	return gs.GroundSatelliteLinkProtocol.Links()
}

func (gs *GroundStation) GetEstablishedLinks() []Link {
	return gs.GroundSatelliteLinkProtocol.Established()
}

func (gs *GroundStation) GetRouter() Router {
	return gs.Router
}
