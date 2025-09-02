package deployment

import (
	"errors"
	"sync"

	"github.com/keniack/stardustGo/pkg/types"
)

// DeploymentOrchestrator orchestrates deployment actions based on specifications.
type DeploymentOrchestrator struct {
	resolver       *DeploymentOrchestratorResolver
	specifications []types.DeploymentSpecification
	mu             sync.Mutex // Protects access to specifications
}

// NewDeploymentOrchestrator creates a new DeploymentOrchestrator.
func NewDeploymentOrchestrator() *DeploymentOrchestrator {
	// Initialize the resolver but do not rely on simulation controller at this stage
	resolver, _ := NewDeploymentOrchestratorResolver([]types.DeploymentOrchestrator{})

	// Just create the DeploymentOrchestrator without trying to inject satellites
	return &DeploymentOrchestrator{
		resolver: resolver,
	}
}

// DeploymentTypes returns the supported deployment types.
func (d *DeploymentOrchestrator) DeploymentTypes() ([]string, error) {
	// This method needs to be implemented
	return nil, errors.New("method not implemented")
}

// CheckRescheduleAsync checks for any rescheduling needs for each specification.
func (d *DeploymentOrchestrator) CheckRescheduleAsync(deployment types.DeploymentSpecification) error {
	// Iterate over all specifications
	d.mu.Lock()
	defer d.mu.Unlock()

	var wg sync.WaitGroup
	for _, spec := range d.specifications {
		wg.Add(1)
		go func(spec types.DeploymentSpecification) {
			defer wg.Done()
			orchestrator, _ := d.resolver.Resolve(spec)
			_ = orchestrator.CheckRescheduleAsync(spec)
		}(spec)
	}
	wg.Wait()
	return nil
}

// CreateDeploymentAsync adds a new deployment specification and creates the deployment.
func (d *DeploymentOrchestrator) CreateDeploymentAsync(deployment types.DeploymentSpecification) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Add deployment to specifications
	d.specifications = append(d.specifications, deployment)

	// Resolve the orchestrator and create the deployment
	orchestrator, _ := d.resolver.Resolve(deployment)
	return orchestrator.CreateDeploymentAsync(deployment)
}

// DeleteDeploymentAsync removes a deployment specification and deletes it.
func (d *DeploymentOrchestrator) DeleteDeploymentAsync(deployment types.DeploymentSpecification) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Remove the deployment specification
	var newSpecifications []types.DeploymentSpecification
	for _, spec := range d.specifications {
		if spec != deployment {
			newSpecifications = append(newSpecifications, spec)
		}
	}
	d.specifications = newSpecifications

	// Resolve the orchestrator and delete the deployment
	orchestrator, _ := d.resolver.Resolve(deployment)
	return orchestrator.DeleteDeploymentAsync(deployment)
}
