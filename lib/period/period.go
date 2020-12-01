package period

import (
	"fmt"
	"sort"
	"time"
)

// Backup represents a number of periodic backups to maintain.
type Backup struct {
	Count    int
	Duration time.Duration
}

func (b Backup) String() string {
	return fmt.Sprintf("{%d backups over %+v}", b.Count, b.Duration)
}

// SortDesc sorts a list of Times in descending order in-place.
func SortDesc(times []time.Time) {
	sort.Slice(times, func(i, j int) bool {
		return times[i].After(times[j])
	})
}

// EarliestTimeWithinInterval picks the earliest time within the interval (base - interval ... base) inclusive.
func EarliestTimeWithinInterval(interval time.Duration, base time.Time, times []time.Time) time.Time {
	result := time.Time{}
	limit := base.Add(-1 * interval)
	for _, cand := range times {
		if cand.Before(limit) || cand.After(base) {
			continue
		}
		result = cand
	}
	return result
}

// EarlierThan filters a list of Times to all times that precede base.
func EarlierThan(base time.Time, times []time.Time) []time.Time {
	result := []time.Time{}
	for _, cand := range times {
		if cand.Before(base) {
			result = append(result, cand)
		}
	}
	return result
}

// EarlierThanOrEqual filters a list of Times to all times that are equal to or precede base.
func EarlierThanOrEqual(base time.Time, times []time.Time) []time.Time {
	result := []time.Time{}
	for _, cand := range times {
		if cand.Before(base) || cand.Equal(base) {
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
func FollowingTimesByInterval(period Backup, keepInitial bool, times []time.Time) []time.Time {
	chain := []time.Time{}
	if keepInitial {
		chain = append(chain, times[0])
	}
	base := times[0]
	for {
		if len(chain) >= period.Count || len(times) == 0 {
			return chain
		}
		earliest := EarliestTimeWithinInterval(period.Duration, base, times)
		if earliest.IsZero() {
			return chain
		}
		base = earliest
		chain = append(chain, earliest)
		times = EarlierThan(earliest, times)
	}
}

// AllTimesForIntervals evaluates a long list of Times and filters it to the ones that match the periods' counts and
// durations, in the order the periods were listed.
// This method modifies the original times slice to sort the times in descending order.
func AllTimesForIntervals(periods []Backup, times []time.Time) []time.Time {
	remaining := times
	SortDesc(remaining)
	allTimes := []time.Time{}
	lastTime := time.Time{}
	for _, period := range periods {
		results := FollowingTimesByInterval(period, lastTime.IsZero(), remaining)
		allTimes = append(allTimes, results...)
		if len(results) == 0 {
			return allTimes
		}
		lastTime = results[len(results)-1]
		remaining = EarlierThanOrEqual(lastTime, remaining)
	}
	return allTimes
}
