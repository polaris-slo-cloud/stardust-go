package types

import (
	"math"
)

// Vector represents a 3D position in space
// Used for coordinates of nodes in the simulation
type Vector struct {
	X float64
	Y float64
	Z float64
}

// NewVector creates a new 3D vector
func NewVector(x, y, z float64) Vector {
	return Vector{X: x, Y: y, Z: z}
}

// Abs returns the magnitude of the vector (Euclidean norm)
func (v Vector) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Dot returns the dot product between two vectors
func (v Vector) Dot(other Vector) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross returns the cross product of two vectors
func (v Vector) Cross(other Vector) Vector {
	return Vector{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

// Normalize returns the unit vector in the same direction
func (v Vector) Normalize() Vector {
	length := v.Abs()
	if length == 0 {
		return Vector{}
	}
	return Vector{
		X: v.X / length,
		Y: v.Y / length,
		Z: v.Z / length,
	}
}

// Equals returns true if both vectors are exactly the same
func (v Vector) Equals(other Vector) bool {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z
}

// Subtract returns the vector difference (other - v)
func (v Vector) Subtract(other Vector) Vector {
	return Vector{
		X: other.X - v.X,
		Y: other.Y - v.Y,
		Z: other.Z - v.Z,
	}
}

// DegreesToRadians converts an angle in degrees to radians
func DegreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// Magnitude returns the Euclidean norm (length) of the vector.
func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}
