package mailer

import (
	"VPSWatchDog/threshold"

	"fmt"
	"github.com/shirou/gopsutil/v4/mem"
)

func memLevelToString(level threshold.MemLevel) string {
	switch level {
	case threshold.MemNormal:
		return "Normal"
	case threshold.MemWarning:
		return "Warning"
	case threshold.MemCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

func SendMemoryAlertMail(vmStat *mem.VirtualMemoryStat, usage float64, level threshold.MemLevel) (string, error) {
	subject := fmt.Sprintf("Memory Alert: %s", memLevelToString(level))
	body := fmt.Sprintf("Memory usage: %.2f%%\nThresholds: warning=%.2f, critical=%.2f\nTotal: %v MB, Used: %v MB, Free: %v MB",
		usage, threshold.MemThresholds.MemWarning, threshold.MemThresholds.MemCritical,
		vmStat.Total/1024/1024, vmStat.Used/1024/1024, vmStat.Free/1024/1024)

	err := SendMail(subject, body)

	return subject, err
}
