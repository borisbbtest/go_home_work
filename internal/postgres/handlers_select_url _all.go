package postgres

import (
	"context"
	"fmt"

	"github.com/borisbbtest/go_home_work/internal/model"
)

const (
	keyPostgresAllSelectURL = "pgsql.select.tb.all.url"
)

// ЗАпрос на предоставления короткого линка.
func (p *Plugin) selectAllURLHandler(conn *postgresConn, key string, params []interface{}) (interface{}, error) {
	buff := []model.ResponseURL{}

	query := `SELECT "URL", "ShortPath" FROM  "storeurl"  WHERE  "UserID"  = $1;`

	rows, err := conn.postgresPool.Query(context.Background(), query, params[0])

	for rows.Next() {
		m := model.ResponseURL{}

		err = rows.Scan(&m.OriginalURL, &m.ShortURL)
		m.ShortURL = fmt.Sprintf("%s/%s", params[1], m.ShortURL)
		if err != nil {
			log.Info("URLs - c: ", err)
			return nil, err
		}
		buff = append(buff, m)
	}
	if err != nil {
		log.Error(err)
		return postgresPingFailed, err
	}

	return buff, nil
}
