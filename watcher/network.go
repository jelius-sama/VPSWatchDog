package watcher

import (
	"VPSWatchDog/logger"
	"VPSWatchDog/mailer"
	"VPSWatchDog/threshold"
	"fmt"
	"github.com/shirou/gopsutil/v4/net"
	"strings"
	"time"
)

func StartNetPoller(interval time.Duration) {
	counters, err := net.IOCounters(true)
	if err != nil {
		logger.TimedError("Error getting network counters:", err)
	}

	body := "Network poller started with interval " + interval.String() + "\n\n"
	for _, c := range counters {
		if shouldIgnoreInterface(c.Name) {
			continue
		}
		body += fmt.Sprintf("Interface: %s\nBytesRecv: %v, BytesSent: %v, PacketsRecv: %v, PacketsSent: %v\n\n",
			c.Name, c.BytesRecv, c.BytesSent, c.PacketsRecv, c.PacketsSent)
	}

	err = mailer.SendMail("Poller Initialized: Network", body)
	if err != nil {
		logger.TimedError("Failed to send Network init email")
	}

	go func() {
		for {
			counters, err := net.IOCounters(true)
			if err != nil {
				logger.TimedError("Error getting network counters:", err)
				time.Sleep(interval)
				continue
			}

			for _, c := range counters {
				if shouldIgnoreInterface(c.Name) {
					continue
				}

				threshold.NetMu.Lock()

				prev, ok := threshold.LastNetCounters[c.Name]
				if ok {
					// compute delta
					seconds := interval.Seconds()
					rxMBps := float64(c.BytesRecv-prev.BytesRecv) / (1024 * 1024) / seconds
					txMBps := float64(c.BytesSent-prev.BytesSent) / (1024 * 1024) / seconds
					totalMBps := rxMBps + txMBps

					var level threshold.NetLevel
					if totalMBps >= threshold.NetThresholds.CriticalMBps {
						level = threshold.NetCritical
					} else if totalMBps >= threshold.NetThresholds.WarningMBps {
						level = threshold.NetWarning
					} else {
						level = threshold.NetNormal
					}

					lastLevel, ok := threshold.LastNetLevels[c.Name]
					if !ok {
						lastLevel = threshold.NetNormal
					}

					if level != lastLevel {
						// Level changed → send email
						subject, err := mailer.SendNetworkAlertMail(c, totalMBps, rxMBps, txMBps, level)
						if err != nil {
							logger.TimedError("Failed to send network alert for", c.Name, ":", err)
						} else {
							logger.TimedInfo("Network alert email sent:", subject)
						}

						threshold.LastNetLevels[c.Name] = level
					}
				}

				// update snapshot
				threshold.LastNetCounters[c.Name] = c

				threshold.NetMu.Unlock()
			}

			time.Sleep(interval)
		}
	}()
}

func shouldIgnoreInterface(name string) bool {
	n := strings.ToLower(name)
	if n == "lo" || strings.HasPrefix(n, "docker") ||
		strings.HasPrefix(n, "veth") || strings.HasPrefix(n, "br-") ||
		strings.HasPrefix(n, "kube") {
		return true
	}
	return false
}
