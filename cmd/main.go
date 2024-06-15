package main

import (
	"github.com/yamauthi/goexpert-load-test/application"
	"github.com/yamauthi/goexpert-load-test/infra/cli"
)

func main() {
	command := cli.NewLoadTestCommand(
		application.NewLoadTest(),
		application.NewLoadTestReportGenerator(),
	)
	command.Execute()
}
