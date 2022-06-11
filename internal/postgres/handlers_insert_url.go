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
	var shortURL string
	query := `
	WITH cte AS (
		INSERT INTO public."storeurl"(
		"Port", "URL", "Path", "ShortPath", "UserID","StatusActive")
		VALUES ($1, $2, $3, $4, $5,1)
		ON CONFLICT ("URL","StatusActive") DO NOTHING
		RETURNING "URL"
	)
	SELECT NULL AS result
	WHERE EXISTS (SELECT 1 FROM cte)
	UNION ALL
    SELECT "ShortPath"  FROM  "storeurl"  WHERE  "URL"  = $2;`

	err = conn.postgresPool.QueryRow(context.Background(), query, params...).Scan(&shortURL)
	log.Info("---->", shortURL, "<---")
	if err != nil {
		return nil, err
	}

	return shortURL, nil
}
