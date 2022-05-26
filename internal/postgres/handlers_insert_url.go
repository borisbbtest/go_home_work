package postgres

import (
	"context"
)

const (
	keyPostgresInserts = "pgsql.insert.tb.url"
)

// connectionsHandler executes select from pg_stat_activity command and returns JSON if all is OK or nil otherwise.
func (p *Plugin) insertURLHandler(conn *postgresConn, key string, params []interface{}) (interface{}, error) {

	var err error
	query := `INSERT INTO public."storeurl"(
		"Port", "URL", "Path", "ShortPath", "UserID")
		VALUES ($1, $2, $3, $4, $5);`

	_, err = conn.postgresPool.Exec(context.Background(), query, params[0], params[1], params[2], params[3], params[4])
	if err != nil {
		return nil, err
	}

	return nil, nil
}
