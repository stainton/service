package main

import "service-sender/cmd/receiver/app"

func main() {
	cmd := app.NewCommand()
	cmd.Execute()
}
