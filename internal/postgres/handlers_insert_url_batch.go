package postgres

import (
	"context"

	"github.com/borisbbtest/go_home_work/internal/model"
	"github.com/jackc/pgx/v4"
)

const (
	keyPostgresInsertsBatch = "pgsql.insert.tb.url.batch"
)

func (p *Plugin) insertBatchURLHandler(conn *postgresConn, key string, params []interface{}) (interface{}, error) {

	ft := params[0].([]model.DataURL)

	b := &pgx.Batch{}

	for _, v := range ft {
		query := `INSERT INTO public."storeurl"(
			"Port", "URL", "Path", "ShortPath", "UserID", "CorrelationId","StatusActive")
			VALUES ($1, $2, $3, $4, $5, $6, 1);`
		b.Queue(query, v.Port, v.URL, v.Path, v.ShortPath, v.UserID, v.CorrelationID)
	}
	batchResults := conn.postgresPool.SendBatch(context.Background(), b)

	var qerr error
	var rows pgx.Rows
	for qerr == nil {
		rows, qerr = batchResults.Query()
		rows.Close()
	}

	return nil, qerr
}
