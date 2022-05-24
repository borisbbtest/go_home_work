package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/borisbbtest/go_home_work/internal/model"
	"github.com/borisbbtest/go_home_work/internal/tools"
)

type StoreDBinFile struct {
	ReadURL  *tools.Consumer
	WriteURL *tools.Producer
	DB       map[string]DataURL
	ListUser map[string][]string
}

func (hook *StoreDBinFile) WriteEvent(event *DataURL) error {
	return hook.WriteURL.Encoder.Encode(&event)
}

func (hook *StoreDBinFile) ReadEvent() (*DataURL, error) {
	event := &DataURL{}
	if err := hook.ReadURL.Decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}

func (hook StoreDBinFile) RestoreDdBackupURL() {

	scanner := bufio.NewScanner(hook.ReadURL.GetFile())
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		var m DataURL
		if err := json.Unmarshal(scanner.Bytes(), &m); err != nil {
			log.Errorf("body error: %v", string(scanner.Bytes()))
			log.Errorf("error decoding message: %v", err)
		}
		hook.DB[m.ShortPath] = m
	}

}

func (hook StoreDBinFile) Put(k string, v DataURL) error {
	hook.DB[k] = v
	hook.ListUser[v.UserID] = append(hook.ListUser[v.UserID], v.ShortPath)
	if hook.WriteURL != nil {
		if hook.WriteURL != nil {
			if err := hook.WriteEvent(&v); err != nil {
				log.Error("Try write ", err)
			}

		}
	}
	return nil
}

func (hook StoreDBinFile) Get(k string) (DataURL, error) {
	if _, ok := hook.DB[k]; ok {
		return hook.DB[k], nil
	} else {
		return DataURL{}, errors.New("key not found")
	}
}

func (hook StoreDBinFile) GetAll(k string, dom string) (model.ResponseURLShortALL, error) {
	buff := model.ResponseURLShortALL{}
	if _, ok := hook.ListUser[k]; ok {
		for i := 0; i < len(hook.ListUser[k]); i++ {
			v := hook.ListUser[k][i]
			if _, ok := hook.DB[v]; ok {
				rp := model.ResponseURL{
					ShortURL:    fmt.Sprintf("%s/%s", dom, hook.DB[v].ShortPath),
					OriginalURL: hook.DB[v].URL,
				}
				buff.ListsURL = append(buff.ListsURL, rp)
			}
		}
		return buff, nil
	} else {
		return model.ResponseURLShortALL{}, errors.New("key not found")
	}
}

func (hook StoreDBinFile) Close() {
	hook.WriteURL.Close()
	hook.ReadURL.Close()
}

func NewFileStorage(filename string) (res *StoreDBinFile, err error) {

	res = &StoreDBinFile{
		DB:       make(map[string]DataURL),
		ListUser: make(map[string][]string),
	}

	if filename != "" {
		res.ReadURL, err = tools.NewConsumer(filename)
		if err != nil {
			log.Fatal(err)
		}
		//defer res.ReadURL.Close()

		res.WriteURL, err = tools.NewProducer(filename)
		if err != nil {
			log.Fatal(err)
		}
		//defer res.WriteURL.Close()
		res.RestoreDdBackupURL()

	} else {
		err = errors.New("file path empty")
	}
	return

}
