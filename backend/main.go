package main

import (
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"os"

	"github.com/kataras/iris/v12"
	"github.com/pquerna/otp/totp"
)

var users = make(map[string]string)

type UserRequest struct {
	Username string `json:"username"`
	OTP      string `json:"otp"`
}

func main() {
	app := iris.Default()
	app.Get("/generate/{name}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		code, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "APP",
			AccountName: name,
		})

		if _, ok := users[name]; !ok {
			fmt.Println("----", code.Secret())
			users[name] = code.Secret()
		}

		fmt.Println(code.URL())

		var buf bytes.Buffer
		img, err := code.Image(200, 200)
		if err != nil {
			panic(err)
		}
		png.Encode(&buf, img)

		// display the QR code to the user.
		os.WriteFile("qr-code.png", buf.Bytes(), 0644)

		if err != nil {
			panic(err)
		}
		ctx.JSON(code)
	})
	app.Post("/validate", func(ctx iris.Context) {
		var user UserRequest
		if err := ctx.ReadBody(&user); err != nil {
			ctx.StopWithError(iris.StatusBadRequest, err)
			return
		}
		isValid := totp.Validate(user.OTP, users[user.Username])
		if !isValid {
			ctx.StopWithError(iris.StatusBadRequest, errors.New("OTP is not valid"))
			return
		}
		ctx.StopWithJSON(iris.StatusOK, nil)
	})
	app.Listen(":8080")
}
