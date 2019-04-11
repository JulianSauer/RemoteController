package main

import (
    "github.com/labstack/echo"
    "github.com/Juliansauer/RemoteController/socket"
    "github.com/Juliansauer/RemoteController/config"
    "github.com/Tinkerforge/go-api-bindings/ipconnection"
    "github.com/Tinkerforge/go-api-bindings/remote_switch_v2_bricklet"
    "github.com/Juliansauer/RemoteController/remote"
)

func main() {
    router := echo.New()
    router.GET("switchTo", socket.SwitchTo)

    configuration, e := config.Load()
    if e != nil {
        router.Logger.Fatal(e.Error())
    }

    connection := ipconnection.New()
    defer connection.Close()
    remoteSwitch, e := remote_switch_v2_bricklet.New(configuration.RemoteSwitchUID, &connection)
    if e != nil {
        router.Logger.Fatal(e.Error())
    }
    connection.Connect(configuration.RemoteSwitchHost + ":" + configuration.RemoteSwitchPort)
    defer connection.Disconnect()
    remoteSwitch.RegisterRemoteStatusACallback(remote.ParseCall)

    router.Logger.Fatal(router.Start(":" + configuration.ServerPort))
}
