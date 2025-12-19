package jobs

import (
	portsvc "backend/internal/core/ports/services"
	"backend/internal/pkgs/mailer"

	"context"
	"log"
	"time"

	cron "github.com/robfig/cron/v3"
)

func StartForecastPeriodAlertCron(svc portsvc.EDISummaryDataService) *cron.Cron {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatalf("load location error: %v", err)
	}

	c := cron.New(cron.WithLocation(loc))

	_, err = c.AddFunc("24 10 * * *", func() {
		ctx := context.Background()

		log.Println("[CRON] Running GetForecastPeriodAlerts...")

		alerts, err := svc.GetForecastPeriodAlerts(ctx)
		if err != nil {
			log.Println("[CRON] GetForecastPeriodAlerts error:", err)
			return
		}

		for _, a := range alerts {
			// log ไว้ debug
			log.Printf("[ALERT] FC=%s Vendor=%s Target=%s Emails=%v\n",
				a.NumberForecast,
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
					a.NumberForecast,
					a.VendorCode,
					a.PeriodFrom,
					a.TargetTime,
				); err != nil {
					log.Printf("[MAIL] send forecast reminder failed: %v (to=%s, fc=%s)\n",
						err, email, a.NumberForecast)
				} else {
					log.Printf("[MAIL] forecast reminder sent to %s (fc=%s)\n", email, a.NumberForecast)
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
