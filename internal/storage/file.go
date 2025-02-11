package storage

import (
	"bufio"
	"encoding/json"
	"github.com/KirillinED/shortener/internal/dto"
	"os"
)

type FileStorage struct {
	file    *os.File
	scanner *bufio.Scanner
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewFileStorage(filepath string) *FileStorage {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}

	return &FileStorage{
		file:    file,
		scanner: bufio.NewScanner(file),
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}
}

func (f *FileStorage) Write(v any) error {
	if err := f.encoder.Encode(v); err != nil {
		return err
	}

	if err := f.file.Sync(); err != nil {
		return err
	}

	return nil
}

func (f *FileStorage) Close() error {
	return f.file.Close()
}

func (f *FileStorage) ReadLink() (*dto.Link, error) {
	link := &dto.Link{}
	err := f.decoder.Decode(link)
	if err != nil {
		return nil, err
	}

	return link, nil
}
