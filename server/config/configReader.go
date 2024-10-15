package config

import (
	"encoding/json"
	"os"
)

type Conf struct {
	JwtSalt []byte `json:"users"`
}

func ConfFromFile(targetFile string) (Conf, error) {
	var conf Conf
	bytes, err := os.ReadFile(targetFile)
	if err != nil {
		return conf, err
	}

	err := json.Unmarshal(bytes, &conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}
