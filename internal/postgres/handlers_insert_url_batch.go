package postgres

import (
	"context"

	"github.com/borisbbtest/go_home_work/internal/model"
)

const (
	keyPostgresInsertsBatch = "pgsql.insert.tb.url.batch"
)

// connectionsHandler executes select from pg_stat_activity command and returns JSON if all is OK or nil otherwise.
func (p *Plugin) insertBatchURLHandler(conn *postgresConn, key string, params []interface{}) (interface{}, error) {

	var err error

	ft := params[0].([]model.DataURL)
	query := `INSERT INTO public."storeurl"(
		"Port", "URL", "Path", "ShortPath", "UserID", "CorrelationId")
		VALUES ($1, $2, $3, $4, $5, $6);`

	for _, v := range ft {
		_, err = conn.postgresPool.Exec(context.Background(), query, v.Port, v.URL, v.Path, v.ShortPath, v.UserID, v.CorrelationId)
		if err != nil {
			log.Error(err)
		}
	}

	return nil, nil
}
