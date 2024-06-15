package application_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yamauthi/goexpert-load-test/application"
	"github.com/yamauthi/goexpert-load-test/domain"
	"github.com/yamauthi/goexpert-load-test/domain/entity"
)

const DefaultStatus = http.StatusOK

type statusRequests struct {
	StatusCode int
	Amount     int
	Current    atomic.Int32
}

type statusCodeReturn struct {
	Distribution []*statusRequests
	Current      int
}

func (s *statusCodeReturn) getStatusCode() int {
	statusResult := DefaultStatus

	if len(s.Distribution) > 0 {
		if s.Current < len(s.Distribution) &&
			s.Distribution[s.Current].Current.Load() < int32(s.Distribution[s.Current].Amount) {

			s.Distribution[s.Current].Current.Add(1)
			statusResult = s.Distribution[s.Current].StatusCode
		} else {
			s.Current++
			return s.getStatusCode()
		}
	}

	return statusResult
}

type LoadTestTestSuite struct {
	suite.Suite
	LoadTestInterface domain.LoadTestInterface
	MockHttpServer    *httptest.Server
	StatusCodeReturn  *statusCodeReturn
}

func (suite *LoadTestTestSuite) SetupTest() {
	suite.MockHttpServer = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(suite.StatusCodeReturn.getStatusCode())
			},
		),
	)

	suite.LoadTestInterface = application.NewLoadTest()
}

func (suite *LoadTestTestSuite) TearDownTest() {
	suite.MockHttpServer.Close()
}

func TestLoadTestTestSuite(t *testing.T) {
	suite.Run(t, new(LoadTestTestSuite))
}

func (suite *LoadTestTestSuite) TestLoadTest_Run() {
	type input struct {
		ConcurrencyCalls int
		TotalRequests    int
		Distribution     []*statusRequests
	}

	type output struct {
		StatusCount   map[int]int
		TotalRequests int
	}

	type TestCase struct {
		TestName string
		Input    input
		Expected output
	}

	testsCase := []TestCase{
		{
			TestName: "Test with TotalRequest parameter NOT divisor of ConcurrencyCalls",
			Input: input{
				TotalRequests:    100,
				ConcurrencyCalls: 13,
			},
			Expected: output{
				StatusCount: map[int]int{
					http.StatusOK:                  50,
					http.StatusBadRequest:          30,
					http.StatusInternalServerError: 20,
				},
				TotalRequests: 100,
			},
		},
		{
			TestName: "Test with TotalRequest parameter divisor of ConcurrencyCalls",
			Input: input{
				TotalRequests:    1000,
				ConcurrencyCalls: 20,
			},
			Expected: output{
				StatusCount: map[int]int{
					http.StatusOK:                  750,
					http.StatusBadRequest:          75,
					http.StatusInternalServerError: 75,
					http.StatusForbidden:           100,
				},
				TotalRequests: 1000,
			},
		},
	}

	for _, tc := range testsCase {
		statusCodeReturn := &statusCodeReturn{}
		for s, v := range tc.Expected.StatusCount {
			statusCodeReturn.Distribution = append(statusCodeReturn.Distribution, &statusRequests{
				StatusCode: s,
				Amount:     v,
			})
		}
		suite.StatusCodeReturn = statusCodeReturn

		suite.Run(tc.TestName, func() {
			config := entity.Config{
				Url:              suite.MockHttpServer.URL,
				ResquestsAmount:  tc.Input.TotalRequests,
				ConcurrencyCalls: tc.Input.ConcurrencyCalls,
			}

			testResult := suite.LoadTestInterface.Run(config)

			suite.NotEmpty(testResult.StartedAt)
			suite.NotEmpty(testResult.FinishedAt)
			suite.Equal(tc.Expected.TotalRequests, testResult.TotalRequests)
			suite.Equal(len(tc.Expected.StatusCount), len(testResult.StatusCount))
			for status, expected := range tc.Expected.StatusCount {
				suite.Equal(expected, testResult.StatusCount[status])
			}
		})
	}
}
