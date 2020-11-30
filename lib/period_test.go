package period_test

import (
	"log"
	"testing"
	"time"

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

func parseYyyyMmDd(s string) time.Time {
	result, err := time.Parse("2006-01-02", s)
	if err != nil {
		log.Panic(err)
	}
	return result
}

func days(from string, to string, interval time.Duration) []time.Time {
	fromD := parseYyyyMmDd(from)
	toD := parseYyyyMmDd(to)
	dates := []time.Time{}
	for {
		if fromD.After(toD) {
			return dates
		}

	}
}

var _ = Describe("period", func() {
	Describe("SortDesc", func() {
		It("sorts in descending order", func() {
			dates := days("2020-05-13", "2020-05-20", day)
			period.SortDesc(dates)
		})
	})
})
