package mailer

import (
	"VPSWatchDog/threshold"
	"fmt"
	"github.com/shirou/gopsutil/v4/net"
	"time"
)

func netLevelToString(level threshold.NetLevel) string {
	switch level {
	case threshold.NetNormal:
		return "Normal"
	case threshold.NetWarning:
		return "Warning"
	case threshold.NetCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

func SendNetworkAlertMail(c net.IOCountersStat, totalMBps float64, rxMBps float64, txMBps float64, level threshold.NetLevel) (string, error) {
	now := time.Now()

	subject := fmt.Sprintf("Network Alert [%s]: %s", c.Name, netLevelToString(level))
	body := fmt.Sprintf("Interface: %s\nTotal Throughput: %.2f MB/s (RX: %.2f, TX: %.2f)\nThresholds: warning=%.2f, critical=%.2f\nTimestamp: %s",
		c.Name, totalMBps, rxMBps, txMBps,
		threshold.NetThresholds.WarningMBps, threshold.NetThresholds.CriticalMBps,
		now.Format(time.RFC1123))

	err := SendMail(subject, body)

	return subject, err
}
