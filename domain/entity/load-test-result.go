package entity

import "time"

type LoadTestResult struct {
	StartedAt     time.Time
	FinishedAt    time.Time
	StatusCount   map[int]int
	TotalRequests int
}
