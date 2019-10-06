package mathutil

import (
	"math"
	"time"
)

// Percentile finds the relative standing in a slice of floats.
// `percent` should be given on the interval [0,100.0).
func Percentile(input []float64, percent float64) float64 {
	if len(input) == 0 {
		return 0
	}

	return PercentileSorted(CopySort(input), percent)
}

// PercentileSorted finds the relative standing in a sorted slice of floats.
// `percent` should be given on the interval [0,100.0).
func PercentileSorted(sortedInput []float64, percent float64) float64 {
	index := (percent / 100.0) * float64(len(sortedInput))
	percentile := float64(0)
	i := int(math.RoundToEven(index))
	if index == float64(int64(index)) {
		percentile = (sortedInput[i-1] + sortedInput[i]) / 2.0
	} else {
		percentile = sortedInput[i-1]
	}

	return percentile
}

// PercentileOfDuration finds the relative standing in a slice of durations
func PercentileOfDuration(input []time.Duration, percentile float64) time.Duration {
	if len(input) == 0 {
		return 0
	}
	return PercentileSortedDurations(CopySortDurations(input), percentile)
}

// PercentileSortedDurations finds the relative standing in a sorted slice of durations
func PercentileSortedDurations(sortedInput []time.Duration, percentile float64) time.Duration {
	index := (percentile / 100.0) * float64(len(sortedInput))
	if index == float64(int64(index)) {
		i := int(RoundPlaces(index, 0))

		if i < 1 {
			return time.Duration(0)
		}

		return MeanDurations([]time.Duration{sortedInput[i-1], sortedInput[i]})
	}

	i := int(RoundPlaces(index, 0))
	if i < 1 {
		return time.Duration(0)
	}

	return sortedInput[i-1]
}
