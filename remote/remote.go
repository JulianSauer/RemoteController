package remote

import (
    "fmt"
    "strconv"
    "github.com/Tinkerforge/go-api-bindings/remote_switch_v2_bricklet"
    "gopkg.in/resty.v1"
    "time"
    "net/url"
    "os/exec"
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
            var call string
            if switchTo == 0 {
                call = button.Off
            } else {
                call = button.On
            }
            fmt.Println("Executing: " + call)
            var response string
            if isValidUrl(call) {
                response = executeRestCall(call)
            } else {
                response = executeShellCommand(call)
            }
            fmt.Println("Response: " + response)
            break
        }
    }

    lastButtonPress = currentTime()
}

func currentTime() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}

func isValidUrl(input string) bool {
    u, e := url.Parse(input)
    return e == nil && u.Scheme != "" && u.Host != ""
}

func executeRestCall(url string) string {
    response, e := resty.R().Get(url)

    if e != nil {
        fmt.Println(e.Error())
    }
    return response.String()
}

func executeShellCommand(cmd string) string {
    response, e := exec.Command(cmd).Output()
    if e != nil {
        fmt.Println(e.Error())
    }
    return fmt.Sprintf("%s", response)
}
