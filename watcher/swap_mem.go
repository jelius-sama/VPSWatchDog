package watcher

import (
	"VPSWatchDog/logger"
	"VPSWatchDog/mailer"
	"VPSWatchDog/threshold"
	"fmt"
	"github.com/shirou/gopsutil/v4/mem"
	"time"
)

func StartSwapPoller(interval time.Duration) {
	swapStat, err := mem.SwapMemory()
	if err != nil {
		logger.TimedError("Error reading swap stats:", err)
	}

	err = mailer.SendMail("Poller Initialized: Swap", fmt.Sprintf("Swap poller started with interval %s\nTotal: %v MB\nUsed: %v MB\nFree: %v MB\nUsedPercent: %.2f%%", interval, swapStat.Total/1024/1024, swapStat.Used/1024/1024, swapStat.Free/1024/1024, swapStat.UsedPercent))
	if err != nil {
		logger.TimedError("Failed to send Swap Memory init email")
	}

	go func() {
		for {
			swapStat, err := mem.SwapMemory()
			if err != nil {
				logger.TimedError("Error reading swap stats:", err)
				time.Sleep(interval)
				continue
			}

			usage := swapStat.UsedPercent
			var level threshold.SwapLevel

			if usage >= threshold.SwapThresholds.SwapCritical {
				level = threshold.SwapCritical
			} else if usage >= threshold.SwapThresholds.SwapWarning {
				level = threshold.SwapWarning
			} else {
				level = threshold.SwapNormal
			}

			threshold.SwapMu.Lock()
			if level != threshold.LastSwapLevel {
				subject, err := mailer.SendSwapMemAlertMail(swapStat, usage, level)
				if err != nil {
					logger.TimedError("Failed to send swap alert:", err)
				} else {
					logger.TimedInfo("Swap alert email sent:", subject)
				}

				threshold.LastSwapLevel = level
			}
			threshold.SwapMu.Unlock()

			time.Sleep(interval)
		}
	}()
}
