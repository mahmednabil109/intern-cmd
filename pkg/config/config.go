package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type File struct {
	Plugins map[string]string `json:"plugins"`
}

func Load(path string) (*File, error) {
	reader, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", path, err)
	}
	defer reader.Close()
	var config File
	if err := json.NewDecoder(reader).Decode(&config); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("failed to serialize config: %w", err)
	}
	return &config, nil
}

func Save(config *File, path string) error {
	writer, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to open config file %s, %w", path, err)
	}
	if err := json.NewEncoder(writer).Encode(config); err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}
	defer writer.Close()
	return err
}
