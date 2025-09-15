package watcher

import (
	"time"

	"VPSWatchDog/logger"
	"VPSWatchDog/mailer"
	"VPSWatchDog/threshold"
	"fmt"
	"github.com/shirou/gopsutil/v4/cpu"
)

func StartCPUPoller(interval time.Duration) {
	percentages, err := cpu.Percent(200*time.Millisecond, false)
	if err != nil || len(percentages) == 0 {
		logger.TimedError("Error reading CPU usage:", err)
	}

	usage := 0.0
	if len(percentages) > 0 {
		usage = percentages[0]
	}
	cores, err := cpu.Counts(true)
	if err != nil {
		logger.TimedError("Error reading CPU core counts:", err)
	}

	err = mailer.SendMail("Poller Initialized: CPU", fmt.Sprintf("CPU poller started with interval %s\nCores: %d\nCurrent Avg Usage: %.2f%%", interval, cores, usage))
	if err != nil {
		logger.TimedError("Failed to send CPU init email")
	}

	go func() {
		for {
			percentages, err := cpu.Percent(interval, false) // avg across all cores
			if err != nil || len(percentages) == 0 {
				logger.TimedError("Error reading CPU usage:", err)
				continue
			}

			avgUsage := percentages[0]
			var level threshold.CPULevel

			if avgUsage >= threshold.CPUThresholds.CPUCritical {
				level = threshold.CPUCritical
			} else if avgUsage >= threshold.CPUThresholds.CPUWarning {
				level = threshold.CPUWarning
			} else {
				level = threshold.CPUNormal
			}

			threshold.CPUMu.Lock()
			if level != threshold.LastCPULevel {
				// Level changed → send email
				subject, err := mailer.SendCPUAlertMail(avgUsage, level)

				if err != nil {
					logger.TimedError("Failed to send CPU alert email:", err)
				} else {
					logger.TimedInfo("CPU alert email sent:", subject)
				}

				threshold.LastCPULevel = level
			}
			threshold.CPUMu.Unlock()
		}
	}()
}
