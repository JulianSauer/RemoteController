package remote

import (
    "fmt"
    "strconv"
    "github.com/Tinkerforge/go-api-bindings/remote_switch_v2_bricklet"
    "gopkg.in/resty.v1"
    "time"
)

const INPUT_DELY = 1000

var lastButtonPress = currentTime()
var lastHouseCode uint8 = 255
var lastReceiverCode uint8 = 255
var lastSwitchTo remote_switch_v2_bricklet.SwitchTo = 255

var remote_config *Remote_Config

func ParseCall(houseCode uint8, receiverCode uint8, switchTo remote_switch_v2_bricklet.SwitchTo, repeats uint16) {
    if currentTime()-lastButtonPress < INPUT_DELY &&
        lastHouseCode == houseCode &&
        lastReceiverCode == receiverCode &&
        lastSwitchTo == switchTo {
        return
    } else {
        lastHouseCode = houseCode
        lastReceiverCode = receiverCode
        lastSwitchTo = switchTo
    }

    if remote_config == nil {
        var e error
        remote_config, e = Load()
        if e != nil {
            fmt.Println(e.Error())
            lastButtonPress = currentTime()
            return
        }
    }

    if remote_config.HouseCode != houseCode {
        fmt.Println("Wrong house code")
        fmt.Println("Expected " + strconv.Itoa(int(remote_config.HouseCode)) + " but got " + strconv.Itoa(int(houseCode)))
    }

    for _, button := range remote_config.Buttons {
        if button.Id == receiverCode {
            var response *resty.Response
            var e error
            if switchTo == 0 {
                response, e = resty.R().Get(button.Off)
            } else {
                response, e = resty.R().Get(button.On)
            }

            if e != nil {
                fmt.Println(e.Error())
            }
            fmt.Println("Executed: " + response.Request.URL)
            fmt.Println("Response: " + response.String())
            break
        }
    }

    lastButtonPress = currentTime()
}

func currentTime() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}
