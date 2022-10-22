package postgres

import (
	"context"
)

const (
	keyPostgresSelectCountUsers = "pgsql.select.tb.users.count"
)

// ЗАпрос на предоставления короткого линка.
func (p *Plugin) selectCountUsersHandler(conn *postgresConn, key string, params []interface{}) (interface{}, error) {

	var buff int32
	query := `SELECT count(distinct "UserID") FROM  "storeurl" group by "UserID" ;`

	err := conn.postgresPool.QueryRow(context.Background(), query).Scan(&buff)
	if err != nil {
		log.Error("SELECT count(User)", err)
		return 0, err
	}

	return buff, nil
}
