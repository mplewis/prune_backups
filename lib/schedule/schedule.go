package schedule

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mplewis/prune_backups/lib/period"
)

var (
	hour         = time.Hour
	day          = 24 * hour
	month        = 30 * day
	year         = 365 * day
	knownPeriods = map[string]time.Duration{
		"h": hour,
		"d": day,
		"m": month,
		"y": year,
	}
)

// Parse parses a schedule string into a series of backup periods.
func Parse(schedule string) ([]period.Backup, error) {
	items := strings.Split(schedule, ",")
	re, err := regexp.Compile(`((\d+)x(\d+)(h|d|m|y))`)
	if err != nil {
		return nil, err
	}

	backups := []period.Backup{}
	for _, item := range items {
		result := re.FindStringSubmatch(item)
		if result == nil {
			return nil, fmt.Errorf("Could not parse schedule item from %s", item)
		}
		count, err := strconv.Atoi(result[2])
		if err != nil {
			return nil, err
		}
		durVal, err := strconv.Atoi(result[3])
		if err != nil {
			return nil, err
		}
		durUnit := result[4]
		durType := knownPeriods[durUnit]
		duration := durType * time.Duration(durVal)
		backups = append(backups, period.Backup{Count: count, Duration: duration})
	}
	return backups, err
}
