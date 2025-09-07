package node

import (
	"math"
	"time"

	"github.com/keniack/stardustGo/internal/links/linktypes"
	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.Node = (*Satellite)(nil) // Ensure Satellite implements Node

// Satellite represents a single satellite node in the simulation.
type Satellite struct {
	BaseNode // Embedding Node struct to satisfy the Node interface

	Inclination          float64
	InclinationRad       float64
	RightAscension       float64
	RightAscensionRad    float64
	Eccentricity         float64
	ArgumentOfPerigee    float64
	ArgumentOfPerigeeRad float64
	MeanAnomaly          float64
	MeanMotion           float64
	SemiMajorAxis        float64
	Epoch                time.Time
	ISLProtocol          types.InterSatelliteLinkProtocol
	GroundLinks          []types.Link
	Position             types.Vector
}

// NewSatellite initializes a new Satellite object with orbital configuration and ISL protocol.
func NewSatellite(name string, inclination, raan, ecc, argPerigee, meanAnomaly, meanMotion float64, epoch time.Time, simTime time.Time, isl types.InterSatelliteLinkProtocol, router types.Router, computing types.Computing) *Satellite {
	inclRad := types.DegreesToRadians(inclination)
	raanRad := types.DegreesToRadians(raan)
	argPerigeeRad := types.DegreesToRadians(argPerigee)

	s := &Satellite{
		BaseNode:             BaseNode{Name: name, Router: router, Computing: computing}, // Embedding Node struct
		Inclination:          inclination,
		InclinationRad:       inclRad,
		RightAscension:       raan,
		RightAscensionRad:    raanRad,
		Eccentricity:         ecc,
		ArgumentOfPerigee:    argPerigee,
		ArgumentOfPerigeeRad: argPerigeeRad,
		MeanAnomaly:          meanAnomaly,
		MeanMotion:           meanMotion,
		Epoch:                epoch,
		ISLProtocol:          isl,
		GroundLinks:          []types.Link{},
	}

	isl.Mount(s)
	router.Mount(s)
	s.UpdatePosition(simTime)
	return s
}

// Implementing Node methods via the embedded Node struct

// GetName returns the name of the satellite (from Node)
func (s *Satellite) GetName() string {
	return s.Name
}

// PositionVector returns the satellite's current position
func (s *Satellite) PositionVector() types.Vector {
	return s.Position
}

// DistanceTo calculates the distance between this satellite and another node (satellite or ground station)
func (s *Satellite) DistanceTo(other types.Node) float64 {
	return s.Position.Sub(other.PositionVector()).Magnitude()
}

func (s *Satellite) GetComputing() types.Computing {
	return s.Computing
}

// GetLinks returns all links connected to the satellite (both ISL and ground links)
func (s *Satellite) GetLinks() []types.Link {
	// Combine inter-satellite links and ground links
	var allLinks []types.Link

	// Add inter-satellite links (ISL links)
	for _, link := range s.ISLProtocol.Links() {
		allLinks = append(allLinks, link)
	}

	// Add ground links
	for _, groundLink := range s.GroundLinks {
		allLinks = append(allLinks, groundLink)
	}

	return allLinks
}

func (s *Satellite) GetEstablishedLinks() []types.Link {
	var establishedLinks []types.Link
	s.ISLProtocol.Established()
	for _, link := range s.ISLProtocol.Established() {
		establishedLinks = append(establishedLinks, link)
	}
	for _, groundLink := range s.GroundLinks {
		establishedLinks = append(establishedLinks, groundLink)
	}
	return establishedLinks
}

// UpdatePosition calculates the satellite's position in the ECI frame based on orbital elements and simulation time
func (s *Satellite) UpdatePosition(simTime time.Time) {
	deltaT := simTime.Sub(s.Epoch).Seconds() // Time since epoch in seconds
	meanMotionRadPerSec := s.MeanMotion * 2.0 * math.Pi / (24 * 3600)
	meanAnomalyCurrent := s.MeanAnomaly + meanMotionRadPerSec*deltaT
	meanAnomalyCurrent = normalizeAngle(meanAnomalyCurrent)
	eccentricAnomaly := solveKeplersEquation(meanAnomalyCurrent, s.Eccentricity)
	trueAnomaly := computeTrueAnomaly(eccentricAnomaly, s.Eccentricity)

	semiMajorAxis := 6790000.0 // Approx. value for LEO satellites
	distance := semiMajorAxis * (1 - s.Eccentricity*math.Cos(eccentricAnomaly))
	xp := distance * math.Cos(trueAnomaly)
	yp := distance * math.Sin(trueAnomaly)
	zp := 0.0

	s.Position = applyOrbitalTransformations(xp, yp, zp, s.InclinationRad, s.ArgumentOfPerigeeRad, s.RightAscensionRad)
}

// ApplyOrbitalTransformations converts orbital plane coordinates into the Earth-Centered Inertial (ECI) frame
func applyOrbitalTransformations(x, y, z, iRad, omegaRad, raanRad float64) types.Vector {
	cosRAAN := math.Cos(raanRad)
	sinRAAN := math.Sin(raanRad)
	cosIncl := math.Cos(iRad)
	sinIncl := math.Sin(iRad)
	cosArgP := math.Cos(omegaRad)
	sinArgP := math.Sin(omegaRad)

	xECI := (cosRAAN*cosArgP-sinRAAN*sinArgP*cosIncl)*x + (-cosRAAN*sinArgP-sinRAAN*cosArgP*cosIncl)*y
	yECI := (sinRAAN*cosArgP+cosRAAN*sinArgP*cosIncl)*x + (-sinRAAN*sinArgP+cosRAAN*cosArgP*cosIncl)*y
	zECI := sinIncl*sinArgP*x + sinIncl*cosArgP*y

	return types.Vector{X: xECI, Y: yECI, Z: zECI}
}

// normalizeAngle wraps an angle in radians into the range [0, 2π].
func normalizeAngle(rad float64) float64 {
	for rad < 0 {
		rad += 2 * math.Pi
	}
	for rad > 2*math.Pi {
		rad -= 2 * math.Pi
	}
	return rad
}

// solveKeplersEquation uses Newton-Raphson iteration to solve for the eccentric anomaly.
func solveKeplersEquation(meanAnomaly, ecc float64) float64 {
	E := meanAnomaly
	delta := 1.0
	tol := 1e-6
	for math.Abs(delta) > tol {
		delta = (E - ecc*math.Sin(E) - meanAnomaly) / (1 - ecc*math.Cos(E))
		E -= delta
	}
	return E
}

// computeTrueAnomaly calculates the true anomaly from the eccentric anomaly.
func computeTrueAnomaly(E, ecc float64) float64 {
	sqrt1me2 := math.Sqrt(1 - ecc*ecc)
	return math.Atan2(sqrt1me2*math.Sin(E), math.Cos(E)-ecc)
}

// ConfigureConstellation configures a constellation of satellites by linking them.
func (s *Satellite) ConfigureConstellation(satellites []*Satellite) {
	for _, satellite := range satellites {
		// Skip if it's the same satellite (this) or if there's already a link
		if satellite == s { // Or add more conditions here if needed (e.g., checking existing links)
			continue
		}

		// Create a new ISL link between the current satellite and the other one
		link := linktypes.NewIslLink(s, satellite)

		// Locking to ensure thread safety while modifying ISLProtocol
		s.ISLProtocol.AddLink(link)         // Add link to this satellite's ISL protocol
		satellite.ISLProtocol.AddLink(link) // Add link to the other satellite's ISL protocol
	}
}
