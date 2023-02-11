package database

import (
	"auth/internal/app"
	"context"
	"database/sql"
	"fmt"
	"github.com/dlmiddlecote/sqlstats"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/qustavo/sqlhooks/v2"
	"github.com/rs/zerolog"
	"time"
)

type adapter struct {
	logger *zerolog.Logger
	config *Config
	db     *sqlx.DB
}

type Hooks struct {
}

// Before hook will print the query with it's args and return the context with the timestamp
func (h *Hooks) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	fmt.Printf("> %s %q", query, args)
	return context.WithValue(ctx, "begin", time.Now()), nil
}

// After hook will get the timestamp registered on the Before hook and print the elapsed time
func (h *Hooks) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	begin := ctx.Value("begin").(time.Time)
	fmt.Printf(". took: %s\n", time.Since(begin))
	return ctx, nil
}

func NewAdapter(logger *zerolog.Logger, config *Config) (app.Database, error) {
	sql.Register("postgresWrapped", sqlhooks.Wrap(&stdlib.Driver{}, &Hooks{}))
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s "+
		"sslmode=disable search_path=public default_query_exec_mode=cache_describe",
		config.Host, config.Port, config.User, config.Name, config.Password)
	db, err := sqlx.Connect("postgresWrapped", dsn)
	if err != nil {
		logger.Err(err).Msg("db connect")
		return nil, err
	}

	// Create a new collector, the name will be used as a label on the metrics
	collector := sqlstats.NewStatsCollector("auth", db)

	// Register it with Prometheus
	prometheus.MustRegister(collector)

	instance, err := postgres.WithInstance(db.DB, &postgres.Config{DatabaseName: config.Name, SchemaName: "public"})
	if err != nil {
		logger.Err(err).Msg("db instance")
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		config.MigrationsURL, config.Name, instance)
	if err != nil {
		logger.Err(err).Msg("db migration create")
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Err(err).Msg("db migrate up")
		return nil, err
	}

	a := &adapter{
		logger: logger,
		config: config,
		db:     db,
	}
	return a, err
}
