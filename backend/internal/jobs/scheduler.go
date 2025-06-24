package jobs

import (
	"github.com/robfig/cron/v3"
	"log"
)

func StartScheduler() {
	c := cron.New()

	_, err := c.AddFunc("0 0 1-7 1 *", func() {
		// This function will run every day at midnight during the first week of January
		ApplyInterestBatch()
	})

	if err != nil {
		log.Fatalf("Cron Add Job Failed: %s", err)
	}
	c.Start()
}
