package brennenstuhl

import (
    "RemoteController/remote"
    "os"
    "encoding/json"
)

const CONFIG_NAME = "brennenstuhl.json"

// Config for a brennenstuhl RCS 1000 N Comfort remote
type config_brennenstuhl struct {
    HouseCode uint8
    A remote.Button
    B remote.Button
    C remote.Button
    D remote.Button
}

func Load() (*config_brennenstuhl, error) {
    file, e := os.Open(CONFIG_NAME)
    if e != nil {
        return nil, e
    }
    defer file.Close()

    config := config_brennenstuhl{}
    decoder := json.NewDecoder(file)
    if e = decoder.Decode(&config); e != nil {
        return nil, e
    }
    return &config, nil
}
