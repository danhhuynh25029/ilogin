package main

import (
	"time"

	"github.com/kataras/iris/v12"
	"github.com/pquerna/otp/totp"
)

func main() {
	app := iris.Default()

	// Ping handler
	app.Get("/ping", func(ctx iris.Context) {
		code, err := totp.GenerateCode("qlt6vmy6svfx4bt4rpmisaiyol6hihca", time.Now().UTC())
		if err != nil {
			panic(err)
		}
		ctx.JSON(code)
	})
	app.Listen(":8080")
}
