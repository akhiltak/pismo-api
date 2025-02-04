package bunorm

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/pandoratoolbox/bun/extra/bunslog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func Connect(ctx context.Context, dns string, ping bool) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dns)))
	if err := sqldb.PingContext(ctx); err != nil {
		panic(fmt.Sprintf("error in connecting to db: %s", err))
	}
	db := bun.NewDB(sqldb, pgdialect.New())
	hook := bunslog.NewQueryHook(
		bunslog.WithQueryLogLevel(slog.LevelInfo),
		bunslog.WithSlowQueryLogLevel(slog.LevelWarn),
		bunslog.WithErrorQueryLogLevel(slog.LevelError),
		bunslog.WithSlowQueryThreshold(3*time.Second),
	)
	db.AddQueryHook(hook)
	if ping {
		if err := db.Ping(); err != nil {
			panic(fmt.Sprintf("error pinging to database %s", err))
		}
	}
	return db
}
