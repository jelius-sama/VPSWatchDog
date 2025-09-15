package watcher

import (
	"VPSWatchDog/logger"
	"VPSWatchDog/mailer"
	"VPSWatchDog/threshold"
	"fmt"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/load"
	"time"
)

func StartLoadPoller(interval time.Duration) {
	avg, err := load.Avg()
	if err != nil {
		logger.TimedError("Error reading load averages:", err)
	}
	cores, err := cpu.Counts(true)
	if err != nil {
		logger.TimedError("Error getting CPU cores:", err)
		cores = 1
	}

	err = mailer.SendMail("Poller Initialized: Load", fmt.Sprintf("Load poller started with interval %s\nCores: %d\nLoad averages: 1m=%.2f, 5m=%.2f, 15m=%.2f", interval, cores, avg.Load1, avg.Load5, avg.Load15))
	if err != nil {
		logger.TimedError("Failed to send Load init email")
	}

	go func() {
		// get number of logical cores for scaling thresholds
		cores, err := cpu.Counts(true)
		if err != nil {
			logger.TimedError("Error getting CPU cores:", err)
			cores = 1
		}

		for {
			avg, err := load.Avg()
			if err != nil {
				logger.TimedError("Error reading load averages:", err)
				time.Sleep(interval)
				continue
			}

			// compare 1-min load average to thresholds * number of cores
			load1 := avg.Load1
			var level threshold.LoadLevel

			if load1 >= threshold.LoadThresholds.CriticalPerCore*float64(cores) {
				level = threshold.LoadCritical
			} else if load1 >= threshold.LoadThresholds.WarningPerCore*float64(cores) {
				level = threshold.LoadWarning
			} else {
				level = threshold.LoadNormal
			}

			threshold.LoadMu.Lock()
			if level != threshold.LastLoadLevel {
				subject, err := mailer.SendLoadAlertMail(avg, cores, level)

				if err != nil {
					logger.TimedError("Failed to send load alert:", err)
				} else {
					logger.TimedInfo("Load alert email sent:", subject)
				}

				threshold.LastLoadLevel = level
			}
			threshold.LoadMu.Unlock()

			time.Sleep(interval)
		}
	}()
}
