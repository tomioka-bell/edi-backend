package jobs

import (
	portsvc "backend/internal/core/ports/services"
	"backend/internal/pkgs/mailer"

	"context"
	"log"
	"time"

	cron "github.com/robfig/cron/v3"
)

func StartOrderPeriodAlertCron(svc portsvc.EDISummaryDataService) *cron.Cron {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatalf("load location error: %v", err)
	}

	c := cron.New(cron.WithLocation(loc))

	_, err = c.AddFunc("24 10 * * *", func() {
		ctx := context.Background()

		log.Println("[CRON] Running GetOrderPeriodAlerts...")

		alerts, err := svc.GetOrderPeriodAlerts(ctx)
		if err != nil {
			log.Println("[CRON] GetOrderPeriodAlerts error:", err)
			return
		}

		for _, a := range alerts {
			// log ไว้ debug
			log.Printf("[ALERT] FC=%s Vendor=%s Target=%s Emails=%v\n",
				a.NumberOrder,
				a.VendorCode,
				a.TargetTime.Format(time.RFC3339),
				a.Emails,
			)

			// ส่งเมลหาแต่ละ email
			for _, email := range a.Emails {
				if email == "" {
					continue
				}

				if err := mailer.SendForecastReadReminder(
					email,
					a.NumberOrder,
					a.VendorCode,
					a.PeriodFrom,
					a.TargetTime,
				); err != nil {
					log.Printf("[MAIL] send forecast reminder failed: %v (to=%s, fc=%s)\n",
						err, email, a.NumberOrder)
				} else {
					log.Printf("[MAIL] forecast reminder sent to %s (fc=%s)\n", email, a.NumberOrder)
				}
			}
		}
	})

	if err != nil {
		log.Fatalf("cron add func error: %v", err)
	}

	c.Start()
	return c
}
