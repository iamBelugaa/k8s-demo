package database

import (
	"context"
	"database/sql"
	"net/url"
	"time"

	"github.com/iamBelugaa/k8s-demo/internal/config"
	"github.com/iamBelugaa/k8s-demo/pkg/logger"
	_ "github.com/lib/pq"
)

func Open(cfg *config.DB) (*sql.DB, error) {
	q := url.Values{}
	q.Set("sslmode", cfg.TLS)

	u := url.URL{
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
		Scheme:   cfg.Scheme,
		User:     url.UserPassword(cfg.User, cfg.Password),
	}

	db, err := sql.Open("postgres", u.String())
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}

func StatusCheck(ctx context.Context, db *sql.DB, log *logger.Logger) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second*10)
		defer cancel()
	}

	for attempts := 1; ; attempts++ {
		if err := db.PingContext(ctx); err == nil {
			break
		} else {
			log.WithTrace(ctx).Infow("db ping error", "error", err)
		}

		time.Sleep(time.Duration(attempts) * 200 * time.Millisecond)

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	const q = `SELECT TRUE`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
