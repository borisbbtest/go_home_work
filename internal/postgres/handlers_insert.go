package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
)

const (
	keyPostgresConnections = "pgsql.connections"
)

// connectionsHandler executes select from pg_stat_activity command and returns JSON if all is OK or nil otherwise.
func (p *Plugin) connectionsHandler(conn *postgresConn, key string, params []string) (interface{}, error) {
	var connectionsJSON string
	var err error
	query := `SELECT row_to_json(T)
	FROM (
		SELECT
			sum(CASE WHEN state = 'active' THEN 1 ELSE 0 END) AS active,
			sum(CASE WHEN state = 'idle' THEN 1 ELSE 0 END) AS idle,
			sum(CASE WHEN state = 'idle in transaction' THEN 1 ELSE 0 END) AS idle_in_transaction,
			sum(CASE WHEN state = 'idle in transaction (aborted)' THEN 1 ELSE 0 END) AS idle_in_transaction_aborted,
			sum(CASE WHEN state = 'fastpath function call' THEN 1 ELSE 0 END) AS fastpath_function_call,
			sum(CASE WHEN state = 'disabled' THEN 1 ELSE 0 END) AS disabled,
			count(*) AS total,
			count(*)*100/(SELECT current_setting('max_connections')::int) AS total_pct,
			sum(CASE WHEN wait_event IS NOT NULL THEN 1 ELSE 0 END) AS waiting,
			(SELECT count(*) FROM pg_prepared_xacts) AS prepared
		FROM pg_stat_activity WHERE datid is not NULL) T;`

	err = conn.postgresPool.QueryRow(context.Background(), query).Scan(&connectionsJSON)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error(err.Error())
			return nil, errorEmptyResult
		}
		log.Error(err.Error())
		return nil, errorCannotFetchData
	}

	return connectionsJSON, nil
}
