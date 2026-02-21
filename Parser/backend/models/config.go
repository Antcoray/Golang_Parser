package models

import (
	"encoding/json"
	"os"
)

type Config struct {
	Operators map[string]bool `json:"operators"`
}

func LoadConfig(fpath string) Config {

	file, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}
