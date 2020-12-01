package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mplewis/prune_backups/lib/schedule"
)

func env(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Unset environment variable: %s\n", key)
	}
	return val
}

func main() {
	backups, err := schedule.Parse(env("SCHEDULE"))
	if err != nil {
		log.Fatal(err)
	}
	for _, backup := range backups {
		fmt.Println(backup)
	}

	// sess := session.Must(session.NewSession())
	// client := s3.New(sess)
	// params := &s3.ListObjectsV2Input{
	// 	Bucket: aws.String(env("BUCKET")),
	// }
	// if val, ok := os.LookupEnv("PREFIX"); ok {
	// 	params.Prefix = aws.String(val)
	// }

	// objs := []*s3.Object{}
	// err := client.ListObjectsV2Pages(params,
	// 	func(page *s3.ListObjectsV2Output, lastPage bool) bool {
	// 		objs = append(objs, page.Contents...)
	// 		return true
	// 	})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(len(objs))

	// dates := []time.Time{}
	// for _, obj := range objs {
	// 	dates = append(dates, *obj.LastModified)
	// }

	// fmt.Println(len(dates))

	// backups := []period.Backup{
	// 	{Count: 7, Duration: day},
	// 	{Count: 4, Duration: 7 * day},
	// 	{Count: 4, Duration: 30 * day},
	// 	{Count: 4, Duration: 90 * day},
	// }
	// results := period.AllTimesForIntervals(backups, dates)
	// fmt.Println(len(results))
	// for _, selectedTime := range results {
	// 	for _, cand := range objs {
	// 		if *cand.LastModified == selectedTime {
	// 			fmt.Println(cand)
	// 			continue
	// 		}
	// 	}
	// }
}
