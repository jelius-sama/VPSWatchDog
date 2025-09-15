package watcher

import (
	"VPSWatchDog/logger"
	"VPSWatchDog/mailer"
	"VPSWatchDog/threshold"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/mem"
)

func StartMemPoller(interval time.Duration) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		logger.TimedError("Error reading memory stats:", err)
	}

	err = mailer.SendMail("Poller Initialized: Memory", fmt.Sprintf("Memory poller started with interval %s\nTotal: %v MB\nUsed: %v MB\nFree: %v MB\nUsedPercent: %.2f%%", interval, vmStat.Total/1024/1024, vmStat.Used/1024/1024, vmStat.Free/1024/1024, vmStat.UsedPercent))
	if err != nil {
		logger.TimedError("Failed to send Memory init email")
	}

	go func() {
		for {
			vmStat, err := mem.VirtualMemory()
			if err != nil {
				logger.TimedError("Error reading memory stats:", err)
				time.Sleep(interval) // avoid tight loop if error
				continue
			}

			usage := vmStat.UsedPercent
			var level threshold.MemLevel
			if usage >= threshold.MemThresholds.MemCritical {
				level = threshold.MemCritical
			} else if usage >= threshold.MemThresholds.MemWarning {
				level = threshold.MemWarning
			} else {
				level = threshold.MemNormal
			}

			threshold.MemMu.Lock()
			if level != threshold.LastMemLevel {
				// Level changed → send email
				subject, err := mailer.SendMemoryAlertMail(vmStat, usage, level)

				if err != nil {
					logger.TimedError("Failed to send memory alert email:", err)
				} else {
					logger.TimedInfo("Memory alert email sent:", subject)
				}

				threshold.LastMemLevel = level
			}
			threshold.MemMu.Unlock()

			time.Sleep(interval)
		}
	}()
}
