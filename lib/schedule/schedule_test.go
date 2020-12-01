package schedule_test

import (
	"testing"
	"time"

	"github.com/mplewis/prune_backups/lib/period"
	"github.com/mplewis/prune_backups/lib/schedule"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSchedule(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Schedule Suite")
}

var _ = Describe("schedule", func() {
	Describe("Parse", func() {
		It("parses a schedule string properly", func() {
			sched := "24x1h,14x1d,8x7d,12x1m,5x1y"
			periods, err := schedule.Parse(sched)
			Expect(err).To(BeNil())
			Expect(periods).To(Equal([]period.Backup{
				{Count: 24, Duration: 1 * time.Hour},
				{Count: 14, Duration: 24 * time.Hour},
				{Count: 8, Duration: 24 * time.Hour * 7},
				{Count: 12, Duration: 24 * time.Hour * 30},
				{Count: 5, Duration: 24 * time.Hour * 365},
			}))
		})

		It("returns parsing errors", func() {
			sched := "FOOxBARd"
			_, err := schedule.Parse(sched)
			Expect(err.Error()).To(Equal("Could not parse schedule item from FOOxBARd"))
		})
	})
})
