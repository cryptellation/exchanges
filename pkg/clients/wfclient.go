package clients

import (
	"github.com/cryptellation/exchanges/api"
	"go.temporal.io/sdk/workflow"
)

// WfClient is a client for the cryptellation exchanges service from a workflow perspective.
type WfClient interface {
	// GetExchange calls the exchange get workflow.
	GetExchange(
		ctx workflow.Context,
		params api.GetExchangeWorkflowParams,
		childWorkflowOptions *workflow.ChildWorkflowOptions,
	) (result api.GetExchangeWorkflowResults, err error)
}

type wfClient struct{}

// NewWfClient creates a new workflow client.
// This client is used to call workflows from within other workflows.
// It is not used to call workflows from outside the workflow environment.
func NewWfClient() WfClient {
	return wfClient{}
}

// GetExchange gets exchange info from Cryptellation service.
func (wfClient) GetExchange(
	ctx workflow.Context,
	params api.GetExchangeWorkflowParams,
	childWorkflowOptions *workflow.ChildWorkflowOptions,
) (result api.GetExchangeWorkflowResults, err error) {
	// Set options
	if childWorkflowOptions == nil {
		childWorkflowOptions = &workflow.ChildWorkflowOptions{}
	}
	ctx = workflow.WithChildOptions(ctx, *childWorkflowOptions)

	// Get exchange info
	err = workflow.ExecuteChildWorkflow(ctx, api.GetExchangeWorkflowName, params).Get(ctx, &result)
	return result, err
}
