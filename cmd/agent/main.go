package main

import (
	"go-metrics-alerting/internal/apps"
	"go-metrics-alerting/internal/apps/cli"
	"os"
)

func main() {
	command := apps.NewAgentCommand()
	code := cli.Run(command)
	os.Exit(code)

}
