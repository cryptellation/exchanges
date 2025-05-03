package clients

import (
	"context"

	"github.com/cryptellation/exchanges/api"
	temporalclient "go.temporal.io/sdk/client"
)

// Client is a client for the cryptellation exchanges service.
type Client interface {
	// GetExchange calls the exchange get workflow.
	GetExchange(ctx context.Context, params api.GetExchangeWorkflowParams) (api.GetExchangeWorkflowResults, error)
	// ListExchanges calls the exchanges list workflow.
	ListExchanges(ctx context.Context, params api.ListExchangesWorkflowParams) (api.ListExchangesWorkflowResults, error)
	// Info calls the service info.
	Info(ctx context.Context) (api.ServiceInfoResults, error)
}

type client struct {
	temporal temporalclient.Client
}

// New creates a new client to execute temporal workflows.
func New(cl temporalclient.Client) Client {
	return &client{temporal: cl}
}

// GetExchange calls the exchange get workflow.
func (c client) GetExchange(
	ctx context.Context,
	params api.GetExchangeWorkflowParams,
) (res api.GetExchangeWorkflowResults, err error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.GetExchangeWorkflowName, params)
	if err != nil {
		return api.GetExchangeWorkflowResults{}, err
	}

	// Get result and return
	err = exec.Get(ctx, &res)
	return res, err
}

// ListExchanges calls the exchanges list workflow.
func (c client) ListExchanges(
	ctx context.Context,
	params api.ListExchangesWorkflowParams,
) (res api.ListExchangesWorkflowResults, err error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.ListExchangesWorkflowName, params)
	if err != nil {
		return api.ListExchangesWorkflowResults{}, err
	}

	// Get result and return
	err = exec.Get(ctx, &res)
	return res, err
}

// Info calls the service info.
func (c client) Info(ctx context.Context) (res api.ServiceInfoResults, err error) {
	workflowOptions := temporalclient.StartWorkflowOptions{
		TaskQueue: api.WorkerTaskQueueName,
	}

	// Execute workflow
	exec, err := c.temporal.ExecuteWorkflow(ctx, workflowOptions, api.ServiceInfoWorkflowName)
	if err != nil {
		return api.ServiceInfoResults{}, err
	}

	// Get result and return
	err = exec.Get(ctx, &res)
	return res, err
}
