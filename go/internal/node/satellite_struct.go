package node

import (
	"math"
	"time"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.Satellite = (*SatelliteStruct)(nil) // Ensure SatelliteStruct implements Satellite

// SatelliteStruct represents a single satellite node in the simulation.
type SatelliteStruct struct {
	// Implementing Node methods via the embedded Node struct
	BaseNode // Embedding BaseNode struct to satisfy the Node interface

	inclination          float64
	inclinationRad       float64
	rightAscension       float64
	rightAscensionRad    float64
	eccentricity         float64
	argumentOfPerigee    float64
	argumentOfPerigeeRad float64
	meanAnomaly          float64
	meanMotion           float64
	epoch                time.Time
	ISLProtocol          types.InterSatelliteLinkProtocol
	groundLinks          []types.Link
}

// NewSatellite initializes a new Satellite object with orbital configuration and ISL protocol.
func NewSatellite(name string, inclination, raan, ecc, argPerigee, meanAnomaly, meanMotion float64, epoch time.Time, simTime time.Time, isl types.InterSatelliteLinkProtocol, router types.Router, computing types.Computing) *SatelliteStruct {
	inclRad := types.DegreesToRadians(inclination)
	raanRad := types.DegreesToRadians(raan)
	argPerigeeRad := types.DegreesToRadians(argPerigee)

	s := &SatelliteStruct{
		BaseNode:             BaseNode{Name: name, Router: router, Computing: computing}, // Embedding Node struct
		inclination:          inclination,
		inclinationRad:       inclRad,
		rightAscension:       raan,
		rightAscensionRad:    raanRad,
		eccentricity:         ecc,
		argumentOfPerigee:    argPerigee,
		argumentOfPerigeeRad: argPerigeeRad,
		meanAnomaly:          meanAnomaly,
		meanMotion:           meanMotion,
		epoch:                epoch,
		ISLProtocol:          isl,
		groundLinks:          []types.Link{},
	}

	isl.Mount(s)
	router.Mount(s)
	s.UpdatePosition(simTime)
	return s
}

// UpdatePosition calculates the satellite's position in the ECI frame based on orbital elements and simulation time
func (s *SatelliteStruct) UpdatePosition(simTime time.Time) {
	deltaT := simTime.Sub(s.epoch).Seconds() // Time since epoch in seconds
	meanMotionRadPerSec := s.meanMotion * 2.0 * math.Pi / (24 * 3600)
	meanAnomalyCurrent := s.meanAnomaly + meanMotionRadPerSec*deltaT
	meanAnomalyCurrent = normalizeAngle(meanAnomalyCurrent)
	eccentricAnomaly := solveKeplersEquation(meanAnomalyCurrent, s.eccentricity)
	trueAnomaly := computeTrueAnomaly(eccentricAnomaly, s.eccentricity)

	semiMajorAxis := 6790000.0 // Approx. value for LEO satellites
	distance := semiMajorAxis * (1 - s.eccentricity*math.Cos(eccentricAnomaly))
	xp := distance * math.Cos(trueAnomaly)
	yp := distance * math.Sin(trueAnomaly)
	zp := 0.0

	s.Position = applyOrbitalTransformations(xp, yp, zp, s.inclinationRad, s.argumentOfPerigeeRad, s.rightAscensionRad)
}

func (s *SatelliteStruct) GetLinkNodeProtocol() types.LinkNodeProtocol {
	return s.ISLProtocol
}

func (s *SatelliteStruct) GetISLProtocol() types.InterSatelliteLinkProtocol {
	return s.ISLProtocol
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
