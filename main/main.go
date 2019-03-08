package main

import (
    "github.com/labstack/echo"
    "RemoteController/socket"
)

func main() {
    router := echo.New()
    router.GET("switchTo", socket.SwitchTo)
    router.Logger.Fatal(router.Start(":8082"))
}
