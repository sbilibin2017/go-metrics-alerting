package main

import (
	"go-metrics-alerting/internal/apps"
	"go-metrics-alerting/internal/apps/cli"
	"os"
)

func main() {
	command := apps.NewServerCommand()
	code := cli.Run(command)
	os.Exit(code)
}
