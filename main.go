package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type backupPeriod struct {
	Count    int
	Duration time.Duration
}

func randInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func sortDesc(times []time.Time) {
	sort.Slice(times, func(i, j int) bool {
		return times[i].After(times[j])
	})
}

func earliestTimeWithinInterval(interval time.Duration, baseTime time.Time, times []time.Time) time.Time {
	result := time.Time{}
	limit := baseTime.Add(-1 * interval)
	for _, cand := range times {
		if result.IsZero() {
			result = cand
			continue
		}
		if cand.Before(limit) {
			continue
		}
		if cand.After(result) {
			continue
		}
		result = cand
	}
	return result
}

func earlierThan(base time.Time, times []time.Time) []time.Time {
	result := []time.Time{}
	for _, cand := range times {
		if cand.Before(base) {
			result = append(result, cand)
		}
	}
	return result
}

func followingTimesByInterval(period backupPeriod, keepInitial bool, times []time.Time) []time.Time {
	chain := []time.Time{}
	if keepInitial {
		chain = append(chain, times[0])
	}
	for {
		if len(chain) >= period.Count || len(times) == 0 {
			return chain
		}
		earliest := earliestTimeWithinInterval(period.Duration, times[0], times)
		chain = append(chain, earliest)
		times = earlierThan(earliest, times)
	}
}

func allTimesForIntervals(periods []backupPeriod, times []time.Time) []time.Time {
	remaining := times
	sortDesc(remaining)
	fmt.Printf("first: %+v\n", remaining[0])
	allTimes := []time.Time{}
	lastTime := time.Time{}
	for _, period := range periods {
		first := lastTime.IsZero()
		results := followingTimesByInterval(period, first, remaining)
		allTimes = append(allTimes, results...)
		// fmt.Printf("%+v\n", period)
		// fmt.Printf("%+v\n", len(results))
		if len(results) == 0 {
			return allTimes
		}
		lastTime = results[len(results)-1]
		remaining = earlierThan(lastTime, remaining)
	}
	return allTimes
}

func main() {
	rand.Seed(time.Now().UnixNano())
	last := time.Now()
	times := []time.Time{last}
	for i := 0; i < 10000; i++ {
		count := time.Duration(randInt(-48, -4))
		last = last.Add(time.Hour * count)
		times = append(times, last)
	}
	fmt.Printf("first: %+v\n", times[0])

	day := time.Hour * 24
	week := day * 7
	month := day * 30
	year := day * 365
	daily := backupPeriod{Count: 4, Duration: day}
	weekly := backupPeriod{Count: 4, Duration: week}
	monthly := backupPeriod{Count: 4, Duration: month}
	yearly := backupPeriod{Count: 4, Duration: year}

	// dailyResults := followingTimesByInterval(daily, true, times)
	// for _, t := range dailyResults {
	// 	fmt.Printf("%+v\n", t)
	// }
	// fmt.Println()
	// lastDay := dailyResults[len(dailyResults)-1]
	// remaining := earlierThan(lastDay, times)
	// weeklyResults := followingTimesByInterval(weekly, false, remaining)
	// for _, t := range weeklyResults {
	// 	fmt.Printf("%+v\n", t)
	// }

	// intervals := []backupPeriod{daily, weekly, monthly, yearly}
	// allResults := allTimesForIntervals(intervals, times)
	// for _, t := range allResults {
	// 	fmt.Printf("%+v\n", t)
	// }

	for i := 0; i < 1; i++ {
		duration := time.Duration(i*12) * time.Hour
		forward := time.Now().Add(duration)
		fmt.Printf("forward: %+v\n", forward)
		fmt.Printf("times: %+v\n", times[0])
		times = append([]time.Time{forward}, times...)
		fmt.Printf("times: %+v\n", times[0])
		fmt.Printf("times: %+v\n", times[1])
		sortDesc(times)
		fmt.Printf("desc: %+v\n", times[0])
		fmt.Printf("desc: %+v\n", times[len(times)-1])
		fmt.Println()
		intervals := []backupPeriod{daily, weekly, monthly, yearly}
		times := allTimesForIntervals(intervals, times)
		for _, t := range times {
			fmt.Printf("%+v\n", t)
		}
	}
}
