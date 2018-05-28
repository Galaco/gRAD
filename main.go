package main

import (
	"github.com/urfave/cli"
	"os"
	"log"
)

func main() {
	app := cli.NewApp()
	app.Name = "Galaco's RADiant"
	app.Usage = "A radiosity implementation written in golang"

	app.Action = Rad

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
