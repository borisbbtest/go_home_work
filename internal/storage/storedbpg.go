package storage

import "github.com/borisbbtest/go_home_work/internal/model"

type StoreDBinPostgreSQL struct {
}

func NewPostgreSQLStorage(filename string) (res *StoreDBinFile, err error) {
	return
}

func (hook *StoreDBinPostgreSQL) Put(k string, v DataURL) error {

	return nil
}

func (hook *StoreDBinPostgreSQL) Get(k string) (DataURL, error) {

	return DataURL{}, nil
}

func (hook *StoreDBinPostgreSQL) GetAll(k string, dom string) ([]model.ResponseURL, error) {

	return []model.ResponseURL{}, nil
}

func (hook *StoreDBinPostgreSQL) Close() {

}
