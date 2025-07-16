package binance

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	client "github.com/adshao/go-binance/v2"
	"github.com/cryptellation/exchanges/pkg/exchange"
	"github.com/cryptellation/exchanges/svc/exchanges"
	"go.temporal.io/sdk/worker"
)

const (
	exchangeName = "binance"
)

var (
	info = exchange.Exchange{
		Name: exchangeName,
		Periods: []string{
			"M1", "M3", "M5", "M15", "M30",
			"H1", "H2", "H4", "H6", "H8", "H12",
			"D1", "D3",
			"W1",
		},
		Fees: 0.1,
	}
)

// Activities is the live Binance activities.
type Activities struct {
	Client *client.Client
}

// New creates a new Binance activities.
func New(apiKey, secretKey string) (*Activities, error) {
	// Validate that API key and secret key are not empty
	if apiKey == "" {
		return nil, errors.New("API key cannot be empty")
	}
	if secretKey == "" {
		return nil, errors.New("secret key cannot be empty")
	}

	c := client.NewClient(apiKey, secretKey)
	c.Logger.SetOutput(io.Discard)

	// Return service
	return &Activities{
		Client: c,
	}, nil
}

// Name returns the name of the Binance activities.
func (a *Activities) Name() string {
	return exchangeName
}

// Register registers the Binance activities with the given worker.
func (a *Activities) Register(_ worker.Worker) {
	// No need to register the Binance activities, they are already registered
	// with its parent.
}

// ListExchangesActivity returns the names of the exchanges.
func (a *Activities) ListExchangesActivity(
	_ context.Context,
	_ exchanges.ListExchangesActivityParams,
) (exchanges.ListExchangesActivityResults, error) {
	return exchanges.ListExchangesActivityResults{
		List: []string{
			exchangeName,
		},
	}, nil
}

// GetExchangeActivity returns the exchange information for the given exchange.
func (a *Activities) GetExchangeActivity(
	ctx context.Context,
	_ exchanges.GetExchangeActivityParams,
) (exchanges.GetExchangeActivityResults, error) {
	exchangeInfos, err := a.Client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		return exchanges.GetExchangeActivityResults{}, err
	}

	pairs := make([]string, len(exchangeInfos.Symbols))
	for i, bs := range exchangeInfos.Symbols {
		pairs[i] = fmt.Sprintf("%s-%s", bs.BaseAsset, bs.QuoteAsset)
	}

	exch := info
	exch.Pairs = pairs
	exch.LastSyncTime = time.Now()

	return exchanges.GetExchangeActivityResults{
		Exchange: exch,
	}, nil
}
