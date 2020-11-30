package period

import (
	"fmt"
	"sort"
	"time"
)

// BackupPeriod represents a number of periodic backups to maintain.
type BackupPeriod struct {
	Count    int
	Duration time.Duration
}

// SortDesc sorts a list of Times in descending order in-place.
func SortDesc(times []time.Time) {
	sort.Slice(times, func(i, j int) bool {
		return times[i].After(times[j])
	})
}

// EarliestTimeWithinInterval picks the earliest time within the interval (base - interval ... base).
// Times must be in descending order.
func EarliestTimeWithinInterval(interval time.Duration, base time.Time, times []time.Time) time.Time {
	result := time.Time{}
	limit := base.Add(-1 * interval)
	for _, cand := range times {
		if result.IsZero() {
			result = cand
			continue
		}
		if cand.Before(limit) {
			// continue
			return result
		}
		if cand.After(result) {
			continue
		}
		result = cand
	}
	return result
}

// EarlierThan filters a list of Times to all that precede base.
func EarlierThan(base time.Time, times []time.Time) []time.Time {
	result := []time.Time{}
	for _, cand := range times {
		if cand.Before(base) {
			result = append(result, cand)
		}
	}
	return result
}

// FollowingTimesByInterval finds a number of Times, moving backward from the first Time in the descending-sorted Time
// slice, keeping at most period.duration time between them, and keeping at most period.count times.
// If keepInitial is true, the first time in the result will be the first time in times. Otherwise, the first time is
// omitted.
// Times must be in descending order.
func FollowingTimesByInterval(period BackupPeriod, keepInitial bool, times []time.Time) []time.Time {
	chain := []time.Time{}
	if keepInitial {
		chain = append(chain, times[0])
	}
	for {
		if len(chain) >= period.Count || len(times) == 0 {
			return chain
		}
		earliest := EarliestTimeWithinInterval(period.Duration, times[0], times)
		chain = append(chain, earliest)
		times = EarlierThan(earliest, times)
	}
}

// AllTimesForIntervals evaluates a long list of Times and filters it to the ones that match the periods' counts and
// durations, in the order the periods were listed.
func AllTimesForIntervals(periods []BackupPeriod, times []time.Time) []time.Time {
	remaining := times
	SortDesc(remaining)
	fmt.Printf("first: %+v\n", remaining[0])
	allTimes := []time.Time{}
	lastTime := time.Time{}
	for _, period := range periods {
		first := lastTime.IsZero()
		results := FollowingTimesByInterval(period, first, remaining)
		allTimes = append(allTimes, results...)
		// fmt.Printf("%+v\n", period)
		// fmt.Printf("%+v\n", len(results))
		if len(results) == 0 {
			return allTimes
		}
		lastTime = results[len(results)-1]
		remaining = EarlierThan(lastTime, remaining)
	}
	return allTimes
}
