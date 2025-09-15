package watcher

import (
	"VPSWatchDog/logger"
	"VPSWatchDog/mailer"
	"VPSWatchDog/threshold"
	"fmt"
	"github.com/shirou/gopsutil/v4/disk"
	"strings"
	"time"
)

func StartDiskPoller(interval time.Duration) {
	parts, err := disk.Partitions(false)
	if err != nil {
		logger.TimedError("Error getting partitions:", err)
	}

	body := "Disk poller started with interval " + interval.String() + "\n\n"
	for _, p := range parts {
		if shouldIgnorePartition(p) {
			continue
		}

		u, err := disk.Usage(p.Mountpoint)
		if err != nil {
			logger.TimedError("Error reading usage for", p.Mountpoint, ":", err)
			continue
		}

		body += fmt.Sprintf("[%s] %s\nTotal: %v GB, Used: %v GB, Free: %v GB (%.2f%%)\n\n",
			p.Mountpoint, p.Fstype,
			u.Total/1024/1024/1024, u.Used/1024/1024/1024, u.Free/1024/1024/1024, u.UsedPercent)
	}

	err = mailer.SendMail("Poller Initialized: Disk", body)
	if err != nil {
		logger.TimedError("Failed to send Disk init email")
	}

	go func() {
		for {
			partitions, err := disk.Partitions(false)
			if err != nil {
				logger.TimedError("Error getting partitions:", err)
				time.Sleep(interval)
				continue
			}

			for _, p := range partitions {
				if shouldIgnorePartition(p) {
					continue
				}

				usage, err := disk.Usage(p.Mountpoint)
				if err != nil {
					logger.TimedError("Error reading usage for", p.Mountpoint, ":", err)
					continue
				}

				var level threshold.DiskLevel

				if usage.UsedPercent >= threshold.DiskThresholds.DiskCritical {
					level = threshold.DiskCritical
				} else if usage.UsedPercent >= threshold.DiskThresholds.DiskWarning {
					level = threshold.DiskWarning
				} else {
					level = threshold.DiskNormal
				}

				threshold.DiskMu.Lock()
				lastLevel, ok := threshold.LastDiskLevels[p.Mountpoint]
				if !ok {
					lastLevel = threshold.DiskNormal
				}

				if level != lastLevel {
					// Level changed → send email
					subject, err := mailer.SendDiskAlertMail(p, level, usage)

					if err != nil {
						logger.TimedError("Failed to send disk alert for", p.Mountpoint, ":", err)
					} else {
						logger.TimedInfo("Disk alert email sent:", subject)
					}

					threshold.LastDiskLevels[p.Mountpoint] = level
				}
				threshold.DiskMu.Unlock()
			}

			time.Sleep(interval)
		}
	}()
}

// skip system/virtual filesystems
func shouldIgnorePartition(p disk.PartitionStat) bool {
	fstype := strings.ToLower(p.Fstype)
	if fstype == "tmpfs" || fstype == "devtmpfs" || fstype == "proc" ||
		fstype == "sysfs" || fstype == "squashfs" || fstype == "overlay" ||
		fstype == "efivarfs" || fstype == "autofs" {
		return true
	}

	// also ignore snap mounts, docker overlay, etc.
	if strings.HasPrefix(p.Mountpoint, "/snap") ||
		strings.HasPrefix(p.Mountpoint, "/var/lib/docker") {
		return true
	}

	return false
}
