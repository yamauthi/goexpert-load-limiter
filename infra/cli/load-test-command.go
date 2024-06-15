package cli

import (
	"flag"

	"github.com/yamauthi/goexpert-load-test/domain"
	"github.com/yamauthi/goexpert-load-test/domain/entity"
)

type LoadTestCommand struct {
	loadTest        domain.LoadTestInterface
	reportGenerator domain.LoadTestReportGeneratorInterface
}

func NewLoadTestCommand(
	loadTest domain.LoadTestInterface,
	reportGenerator domain.LoadTestReportGeneratorInterface) *LoadTestCommand {
	return &LoadTestCommand{
		loadTest:        loadTest,
		reportGenerator: reportGenerator,
	}
}

func (c *LoadTestCommand) Execute() {
	var config entity.Config
	flag.StringVar(&config.Url, "url", "", "Endpoint URL to be tested")
	flag.IntVar(&config.ResquestsAmount, "requests", 0, "Total amount of requests")
	flag.IntVar(&config.ConcurrencyCalls, "concurrency", 1, "Number of concurrent endpoint calls")
	flag.Parse()

	c.reportGenerator.Generate(
		c.loadTest.Run(config),
	)
}
