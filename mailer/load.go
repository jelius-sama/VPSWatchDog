package mailer

import (
	"VPSWatchDog/threshold"
	"fmt"
	"github.com/shirou/gopsutil/v4/load"
)

func loadLevelToString(level threshold.LoadLevel) string {
	switch level {
	case threshold.LoadNormal:
		return "Normal"
	case threshold.LoadWarning:
		return "Warning"
	case threshold.LoadCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

func SendLoadAlertMail(avg *load.AvgStat, cores int, level threshold.LoadLevel) (string, error) {
	subject := fmt.Sprintf("Load Alert: %s", loadLevelToString(level))
	body := fmt.Sprintf("Load averages: 1m=%.2f, 5m=%.2f, 15m=%.2f\nCores=%d\nThresholds per core: warning=%.2f, critical=%.2f",
		avg.Load1, avg.Load5, avg.Load15,
		cores, threshold.LoadThresholds.WarningPerCore, threshold.LoadThresholds.CriticalPerCore)

	err := SendMail(subject, body)

	return subject, err
}
