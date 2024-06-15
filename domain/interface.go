package domain

import "github.com/yamauthi/goexpert-load-test/domain/entity"

type LoadTestInterface interface {
	Run(config entity.Config) entity.LoadTestResult
}

type LoadTestReportGeneratorInterface interface {
	Generate(result entity.LoadTestResult)
}
