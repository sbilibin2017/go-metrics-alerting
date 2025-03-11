package main

import (
	"go-metrics-alerting/internal/apps"
	"go-metrics-alerting/internal/engines/cli"

	"os"
)

func main() {
	command := apps.NewAgentAppCommand()
	code := cli.RunCommand(command)
	os.Exit(code)
}
