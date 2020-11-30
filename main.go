package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mplewis/prune_backups/lib/period"
)

func randInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
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
	daily := period.Backup{Count: 4, Duration: day}
	weekly := period.Backup{Count: 4, Duration: week}
	monthly := period.Backup{Count: 4, Duration: month}
	yearly := period.Backup{Count: 4, Duration: year}

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
		times = append([]time.Time{time.Now().Add(time.Duration(i*12) * time.Hour)}, times...)
		intervals := []period.Backup{daily, weekly, monthly, yearly}
		times := period.AllTimesForIntervals(intervals, times)
		for _, t := range times {
			fmt.Printf("%+v\n", t)
		}
	}
}
