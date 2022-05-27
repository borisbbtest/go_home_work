package postgres

import (
	"context"
)

const (
	keyPostgresCreateDdURL = "pgsql.create.tb.url"
)

// pingHandler executes 'SELECT 1 as pingOk' commands and returns pingOK if a connection is alive or postgresPingFailed otherwise.
func (p *Plugin) CreateTbURLHandler(conn *postgresConn, key string, params []interface{}) (interface{}, error) {
	query := `
				CREATE TABLE IF NOT EXISTS public."storeurl"
					(
						"Port" "text",
						"URL" "text" NOT NULL,
						"Path" "text",
						"ShortPath" "text" NOT NULL,
						"UserID" "text",
						"CorrelationId" "text",
						CONSTRAINT "storeurl_pkey" PRIMARY KEY ("URL")
					)

					TABLESPACE pg_default;

					ALTER TABLE IF EXISTS public."storeurl"
						OWNER to postgres;

					COMMENT ON TABLE public."storeurl"
						IS '	Port      string json:"Port"
						URL       string json:"URL"
						Path      string json:"Path"
						ShortPath string json:"ShortPath"
						UserID    string json:"UserID"';
			`

	if _, err := conn.postgresPool.Exec(context.Background(), query); err != nil {
		return 0, err
	}
	return 1, nil
}
