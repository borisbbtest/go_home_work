package storage

import (
	"bufio"
	"encoding/json"
	"errors"

	"github.com/borisbbtest/go_home_work/internal/tools"
)

type StoreDBinFile struct {
	ReadURL  *tools.Consumer
	WriteURL *tools.Producer
	DB       map[string]DataURL
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
		hook.Put(m.ShortPath, m)
	}

}

func (hook StoreDBinFile) Put(k string, v DataURL) error {
	hook.DB[k] = v
	if hook.WriteURL != nil {
		if hook.WriteURL != nil {
			hook.WriteEvent(&v)
		}
	}
	return nil
}
func (hook StoreDBinFile) Get(k string) (DataURL, error) {
	if _, ok := hook.DB[k]; ok {
		return hook.DB[k], nil
	} else {
		return DataURL{}, errors.New("Key not found")
	}
}

func NewFileStorage(filename string) (res *StoreDBinFile, err error) {

	res = &StoreDBinFile{
		DB: make(map[string]DataURL),
	}
	if filename != "" {
		res.ReadURL, err = tools.NewConsumer(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer res.ReadURL.Close()

		res.WriteURL, err = tools.NewProducer(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer res.WriteURL.Close()
		res.RestoreDdBackupURL()

	} else {
		err = errors.New("File path empty")
	}
	return

}
