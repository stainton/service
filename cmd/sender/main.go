package main

import "service-sender/cmd/sender/app"

func main() {
	app.NewCommand().Execute()
}
