package main

import (
	"github.com/ebcp-dev/gorest-api/app"
)

func main() {
	a := app.App{}

	a.Initialize()
	a.Run(":8010")
}
