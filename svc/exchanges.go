package svc

import (
	"github.com/cryptellation/exchanges/api"
	"github.com/cryptellation/exchanges/svc/db"
	exchangesadapter "github.com/cryptellation/exchanges/svc/exchanges"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// Exchanges is the exchanges domain.
type Exchanges interface {
	Register(w worker.Worker)

	GetExchangeWorkflow(
		ctx workflow.Context,
		params api.GetExchangeWorkflowParams,
	) (api.GetExchangeWorkflowResults, error)

	ListExchangesWorkflow(
		ctx workflow.Context,
		params api.ListExchangesWorkflowParams,
	) (api.ListExchangesWorkflowResults, error)
}

// Check that the workflows implements the Exchanges interface.
var _ Exchanges = &workflows{}

type workflows struct {
	db        db.DB
	exchanges exchangesadapter.Exchanges
}

// New creates a new exchanges workflows.
func New(db db.DB, exchs exchangesadapter.Exchanges) Exchanges {
	return &workflows{
		db:        db,
		exchanges: exchs,
	}
}

// Register registers the candlesticks workflows to the worker.
func (wf *workflows) Register(w worker.Worker) {
	w.RegisterWorkflowWithOptions(wf.GetExchangeWorkflow, workflow.RegisterOptions{
		Name: api.GetExchangeWorkflowName,
	})
	w.RegisterWorkflowWithOptions(wf.ListExchangesWorkflow, workflow.RegisterOptions{
		Name: api.ListExchangesWorkflowName,
	})

	w.RegisterWorkflowWithOptions(ServiceInfoWorkflow, workflow.RegisterOptions{
		Name: api.ServiceInfoWorkflowName,
	})
}
