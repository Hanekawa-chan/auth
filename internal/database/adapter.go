package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Hanekawa-chan/kanji-auth/internal/app"
	"github.com/dlmiddlecote/sqlstats"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/qustavo/sqlhooks/v2"
	"github.com/rs/zerolog"
	"time"
)

type adapter struct {
	logger *zerolog.Logger
	config *app.Config
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

func NewAdapter(logger *zerolog.Logger, config *app.Config) (app.Database, error) {
	sql.Register("postgresWrapped", sqlhooks.Wrap(&pq.Driver{}, &Hooks{}))
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.Name, config.DB.Password)
	db, err := sqlx.Connect("postgresWrapped", dsn)

	// Create a new collector, the name will be used as a label on the metrics
	collector := sqlstats.NewStatsCollector("auth", db)

	// Register it with Prometheus
	prometheus.MustRegister(collector)

	a := &adapter{
		logger: logger,
		config: config,
		db:     db,
	}
	return a, err
}

var (
	ErrNotFound = errors.New("rows not found")
)
