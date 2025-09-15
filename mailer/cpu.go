package mailer

import (
	"VPSWatchDog/threshold"
	"fmt"
)

func cpuLevelToString(level threshold.CPULevel) string {
	switch level {
	case threshold.CPUNormal:
		return "Normal"
	case threshold.CPUWarning:
		return "Warning"
	case threshold.CPUCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

func SendCPUAlertMail(avgUsage float64, level threshold.CPULevel) (string, error) {
	subject := fmt.Sprintf("CPU Alert: %s", cpuLevelToString(level))
	body := fmt.Sprintf("Average CPU usage: %.2f%%\nThresholds: warning=%.2f, critical=%.2f",
		avgUsage, threshold.CPUThresholds.CPUWarning, threshold.CPUThresholds.CPUCritical)

	err := SendMail(subject, body)

	return subject, err
}
