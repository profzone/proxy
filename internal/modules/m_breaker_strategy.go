package modules

import "github.com/sony/gobreaker"

func BreakerStrategyConsecutiveFailures(counts gobreaker.Counts) bool {
	return counts.ConsecutiveFailures > 10
}

func BreakerStrategyTotalFailures(counts gobreaker.Counts) bool {
	failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
	return counts.Requests >= 3 && failureRatio >= 0.5
}
