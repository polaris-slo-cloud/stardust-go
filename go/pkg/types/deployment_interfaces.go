package types

// IDeployedService defines the structure for a deployed service.
type DeployableService interface {
	// GetServiceName returns the name of the service.
	GetServiceName() string

	// GetCpuUsage returns the current CPU usage of the deployed service.
	GetCpuUsage() float64

	// GetMemoryUsage returns the current memory usage of the deployed service.
	GetMemoryUsage() float64

	// IsDeployed checks if the service has been successfully deployed.
	IsDeployed() bool

	// Deploy starts the service deployment process.
	Deploy() error

	// Remove stops the service and removes the deployment.
	Remove() error
}

// DeploymentSpecification defines the structure for a deployment specification.
type DeploymentSpecification interface {
	// Type returns the type of the deployment.
	Type() string

	// Service returns the DeployableService associated with the deployment.
	Service() DeployableService
}

// DeploymentOrchestrator defines the structure for a deployment orchestrator.
type DeploymentOrchestrator interface {
	// DeploymentTypes returns the list of supported deployment types.
	DeploymentTypes() []string

	// CreateDeploymentAsync initiates the creation of a deployment.
	CreateDeploymentAsync(deployment DeploymentSpecification) error

	// DeleteDeploymentAsync initiates the deletion of a deployment.
	DeleteDeploymentAsync(deployment DeploymentSpecification) error

	// CheckRescheduleAsync checks if a deployment needs to be rescheduled.
	CheckRescheduleAsync(deployment DeploymentSpecification) error
}
