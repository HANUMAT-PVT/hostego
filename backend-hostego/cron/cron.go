package cron

import (
	"backend-hostego/controllers"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func InitCronJobs() {
	// Create a new cron with IST timezone
	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Printf("Error loading IST timezone: %v\n", err)
		return
	}
	c := cron.New(cron.WithLocation(ist))

	// Run every Sunday at 02:00 (2 AM) IST
	_, err = c.AddFunc("0 2 * * 0", func() {
		log.Println("Starting weekly wallet withdrawal requests generation...")

		// Call the non-HTTP version directly
		if err := controllers.CreateWalletWithdrawalRequests(); err != nil {
			log.Printf("Error generating wallet withdrawal requests: %v\n", err)
		} else {
			log.Println("Successfully generated wallet withdrawal requests")
		}
	})

	if err != nil {
		log.Printf("Error setting up cron job: %v\n", err)
		return
	}

	c.Start()
	log.Println("Cron jobs initialized - Wallet withdrawals scheduled for Sunday 2 AM IST")
}
