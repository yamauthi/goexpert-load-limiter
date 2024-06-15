package application

import (
	"net/http"
	"sync"
	"time"

	"github.com/yamauthi/goexpert-load-test/domain/entity"
)

type LoadTest struct {
	statusChan chan int
	TestResult entity.LoadTestResult
}

type loadWorker struct {
	requestsAmount int
	statusChan     chan int
	wg             *sync.WaitGroup
}

func NewLoadTest() *LoadTest {
	return &LoadTest{}
}

func (lt *LoadTest) Run(config entity.Config) entity.LoadTestResult {
	lt.TestResult = entity.LoadTestResult{
		StartedAt:   time.Now(),
		StatusCount: make(map[int]int),
	}

	workerWg := &sync.WaitGroup{}
	counterWg := &sync.WaitGroup{}
	counterWg.Add(1)
	lt.statusChan = make(chan int)
	workerRequests := config.ResquestsAmount / config.ConcurrencyCalls
	distribute := config.ResquestsAmount % config.ConcurrencyCalls

	go lt.startCounter(counterWg)
	for range config.ConcurrencyCalls {
		workerWg.Add(1)
		worker := loadWorker{
			requestsAmount: workerRequests,
			statusChan:     lt.statusChan,
			wg:             workerWg,
		}

		if distribute > 0 {
			worker.requestsAmount++
			distribute--
		}

		go worker.execute(config.Url)
	}

	workerWg.Wait()
	close(lt.statusChan)
	counterWg.Wait()

	lt.TestResult.FinishedAt = time.Now()
	return lt.TestResult
}

func (lt *LoadTest) startCounter(wg *sync.WaitGroup) {
	defer wg.Done()
	for status := range lt.statusChan {
		lt.TestResult.TotalRequests++
		lt.TestResult.StatusCount[status]++
	}
}

func (w *loadWorker) execute(url string) {
	defer w.wg.Done()
	for range w.requestsAmount {
		resp, err := http.Get(url)
		if err != nil {
			w.statusChan <- 0
			continue
		}
		defer resp.Body.Close()

		w.statusChan <- resp.StatusCode
	}
}
