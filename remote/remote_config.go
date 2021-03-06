package remote

import (
    "os"
    "encoding/json"
    "strings"
)

const CONFIG_NAME = "remote.json"

// Config for a brennenstuhl RCS 1000 N Comfort remote
type Remote_Config struct {
    HouseCode uint8
    Buttons   []Button
}

func Load() (*Remote_Config, error) {
    file, e := os.Open(CONFIG_NAME)
    if e != nil {
        return nil, e
    }
    defer file.Close()

    config := Remote_Config{}
    decoder := json.NewDecoder(file)
    if e = decoder.Decode(&config); e != nil {
        return nil, e
    }
    for _, button := range config.Buttons {
        button.On = strings.Replace(button.On, "\\", "", -1)
        button.Off = strings.Replace(button.Off, "\\", "", -1)
    }
    return &config, nil
}
