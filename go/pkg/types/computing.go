package types

// Computing is an interface for computing resources, including managing services, available resources, and service placement.
type Computing interface {
	// Mount attaches the computing unit to a node
	Mount(node *Node) error

	// TryPlaceDeploymentAsync tries to place a service on this computing unit
	TryPlaceDeploymentAsync(service IDeployableService) (bool, error)

	// RemoveDeploymentAsync removes a deployed service from the computing unit
	RemoveDeploymentAsync(service IDeployableService) error

	// CanPlace checks if the service can be placed on this computing unit
	CanPlace(service IDeployableService) bool

	// HostsService checks if the computing unit hosts a service by name
	HostsService(serviceName string) bool

	// CpuAvailable returns the remaining CPU available
	CpuAvailable() float64

	// MemoryAvailable returns the remaining memory available
	MemoryAvailable() float64

	// Clone creates a new copy of the current computing unit
	Clone() Computing

	// GetServices returns the list of deployed services
	GetServices() []IDeployableService
}
