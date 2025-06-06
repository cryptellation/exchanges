package svc

import (
	"github.com/cryptellation/exchanges/api"
	exchangesactivity "github.com/cryptellation/exchanges/svc/exchanges"
	"go.temporal.io/sdk/workflow"
)

// ListExchangesWorkflow will list the exchanges.
func (wf *workflows) ListExchangesWorkflow(
	ctx workflow.Context,
	_ api.ListExchangesWorkflowParams,
) (api.ListExchangesWorkflowResults, error) {
	// Get the list of exchanges from the services
	var r exchangesactivity.ListExchangesActivityResults
	err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, exchangesactivity.DefaultActivityOptions()),
		wf.exchanges.ListExchangesActivity,
		exchangesactivity.ListExchangesActivityParams{}).Get(ctx, &r)
	if err != nil {
		return api.ListExchangesWorkflowResults{}, err
	}

	// Return the result
	return api.ListExchangesWorkflowResults{
		List: r.List,
	}, nil
}
