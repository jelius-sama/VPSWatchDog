package mailer

import (
	"VPSWatchDog/threshold"
	"fmt"
	"github.com/shirou/gopsutil/v4/mem"
)

func swapLevelToString(level threshold.SwapLevel) string {
	switch level {
	case threshold.SwapNormal:
		return "Normal"
	case threshold.SwapWarning:
		return "Warning"
	case threshold.SwapCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

func SendSwapMemAlertMail(swapStat *mem.SwapMemoryStat, usage float64, level threshold.SwapLevel) (string, error) {
	subject := fmt.Sprintf("Swap Alert: %s", swapLevelToString(level))
	body := fmt.Sprintf("Swap usage: %.2f%%\nThresholds: warning=%.2f, critical=%.2f\nTotal: %v MB, Used: %v MB, Free: %v MB",
		usage, threshold.SwapThresholds.SwapWarning, threshold.SwapThresholds.SwapCritical,
		swapStat.Total/1024/1024, swapStat.Used/1024/1024, swapStat.Free/1024/1024)

	err := SendMail(subject, body)

	return subject, err
}
