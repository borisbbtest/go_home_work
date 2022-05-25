package postgres

import (
	"context"
)

const (
	keyPostgresCreateDdURL = "pgsql.create.db.url"
)

// pingHandler executes 'SELECT 1 as pingOk' commands and returns pingOK if a connection is alive or postgresPingFailed otherwise.
func (p *Plugin) CreateDdURLHandler(conn *postgresConn, key string, params []string) (interface{}, error) {
	query := `
				CREATE TABLE IF NOT EXISTS public."StoreURL"
					(
						"Port" "char"[],
						"URL" "char"[] NOT NULL,
						"Path" "char"[],
						"ShortPath" "char"[] NOT NULL,
						"UserID" "char"[],
						CONSTRAINT "StoreURL_pkey" PRIMARY KEY ("ShortPath")
					)

					TABLESPACE pg_default;

					ALTER TABLE IF EXISTS public."StoreURL"
						OWNER to postgres;

					COMMENT ON TABLE public."StoreURL"
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
