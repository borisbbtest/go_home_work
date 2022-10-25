package postgres

import (
	"context"

	"github.com/borisbbtest/go_home_work/internal/model"
	"github.com/jackc/pgx/v4"
)

const (
	keyPostgresInsertsBatch = "pgsql.insert.tb.url.batch"
)

// Массовая вставка линков
func (p *Plugin) insertBatchURLHandler(conn *postgresConn, key string, params []interface{}) (interface{}, error) {

	ft := params[0].([]model.DataURL)

	ctx := context.Background()

	tx, err := conn.postgresPool.Begin(ctx)
	b := &pgx.Batch{}

	for _, v := range ft {
		query := `INSERT INTO public."storeurl"(
			"Port", "URL", "Path", "ShortPath", "UserID", "CorrelationId","StatusActive")
			VALUES ($1, $2, $3, $4, $5, $6, 1);`
		b.Queue(query, v.Port, v.URL, v.Path, v.ShortPath, v.UserID, v.CorrelationID)
	}

	batchResults := tx.SendBatch(ctx, b)

	var qerr error
	var rows pgx.Rows
	for qerr == nil {
		rows, qerr = batchResults.Query()
		log.Info(qerr)
		rows.Close()
	}
	tx.Commit(ctx)

	return nil, err
}
