package application

import (
	"fmt"
	"net/http"

	"github.com/yamauthi/goexpert-load-test/domain/entity"
)

type LoadTestReportGenerator struct{}

func NewLoadTestReportGenerator() *LoadTestReportGenerator {
	return &LoadTestReportGenerator{}
}

func (rg *LoadTestReportGenerator) Generate(result entity.LoadTestResult) {
	fmt.Println("----------------Load Test Result----------------")
	fmt.Println("-- Test duration: ", result.FinishedAt.Sub(result.StartedAt))
	fmt.Println("-- Total requests: ", result.TotalRequests)
	fmt.Println("-- Status OK(200) requests: ", result.StatusCount[http.StatusOK])

	for status, value := range result.StatusCount {
		if status != http.StatusOK && value > 0 {
			statusTxt := http.StatusText(status)
			fmt.Printf("-- Status %s(%v) requests: %v\n", statusTxt, status, value)
		}
	}
	fmt.Println("------------------------------------------------")
}
