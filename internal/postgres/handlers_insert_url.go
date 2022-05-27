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
	var short_url string
	query := `
	WITH cte AS (
		INSERT INTO public."storeurl"(
		"Port", "URL", "Path", "ShortPath", "UserID")
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT ("URL") DO NOTHING
		RETURNING "URL"
	)
    SELECT "ShortPath"  FROM  "storeurl"  WHERE  "URL"  = $1;`

	err = conn.postgresPool.QueryRow(context.Background(), query, params[0], params[1], params[2], params[3], params[4]).Scan(&short_url)
	log.Info("---->", short_url, "<---")
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return nil, nil
}
