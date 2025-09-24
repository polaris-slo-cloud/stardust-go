package deployment

import (
	"errors"
	"fmt"

	"github.com/keniack/stardustGo/pkg/types"
)

// DeploymentOrchestratorResolver resolves the correct IDeploymentOrchestrator based on the type.
type DeploymentOrchestratorResolver struct {
	orchestrators map[string]types.DeploymentOrchestrator
}

// NewDeploymentOrchestratorResolver creates a new DeploymentOrchestratorResolver.
func NewDeploymentOrchestratorResolver(orchestrators []types.DeploymentOrchestrator) (*DeploymentOrchestratorResolver, error) {
	orchestratorMap := make(map[string]types.DeploymentOrchestrator)

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
func (r *DeploymentOrchestratorResolver) Resolve(specification types.DeploymentSpecification) (types.DeploymentOrchestrator, error) {
	orchestrator, exists := r.orchestrators[specification.Type()]
	if !exists {
		return nil, errors.New("orchestrator not found for type: " + specification.Type())
	}
	return orchestrator, nil
}
