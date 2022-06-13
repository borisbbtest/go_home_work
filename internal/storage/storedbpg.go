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
	_, err = res.pgp.NewDBConn("pgsql.create.tb.url", []string{}, connStr, []interface{}{})
	if err != nil {
		log.Error("pgsql.create.tb.url", err)
	}
	return
}

func (hook *StoreDBinPostgreSQL) Put(k string, v model.DataURL) (string, error) {
	buff := []interface{}{v.Port, v.URL, v.Path, v.ShortPath, v.UserID}
	res, err := hook.pgp.NewDBConn("pgsql.insert.tb.url", []string{}, hook.connStr, buff)
	if err != nil {
		return "", err
	}
	return res.(string), err
}
func (hook *StoreDBinPostgreSQL) PutBatch(k string, v []model.DataURL) error {

	buff := []interface{}{v}
	_, err := hook.pgp.NewDBConn("pgsql.insert.tb.url.batch", []string{}, hook.connStr, buff)
	if err != nil {
		log.Error("pgsql.insert.tb.url", err)
	}
	return err
}

func (hook *StoreDBinPostgreSQL) DeletedURLBatch(k string, v []model.DataURL) error {

	buff := []interface{}{v}
	_, err := hook.pgp.NewDBConn("pgsql.deleted.tb.short.url.batch", []string{}, hook.connStr, buff)
	if err != nil {
		log.Error("pgsql.deleted.tb.short.url.batch", err)
	}
	return err
}

func (hook *StoreDBinPostgreSQL) Get(k string) (model.DataURL, error) {

	buff := []interface{}{k}
	res, err := hook.pgp.NewDBConn("pgsql.select.tb.url", []string{}, hook.connStr, buff)
	if err != nil {
		log.Error("pgsql.select.tb.url", err)
		return model.DataURL{}, err
	}

	return res.(model.DataURL), nil
}

func (hook *StoreDBinPostgreSQL) GetAll(k string, dom string) ([]model.ResponseURL, error) {

	return []model.ResponseURL{}, nil
}

func (hook *StoreDBinPostgreSQL) Close() {
	hook.pgp.Stop()
}
