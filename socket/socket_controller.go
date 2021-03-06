package socket

import (
    "github.com/JulianSauer/RemoteController/config"
    "github.com/Tinkerforge/go-api-bindings/ipconnection"
    "github.com/Tinkerforge/go-api-bindings/remote_switch_v2_bricklet"
    "github.com/labstack/echo"
    "net/http"
    "strconv"
    "time"
)

var configuration *config.Config

// Rest interface for switching sockets
func SwitchTo(context echo.Context) error {
    if e := checkParameter("houseCode", context); e != nil {
        return e
    }
    if e := checkParameter("receiverCode", context); e != nil {
        return e
    }
    if e := checkParameter("switchTo", context); e != nil {
        return e
    }

    var houseCode uint64
    var receiverCode uint64
    var switchTo bool
    var timeout int
    var e error
    if houseCode, e = strconv.ParseUint(context.QueryParam("houseCode"), 10, 8); e != nil {
        return context.String(http.StatusBadGateway, e.Error())
    }
    if receiverCode, e = strconv.ParseUint(context.QueryParam("receiverCode"), 10, 8); e != nil {
        return context.String(http.StatusBadGateway, e.Error())
    }
    if switchTo, e = strconv.ParseBool(context.QueryParam("switchTo")); e != nil {
        return context.String(http.StatusBadGateway, e.Error())
    }
    if context.QueryParam("timeout") == "" {
        timeout = -1
    } else if timeout, e = strconv.Atoi(context.QueryParam("timeout")); e != nil {
        return context.String(http.StatusBadGateway, e.Error())
    }

    if e = switchSocketTo(uint8(houseCode), uint8(receiverCode), switchTo); e != nil {
        return context.String(http.StatusBadGateway, e.Error())
    }
    if timeout > 0 {
        switchSocketIn(timeout, uint8(houseCode), uint8(receiverCode), !switchTo, context.Logger())
    }
    return nil
}

// Returns a 400 if the parameter is missing
func checkParameter(parameterKey string, context echo.Context) error {
    if context.QueryParam(parameterKey) == "" {
        errorMessage := "Missing parameter: " + parameterKey
        context.Logger().Warn(errorMessage)
        return context.String(http.StatusBadRequest, errorMessage)
    }
    return nil
}

// Switches a socket of type A using the Tinkerforge Remote Switch Bricklet 2.0
func switchSocketTo(houseCode uint8, receiverCode uint8, switchTo bool) error {
    if configuration == nil {
        var e error
        configuration, e = config.Load()
        if e != nil {
            return e
        }
    }

    connection := ipconnection.New()
    defer connection.Close()

    remoteSwitch, e := remote_switch_v2_bricklet.New(configuration.RemoteSwitchUID, &connection)
    if e != nil {
        return e
    }

    connection.Connect(configuration.RemoteSwitchHost + ":" + configuration.RemoteSwitchPort)
    defer connection.Disconnect()

    if switchTo {
        return remoteSwitch.SwitchSocketA(houseCode, receiverCode, remote_switch_v2_bricklet.SwitchToOn)
    } else {
        return remoteSwitch.SwitchSocketA(houseCode, receiverCode, remote_switch_v2_bricklet.SwitchToOff)
    }
}

// Switches a socket after a given time in minutes
func switchSocketIn(timeout int, houseCode uint8, receiverCode uint8, switchTo bool, logger echo.Logger) {
    go func() {
        time.Sleep(time.Duration(timeout) * time.Minute)
        if e := switchSocketTo(houseCode, receiverCode, switchTo); e != nil {
            logger.Warn("Error during scheduled switching: " + e.Error())
        }
    }()
}
