package main

import (
	"AuthInGo/app"
)

func main() {

	cfg := app.NewConfig()
	app := app.NewApplication(cfg)

	app.Run()

}
