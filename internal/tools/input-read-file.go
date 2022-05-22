package tools

import (
	"encoding/json"
	"os"
)

type Consumer struct {
	File    *os.File
	Decoder *json.Decoder
}

type Producer struct {
	File    *os.File
	Encoder *json.Encoder
}

func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &Producer{
		File:    file,
		Encoder: json.NewEncoder(file),
	}, nil
}
func (p *Producer) GetFile() *os.File {
	return p.File
}

func (p *Producer) Close() error {
	return p.File.Close()
}

func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		File:    file,
		Decoder: json.NewDecoder(file),
	}, nil
}
func (c *Consumer) GetFile() *os.File {
	return c.File
}

func (c *Consumer) Close() error {
	return c.File.Close()
}
