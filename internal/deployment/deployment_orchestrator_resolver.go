package deployment

import (
	"errors"
	"fmt"
	"stardustGo/pkg/types"
)

// DeploymentOrchestratorResolver resolves the correct IDeploymentOrchestrator based on the type.
type DeploymentOrchestratorResolver struct {
	orchestrators map[string]types.IDeploymentOrchestrator
}

// NewDeploymentOrchestratorResolver creates a new DeploymentOrchestratorResolver.
func NewDeploymentOrchestratorResolver(orchestrators []types.IDeploymentOrchestrator) (*DeploymentOrchestratorResolver, error) {
	orchestratorMap := make(map[string]types.IDeploymentOrchestrator)

	// Add orchestrators to the map
	for _, orchestrator := range orchestrators {
		for _, orchestratorType := range orchestrator.DeploymentTypes() {
			if _, exists := orchestratorMap[orchestratorType]; exists {
				return nil, fmt.Errorf("type %s is duplicated", orchestratorType)
			}
			orchestratorMap[orchestratorType] = orchestrator
		}
	}

	return &DeploymentOrchestratorResolver{
		orchestrators: orchestratorMap,
	}, nil
}

// Resolve finds the correct IDeploymentOrchestrator based on the specification type.
func (r *DeploymentOrchestratorResolver) Resolve(specification types.IDeploymentSpecification) (types.IDeploymentOrchestrator, error) {
	orchestrator, exists := r.orchestrators[specification.Type()]
	if !exists {
		return nil, errors.New("orchestrator not found for type: " + specification.Type())
	}
	return orchestrator, nil
}
