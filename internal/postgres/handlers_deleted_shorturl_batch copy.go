package postgres

import (
	"context"

	"github.com/borisbbtest/go_home_work/internal/model"
)

const (
	keyPostgresDeletedShortURLBatch = "pgsql.deleted.tb.short.url.batch"
)

// connectionsHandler executes select from pg_stat_activity command and returns JSON if all is OK or nil otherwise.
func (p *Plugin) DeletedShortURLBatchURLHandler(conn *postgresConn, key string, params []interface{}) (interface{}, error) {

	var err error

	ft := params[0].([]model.DataURL)
	query := `UPDATE  public."storeurl"
	          SET "StatusActive"= 2
			  WHERE ShortPath = $1;`

	for _, v := range ft {
		_, err = conn.postgresPool.Exec(context.Background(), query, v.ShortPath)
		if err != nil {
			log.Error(err)
		}
	}

	return nil, nil
}
