package persistence

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

type Persistor struct {
	file *os.File
}

func NewPersistor(filePath string) (*Persistor, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &Persistor{file: file}, nil
}

func (p *Persistor) Close() error {
	return p.file.Close()
}

func (p *Persistor) LoadData(requests *[]time.Time) error {
	if err := json.NewDecoder(p.file).Decode(requests); err != nil {
		if err == io.EOF {
			// If the file is empty, return an empty slice
			log.Print("Info: No data to load, skip loading")

			return nil

		}
		return err
	}

	log.Print("Info: Loaded data from saved file, ", rune(len(*requests)), " request found")
	return nil
}

func (p *Persistor) PersistData(data []time.Time) error {
	if len(data) < 1 {
		return nil
	}

	// start from the beginning of the file
	if _, err := p.file.Seek(0, 0); err != nil {
		return err
	}

	// Truncate the file before writing new data
	if err := p.file.Truncate(0); err != nil {
		return err
	}

	log.Print("Info: Persist data to the file")

	return json.NewEncoder(p.file).Encode(data)
}
