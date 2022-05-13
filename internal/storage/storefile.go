package storage

import (
	"encoding/json"
	"os"
)

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

type InitStoreDBinFile struct {
	ReadURL  *consumer
	WriteURL *producer
}

type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewProducer(fileName string) (*producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}
func (p *producer) GetFile() *os.File {
	return p.file
}
func (p *producer) WriteEvent(event *StorageURL) error {
	return p.encoder.Encode(&event)
}
func (p *producer) Close() error {
	return p.file.Close()
}

func NewConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}
func (c *consumer) GetFile() *os.File {
	return c.file
}
func (c *consumer) ReadEvent() (*StorageURL, error) {
	event := &StorageURL{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}
func (c *consumer) Close() error {
	return c.file.Close()
}
