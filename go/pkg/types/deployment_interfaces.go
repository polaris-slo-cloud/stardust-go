package types

// IDeployedService defines the structure for a deployed service.
type IDeployableService interface {
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

// IDeploymentSpecification defines the structure for a deployment specification.
type IDeploymentSpecification interface {
	// Type returns the type of the deployment.
	Type() string

	// Service returns the DeployableService associated with the deployment.
	Service() IDeployableService
}

// IDeploymentOrchestrator defines the structure for a deployment orchestrator.
type IDeploymentOrchestrator interface {
	// DeploymentTypes returns the list of supported deployment types.
	DeploymentTypes() []string

	// CreateDeploymentAsync initiates the creation of a deployment.
	CreateDeploymentAsync(deployment IDeploymentSpecification) error

	// DeleteDeploymentAsync initiates the deletion of a deployment.
	DeleteDeploymentAsync(deployment IDeploymentSpecification) error

	// CheckRescheduleAsync checks if a deployment needs to be rescheduled.
	CheckRescheduleAsync(deployment IDeploymentSpecification) error
}
