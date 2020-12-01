package period_test

import (
	"log"
	"testing"
	"time"

	"github.com/mplewis/prune_backups/lib/period"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	day = 24 * time.Hour
)

func TestPeriod(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Period Suite")
}

func ymd(s string) time.Time {
	result, err := time.Parse("2006-01-02", s)
	if err != nil {
		log.Panic(err)
	}
	return result
}

func days(from string, to string) []time.Time {
	fromD := ymd(from)
	toD := ymd(to)
	dates := []time.Time{}
	for {
		if fromD.After(toD) {
			return dates
		}
		dates = append(dates, fromD)
		fromD = fromD.Add(day)
	}
}

var _ = Describe("period", func() {
	// Describe("SortDesc", func() {
	// 	It("sorts dates in descending order", func() {
	// 		dates := days("2020-05-10", "2020-05-14")
	// 		period.SortDesc(dates)
	// 		Expect(dates[0]).To(Equal(ymd("2020-05-14")))
	// 		Expect(dates[4]).To(Equal(ymd("2020-05-10")))
	// 	})
	// })

	Describe("EarliestTimeWithinInterval", func() {
		It("picks the proper time for a descending-order date range", func() {
			dates := days("2020-01-01", "2020-03-01")
			period.SortDesc(dates)
			base := ymd("2020-02-15")
			Expect(period.EarliestTimeWithinInterval(7*day, base, dates)).To(Equal(ymd("2020-02-08")))
			Expect(period.EarliestTimeWithinInterval(14*day, base, dates)).To(Equal(ymd("2020-02-01")))
		})
	})

	Describe("EarlierThan", func() {
		It("filters a list of dates properly", func() {
			dates := days("2020-01-01", "2020-03-01")
			Expect(period.EarlierThan(ymd("2020-01-08"), dates)).To(Equal([]time.Time{
				ymd("2020-01-01"),
				ymd("2020-01-02"),
				ymd("2020-01-03"),
				ymd("2020-01-04"),
				ymd("2020-01-05"),
				ymd("2020-01-06"),
				ymd("2020-01-07"),
			}))
		})
	})

	Describe("EarlierThanOrEqual", func() {
		It("filters a list of dates properly", func() {
			dates := days("2020-01-01", "2020-03-01")
			Expect(period.EarlierThanOrEqual(ymd("2020-01-08"), dates)).To(Equal([]time.Time{
				ymd("2020-01-01"),
				ymd("2020-01-02"),
				ymd("2020-01-03"),
				ymd("2020-01-04"),
				ymd("2020-01-05"),
				ymd("2020-01-06"),
				ymd("2020-01-07"),
				ymd("2020-01-08"),
			}))
		})
	})

	Describe("FollowingTimesByInterval", func() {
		It("picks periodic times for an interval", func() {
			dates := days("2020-01-01", "2020-03-01")
			period.SortDesc(dates)
			Expect(period.FollowingTimesByInterval(
				period.Backup{Count: 4, Duration: 7 * day}, true, dates,
			)).To(Equal([]time.Time{
				ymd("2020-03-01"),
				ymd("2020-02-23"),
				ymd("2020-02-16"),
				ymd("2020-02-09"),
			}))
			Expect(period.FollowingTimesByInterval(
				period.Backup{Count: 4, Duration: 7 * day}, false, dates,
			)).To(Equal([]time.Time{
				ymd("2020-02-23"),
				ymd("2020-02-16"),
				ymd("2020-02-09"),
				ymd("2020-02-02"),
			}))
			Expect(period.FollowingTimesByInterval(
				period.Backup{Count: 10, Duration: 30 * day}, true, dates,
			)).To(Equal([]time.Time{
				ymd("2020-03-01"),
				ymd("2020-01-31"),
				ymd("2020-01-01"),
			}))
		})

		It("handles the case when no more dates are found within the interval", func() {
			dates := []time.Time{
				ymd("2020-01-28"),
				ymd("2020-01-24"),
				ymd("2020-01-01"),
			}
			Expect(period.FollowingTimesByInterval(
				period.Backup{Count: 4, Duration: 7 * day}, true, dates,
			)).To(Equal([]time.Time{
				ymd("2020-01-28"),
				ymd("2020-01-24"),
			}))
		})
	})

	Describe("AllTimesForIntervals", func() {
		It("picks periodic times for all intervals", func() {
			// dates := days("2020-09-01", "2020-12-31")
			// periods := []period.Backup{{Count: 3, Duration: day}}
			dates := days("2016-01-01", "2020-12-31")
			fiveDaily := period.Backup{Count: 5, Duration: day}
			fourWeekly := period.Backup{Count: 4, Duration: day * 7}
			fiveMonthly := period.Backup{Count: 5, Duration: day * 30}
			threeYearly := period.Backup{Count: 3, Duration: day * 365}
			periods := []period.Backup{fiveDaily, fourWeekly, fiveMonthly, threeYearly}
			Expect(period.AllTimesForIntervals(periods, dates)).To(Equal([]time.Time{
				// daily
				ymd("2020-12-31"),
				ymd("2020-12-30"),
				ymd("2020-12-29"),
				ymd("2020-12-28"),
				ymd("2020-12-27"),
				// weekly
				ymd("2020-12-20"),
				ymd("2020-12-13"),
				ymd("2020-12-06"),
				ymd("2020-11-29"),
				// monthly
				ymd("2020-10-30"),
				ymd("2020-09-30"),
				ymd("2020-08-31"),
				ymd("2020-08-01"),
				ymd("2020-07-02"),
				// yearly
				ymd("2019-07-03"),
				ymd("2018-07-03"),
				ymd("2017-07-03"),
			}))
		})
	})
})
