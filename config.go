package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
)

type Config struct {
    Api_key string
}

func ReadConfig(path string) (Config, error) {
    config := Config{}
    config_file, err := os.Open(path)
    if err != nil {
        return config, errors.New(fmt.Sprintf("Unable to open config file '%s': %s\n", path, err))
    }

    json_decoder := json.NewDecoder(config_file)
    err = json_decoder.Decode(&config)
    if err != nil {
        return config, errors.New(fmt.Sprintf("Failed to decode config file '%s': %s\n", path, err))
    }
    return config, nil
}
