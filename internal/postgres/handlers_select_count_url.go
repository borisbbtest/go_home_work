package postgres

import (
	"context"
)

const (
	keyPostgresSelectCountURL = "pgsql.select.tb.url.count"
)

// ЗАпрос на предоставления короткого линка.
func (p *Plugin) selectCountURLlHandler(conn *postgresConn, key string, params []interface{}) (interface{}, error) {

	var buff int32
	query := `SELECT count(distinct "ShortPath")  FROM  "storeurl" ;`

	err := conn.postgresPool.QueryRow(context.Background(), query).Scan(&buff)
	if err != nil {
		log.Error("SELECT count( URL)", err)
		return 0, err
	}

	return buff, nil
}
