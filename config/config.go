package config

import (
    "os"
    "encoding/json"
)

const CONFIG_NAME = "config.json"

type Config struct {
    ServerPort       string
    RemoteSwitchHost string
    RemoteSwitchPort string
    RemoteSwitchUID  string
}

func Load() (*Config, error) {
    file, e := os.Open(CONFIG_NAME)
    if e != nil {
        return nil, e
    }
    defer file.Close()

    config := Config{}
    decoder := json.NewDecoder(file)
    if e = decoder.Decode(&config); e != nil {
        return nil, e
    }
    return &config, nil
}
