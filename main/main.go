package main

import (
    "github.com/labstack/echo"
    "RemoteController/socket"
    "RemoteController/config"
)

func main() {
    router := echo.New()
    router.GET("switchTo", socket.SwitchTo)

    configuration, e := config.Load()
    if e != nil {
        router.Logger.Fatal(e.Error())
    }

    router.Logger.Fatal(router.Start(":" + configuration.ServerPort))
}
