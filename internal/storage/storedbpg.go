package storage

import (
	"github.com/borisbbtest/go_home_work/internal/model"
	"github.com/borisbbtest/go_home_work/internal/postgres"
)

type StoreDBinPostgreSQL struct {
	pgp     postgres.Plugin
	connStr string
}

func NewPostgreSQLStorage(connStr string) (res *StoreDBinPostgreSQL, err error) {
	res = &StoreDBinPostgreSQL{}
	res.connStr = connStr
	res.pgp.Start()
	_, err = res.pgp.NewDBConn("pgsql.create.tb.url", []string{}, connStr)
	return
}

func (hook *StoreDBinPostgreSQL) Put(k string, v DataURL) error {

	_, err := hook.pgp.NewDBConn("pgsql.create.tb.url", []string{}, hook.connStr)
	if err != nil {
		log.Error("pgsql.create.tb.url", err)
	}
	return err
}

func (hook *StoreDBinPostgreSQL) Get(k string) (DataURL, error) {

	return DataURL{}, nil
}

func (hook *StoreDBinPostgreSQL) GetAll(k string, dom string) ([]model.ResponseURL, error) {

	return []model.ResponseURL{}, nil
}

func (hook *StoreDBinPostgreSQL) Close() {
	hook.pgp.Stop()
}
