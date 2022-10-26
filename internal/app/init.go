package app

import (
	"github.com/borisbbtest/go_home_work/internal/config"
	"github.com/borisbbtest/go_home_work/internal/storage"
)

type ServiceShortURL struct {
	ServerConf *config.ServiceShortURLConfig
	Storage    storage.Storage
}

func Init(cfg *config.ServiceShortURLConfig) (res *ServiceShortURL, err error) {

	res = &ServiceShortURL{}
	res.Storage, err = storage.NewPostgreSQLStorage(cfg.DataBaseDSN)
	if err != nil {
		res.Storage, err = storage.NewFileStorage(cfg.FileStorePath)
		if err != nil {
			log.Error(err)
			return
		}
	}
	res.ServerConf = cfg
	return
}

func (hook *ServiceShortURL) Start() (err error) {

	log.Info("Start RPC")
	go NewRPC(hook.ServerConf, hook.Storage).Start()

	log.Info("Start HTTP")
	err = NewHTTP(hook.ServerConf, hook.Storage).Start()
	if err != nil {
		log.Fatal(err)
		return
	}
	return

}
