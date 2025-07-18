package main

import (
	"github.com/cryptellation/exchanges/dagger/internal/dagger"
)

// PostgresService returns a service running Postgres initialized for integration tests.
func PostgresService(dag *dagger.Client, sourceDir *dagger.Directory) *dagger.Service {
	// Create the Postgres container with its environment variables
	c := dag.Container().
		From("postgres:15-alpine").
		WithEnvVariable("POSTGRES_PASSWORD", "postgres").
		WithEnvVariable("POSTGRES_USER", "postgres").
		WithEnvVariable("PGUSER", "postgres").
		WithEnvVariable("PGPASSWORD", "postgres").
		WithEnvVariable("POSTGRES_DB", "postgres")

	// Mount the initialization SQL directory
	initSQLDir := sourceDir.Directory("deployments/docker-compose/postgresql")
	c = c.WithMountedDirectory("/docker-entrypoint-initdb.d", initSQLDir)

	// Expose the default Postgres port
	c = c.WithExposedPort(5432)

	return c.AsService()
}

// TemporalService returns a Temporal service configured for Postgres, mounting dynamic config,
// and waiting for Postgres.
func TemporalService(dag *dagger.Client, sourceDir *dagger.Directory, db *dagger.Service) *dagger.Service {
	// Build the Temporal container with the official temporal image
	container := dag.Container().From("temporalio/auto-setup:1.25")

	// Bind the shared Postgres service to the container
	container = container.WithServiceBinding("postgresql", db)
	container = container.WithEnvVariable("DB", "postgres12")
	container = container.WithEnvVariable("DB_PORT", "5432")
	container = container.WithEnvVariable("POSTGRES_USER", "temporal")
	container = container.WithEnvVariable("POSTGRES_PWD", "temporal")
	container = container.WithEnvVariable("POSTGRES_SEEDS", "postgresql")
	container = container.WithEnvVariable("BIND_ON_IP", "0.0.0.0")
	container = container.WithEnvVariable("TEMPORAL_BROADCAST_ADDRESS", "127.0.0.1")

	// Set the dynamic config file
	temporalConfigDir := sourceDir.Directory("deployments/docker-compose/temporal")
	container = container.WithEnvVariable("DYNAMIC_CONFIG_FILE_PATH", "config/dynamicconfig/development-sql.yaml")
	container = container.WithMountedDirectory("/etc/temporal/config/dynamicconfig", temporalConfigDir)

	// Expose the Temporal frontend port
	container = container.WithExposedPort(7233)

	return container.AsService()
}
