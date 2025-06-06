package sql

import (
	"context"
	"fmt"

	"github.com/cryptellation/exchanges/pkg/exchange"
	"github.com/cryptellation/exchanges/svc/db"
	"github.com/cryptellation/exchanges/svc/db/sql/entities"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostGres driver
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

var _ db.DB = (*Activities)(nil)

// Activities is a struct that contains all the methods to interact with the
// activities table in the database.
type Activities struct {
	db *sqlx.DB
}

// New creates new SQL activities.
func New(ctx context.Context, dsn string) (*Activities, error) {
	// Create embedded database access
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Create a structure
	a := &Activities{
		db: db,
	}

	return a, nil
}

// Register registers the activities.
func (a *Activities) Register(w worker.Worker) {
	w.RegisterActivityWithOptions(
		a.CreateExchangesActivity,
		activity.RegisterOptions{Name: db.CreateExchangesActivityName},
	)
	w.RegisterActivityWithOptions(
		a.ReadExchangesActivity,
		activity.RegisterOptions{Name: db.ReadExchangesActivityName},
	)
	w.RegisterActivityWithOptions(
		a.UpdateExchangesActivity,
		activity.RegisterOptions{Name: db.UpdateExchangesActivityName},
	)
	w.RegisterActivityWithOptions(
		a.DeleteExchangesActivity,
		activity.RegisterOptions{Name: db.DeleteExchangesActivityName},
	)
}

// Reset will reset the database.
func (a *Activities) Reset(ctx context.Context) error {
	_, err := a.db.ExecContext(ctx, "DELETE FROM exchanges")
	if err != nil {
		return fmt.Errorf("deleting exchanges rows: %w", err)
	}

	return nil
}

// CreateExchangesActivity will create the exchanges in the database.
func (a *Activities) CreateExchangesActivity(
	ctx context.Context,
	params db.CreateExchangesActivityParams,
) (db.CreateExchangesActivityResults, error) {
	// Change models to entities
	exchanges := make([]entities.Exchange, len(params.Exchanges))
	for i, m := range params.Exchanges {
		e, err := entities.ExchangeFromModel(m)
		if err != nil {
			return db.CreateExchangesActivityResults{}, fmt.Errorf("creating exchange entity: %w", err)
		}
		exchanges[i] = e
	}

	// Insert the exchanges
	_, err := a.db.NamedExecContext(
		ctx,
		`INSERT INTO exchanges (name, data)
		VALUES (:name, :data)`,
		entities.FromEntitiesToMap(exchanges),
	)
	if err != nil {
		return db.CreateExchangesActivityResults{}, fmt.Errorf("inserting exchanges: %w", err)
	}

	return db.CreateExchangesActivityResults{}, nil
}

// ReadExchangesActivity will read the exchanges from the database.
func (a *Activities) ReadExchangesActivity(
	ctx context.Context,
	params db.ReadExchangesActivityParams,
) (db.ReadExchangesActivityResults, error) {
	var query string
	var args []interface{}
	var err error

	// Build the exchanges query
	if len(params.Names) == 0 {
		query = "SELECT name, data FROM exchanges"
	} else {
		query, args, err = sqlx.In("SELECT name, data FROM exchanges WHERE name IN (?)", params.Names)
	}
	if err != nil {
		return db.ReadExchangesActivityResults{}, fmt.Errorf("building query: %w", err)
	}

	// Execute the query
	query = a.db.Rebind(query)
	rows, err := a.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return db.ReadExchangesActivityResults{}, fmt.Errorf("querying exchanges: %w", err)
	}
	defer rows.Close()

	// Convert the rows to entities
	exchanges := make([]entities.Exchange, 0)
	for rows.Next() {
		var e entities.Exchange
		if err := rows.StructScan(&e); err != nil {
			return db.ReadExchangesActivityResults{}, fmt.Errorf("scanning exchange: %w", err)
		}
		exchanges = append(exchanges, e)
	}

	// Convert the entities to models
	models := make([]exchange.Exchange, 0)
	for _, e := range exchanges {
		m := e.ToModel()
		models = append(models, m)
	}

	return db.ReadExchangesActivityResults{Exchanges: models}, nil
}

// UpdateExchangesActivity will update the exchanges in the database.
func (a *Activities) UpdateExchangesActivity(
	ctx context.Context,
	params db.UpdateExchangesActivityParams,
) (db.UpdateExchangesActivityResults, error) {
	// Change models to entities
	exchanges := make([]entities.Exchange, len(params.Exchanges))
	for i, m := range params.Exchanges {
		e, err := entities.ExchangeFromModel(m)
		if err != nil {
			return db.UpdateExchangesActivityResults{}, fmt.Errorf("updating exchange entity: %w", err)
		}
		exchanges[i] = e
	}

	// Update the exchanges
	for _, e := range exchanges {
		_, err := a.db.NamedExecContext(
			ctx,
			`UPDATE exchanges
			SET data = :data
			WHERE name = :name`,
			e,
		)
		if err != nil {
			return db.UpdateExchangesActivityResults{}, fmt.Errorf("updating exchange: %w", err)
		}
	}

	return db.UpdateExchangesActivityResults{}, nil
}

// DeleteExchangesActivity will delete the exchanges from the database.
func (a *Activities) DeleteExchangesActivity(
	ctx context.Context,
	params db.DeleteExchangesActivityParams,
) (db.DeleteExchangesActivityResults, error) {
	// Build query
	query, args, err := sqlx.In(
		"DELETE FROM exchanges WHERE name IN (?)",
		params.Names,
	)
	if err != nil {
		return db.DeleteExchangesActivityResults{}, fmt.Errorf("building query: %w", err)
	}

	// Execute query
	query = a.db.Rebind(query)
	_, err = a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return db.DeleteExchangesActivityResults{}, fmt.Errorf("deleting exchanges: %w", err)
	}

	return db.DeleteExchangesActivityResults{}, nil
}
