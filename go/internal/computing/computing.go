package computing

import (
	"fmt"
	"sync"
	"time"

	"github.com/keniack/stardustGo/pkg/types"
)

// Computing represents the computing resources of a node.
type Computing struct {
	Cpu         float64                   // Total CPU available
	Memory      float64                   // Total memory available
	Type        types.ComputingType       // Type of the computing unit
	CpuUsage    float64                   // Current CPU usage
	MemoryUsage float64                   // Current memory usage
	Services    []types.DeployableService // List of deployed services (using IDeployedService)
	mu          sync.Mutex                // Mutex to ensure thread safety
	node        types.Node                // Node to which this computing is mounted
}

func (c *Computing) GetServices() []types.DeployableService {
	return c.Services
}

// NewComputing creates a new Computing instance with the provided CPU, memory, and type.
func NewComputing(cpu, memory float64, ctype types.ComputingType) *Computing {
	return &Computing{
		Cpu:      cpu,
		Memory:   memory,
		Type:     ctype,
		Services: []types.DeployableService{},
	}
}

// Mount attaches this computing unit to a node
func (c *Computing) Mount(node types.Node) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.node != nil {
		return fmt.Errorf("computing is already mounted to node")
	}
	c.node = node
	return nil
}

// TryPlaceDeploymentAsync tries to place a service on this computing unit
func (c *Computing) TryPlaceDeploymentAsync(service types.DeployableService) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.node == nil {
		return false, fmt.Errorf("computing must be mounted to node before it can be used")
	}

	if !c.CanPlace(service) {
		return false, nil
	}

	c.Services = append(c.Services, service)
	c.CpuUsage += service.GetCpuUsage()
	c.MemoryUsage += service.GetMemoryUsage()

	// Simulate advertising the new service
	go func() {
		time.Sleep(1 * time.Second) // Simulate async operation
		// For example, advertising the service to other nodes or components
		// c.node.Router.AdvertiseNewServiceAsync(service.GetServiceName())
	}()

	return true, nil
}

// RemoveDeploymentAsync removes a deployed service from the computing unit
func (c *Computing) RemoveDeploymentAsync(service types.DeployableService) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Find and remove the service
	for i, s := range c.Services {
		if s.GetServiceName() == service.GetServiceName() {
			c.Services = append(c.Services[:i], c.Services[i+1:]...)
			c.CpuUsage -= service.GetCpuUsage()
			c.MemoryUsage -= service.GetMemoryUsage()
			return nil
		}
	}
	return fmt.Errorf("service %s not found", service.GetServiceName())
}

// CanPlace checks if the service can be placed on this computing unit
func (c *Computing) CanPlace(service types.DeployableService) bool {
	if service.GetCpuUsage() > c.CpuAvailable() {
		return false
	}
	if service.GetMemoryUsage() > c.MemoryAvailable() {
		return false
	}
	for _, s := range c.Services {
		if s.GetServiceName() == service.GetServiceName() {
			return false
		}
	}
	return true
}

// HostsService checks if the computing unit hosts a service by name
func (c *Computing) HostsService(serviceName string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, s := range c.Services {
		if s.GetServiceName() == serviceName {
			return true
		}
	}
	return false
}

// CpuAvailable returns the remaining CPU available
func (c *Computing) CpuAvailable() float64 {
	return c.Cpu - c.CpuUsage
}

// MemoryAvailable returns the remaining memory available
func (c *Computing) MemoryAvailable() float64 {
	return c.Memory - c.MemoryUsage
}

// Clone creates a new copy of the current computing unit and returns it as IComputing.
func (c *Computing) Clone() types.Computing {
	// Clone each deployed service
	servicesClone := make([]types.DeployableService, len(c.Services))
	copy(servicesClone, c.Services)

	return &Computing{
		Cpu:         c.Cpu,
		Memory:      c.Memory,
		Type:        c.Type,
		CpuUsage:    c.CpuUsage,
		MemoryUsage: c.MemoryUsage,
		Services:    servicesClone,
	}
}

func (c *Computing) GetComputingType() types.ComputingType {
	return c.Type
}
