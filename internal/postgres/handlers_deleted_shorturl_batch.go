package postgres

import (
	"context"

	"github.com/borisbbtest/go_home_work/internal/model"
	"github.com/jackc/pgx/v4"
)

const (
	keyPostgresDeletedShortURLBatch = "pgsql.deleted.tb.short.url.batch"
)

func (p *Plugin) DeletedShortURLBatchURLHandler(conn *postgresConn, key string, params []interface{}) (interface{}, error) {

	ft := params[0].([]model.DataURL)

	ctx := context.Background()

	tx, err := conn.postgresPool.Begin(ctx)
	b := &pgx.Batch{}

	for _, v := range ft {
		query := `UPDATE  public."storeurl"
		SET "StatusActive"= 2
		WHERE "ShortPath" = $1;`
		b.Queue(query, v.ShortPath)
	}

	batchResults := tx.SendBatch(ctx, b)

	var qerr error
	var rows pgx.Rows
	for qerr == nil {
		rows, qerr = batchResults.Query()
		rows.Close()
	}
	tx.Commit(ctx)

	return nil, err
}
