package postgres

import (
	"context"

	"github.com/borisbbtest/go_home_work/internal/model"
)

const (
	keyPostgresSelectURLtoURL = "pgsql.select.tb.url.then"
)

// Запрос на вывод из БД  по users
func (p *Plugin) selectURLtoURLHandler(conn *postgresConn, key string, params []interface{}) (interface{}, error) {

	buff := model.DataURL{}
	query := `SELECT "Port", "URL", "Path", "ShortPath", "UserID", "StatusActive" FROM  "storeurl"  WHERE  "ShortPath"  = $1`

	err := conn.postgresPool.QueryRow(context.Background(), query, params[0]).Scan(&buff.Port, &buff.URL, &buff.Path, &buff.ShortPath, &buff.UserID, &buff.StatusActive)
	if err != nil {
		log.Error(err)
		return postgresPingFailed, err
	}

	return buff, nil
}
