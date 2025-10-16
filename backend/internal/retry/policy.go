package retry

import (
	"math"
	"time"
)

const (
	BaseDelay = 1 * time.Second
	MaxDelay  = 5 * time.Minute
)

// CalculateBackoff calculates exponential backoff delay
func CalculateBackoff(retryCount int) time.Duration {
	delay := time.Duration(math.Pow(2, float64(retryCount))) * BaseDelay
	if delay > MaxDelay {
		delay = MaxDelay
	}
	return delay
}

// ShouldRetry determines if a job should be retried
func ShouldRetry(retryCount, maxRetries int) bool {
	return retryCount < maxRetries
}
