package api

import "github.com/cryptellation/exchanges/pkg/exchange"

const (
	// WorkerTaskQueueName is the name of the task queue for the cryptellation worker.
	WorkerTaskQueueName = "CryptellationExchangesTaskQueue"
)

const (
	// GetExchangeWorkflowName is the name of the workflow to get an exchange.
	GetExchangeWorkflowName = "GetExchangeWorkflow"
)

type (
	// GetExchangeWorkflowParams is the parameters of the GetExchange workflow.
	GetExchangeWorkflowParams struct {
		Name string
	}

	// GetExchangeWorkflowResults is the result of the GetExchange workflow.
	GetExchangeWorkflowResults struct {
		Exchange exchange.Exchange
	}
)

const (
	// ListExchangesWorkflowName is the name of the workflow to list exchanges.
	ListExchangesWorkflowName = "ListExchangesWorkflow"
)

type (
	// ListExchangesWorkflowParams is the parameters of the ListExchanges workflow.
	ListExchangesWorkflowParams struct{}

	// ListExchangesWorkflowResults is the result of the ListExchanges workflow.
	ListExchangesWorkflowResults struct {
		List []string
	}
)

const (
	// ServiceInfoWorkflowName is the name of the workflow to get the service info.
	ServiceInfoWorkflowName = "ServiceInfoWorkflow"
)

type (
	// ServiceInfoParams contains the parameters of the service info workflow.
	ServiceInfoParams struct{}

	// ServiceInfoResults contains the result of the service info workflow.
	ServiceInfoResults struct {
		Version string
	}
)
