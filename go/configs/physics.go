package configs

import "math"

// Physics constants used across the simulation.
const (
	// MU Earth's gravitational parameter in m^3/s^2
	MU = 398_600_441_800_000

	// EarthRadius Earth's radius in meters
	EarthRadius = 6_378_000

	// EarthRotationSpeed Earth's rotation speed in radians per second (2π / 86400)
	EarthRotationSpeed = 2 * math.Pi / 86400

	// MaxISLDistance Maximal distance for two satellites to communicate in meters
	MaxISLDistance = EarthRadius / 3

	// SpeedOfLight Speed of light in m/s
	SpeedOfLight = 299_792_000
)
