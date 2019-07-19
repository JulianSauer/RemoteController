package remote

import (
    "fmt"
    "github.com/Tinkerforge/go-api-bindings/remote_switch_v2_bricklet"
    "gopkg.in/resty.v1"
    "net/url"
    "os/exec"
    "strconv"
    "time"
)

var remoteConfig *Remote_Config

const INPUT_DELAY = 1000

var lastHouseCode uint8
var lastReceiverCode uint8
var lastSwitchTo remote_switch_v2_bricklet.SwitchTo
var lastCurrentTime int64

func ParseCall(houseCode uint8, receiverCode uint8, switchTo remote_switch_v2_bricklet.SwitchTo, repeats uint16) {
    if currentTime()-lastCurrentTime < INPUT_DELAY &&
        lastHouseCode == houseCode &&
        lastReceiverCode == receiverCode &&
        lastSwitchTo == switchTo {
        return
    }

    lastHouseCode = houseCode
    lastReceiverCode = receiverCode
    lastSwitchTo = switchTo
    lastCurrentTime = currentTime()

    go func() {
        if remoteConfig == nil {
            var e error
            remoteConfig, e = Load()
            if e != nil {
                fmt.Println(e.Error())
                return
            }
        }

        if remoteConfig.HouseCode != houseCode {
            fmt.Println("Wrong house code")
            fmt.Println("Expected " + strconv.Itoa(int(remoteConfig.HouseCode)) + " but got " + strconv.Itoa(int(houseCode)))
        }

        for _, button := range remoteConfig.Buttons {
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
    }()
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

func currentTime() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}
