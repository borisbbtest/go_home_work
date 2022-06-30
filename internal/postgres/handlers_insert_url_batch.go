package postgres

import (
	"context"
	"time"

	"github.com/borisbbtest/go_home_work/internal/model"
	"github.com/jackc/pgx/v4"
)

const (
	keyPostgresInsertsBatch = "pgsql.insert.tb.url.batch"
)

func (p *Plugin) insertBatchURLHandler(conn *postgresConn, key string, params []interface{}) (interface{}, error) {

	var err error

	ft := params[0].([]model.DataURL)

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()

	b := &pgx.Batch{}

	for _, v := range ft {
		query := `INSERT INTO public."storeurl"(
			"Port", "URL", "Path", "ShortPath", "UserID", "CorrelationId","StatusActive")
			VALUES ($1, $2, $3, $4, $5, $6, 1);`
		b.Queue(query, v.Port, v.URL, v.Path, v.ShortPath, v.UserID, v.CorrelationID)
	}
	conn.postgresPool.SendBatch(ctx, b)

	if err != nil {
		log.Error(err)
	}

	return nil, nil
}
