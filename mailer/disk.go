package mailer

import (
	"VPSWatchDog/threshold"
	"fmt"
	"github.com/shirou/gopsutil/v4/disk"
)

func diskLevelToString(level threshold.DiskLevel) string {
	switch level {
	case threshold.DiskNormal:
		return "Normal"
	case threshold.DiskWarning:
		return "Warning"
	case threshold.DiskCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

func SendDiskAlertMail(p disk.PartitionStat, level threshold.DiskLevel, usage *disk.UsageStat) (string, error) {
	subject := fmt.Sprintf("Disk Alert [%s]: %s", p.Mountpoint, diskLevelToString(level))
	body := fmt.Sprintf("Mount: %s\nFilesystem: %s\nUsage: %.2f%%\nThresholds: warning=%.2f, critical=%.2f\nTotal: %v GB, Used: %v GB, Free: %v GB",
		p.Mountpoint, p.Fstype,
		usage.UsedPercent, threshold.DiskThresholds.DiskWarning, threshold.DiskThresholds.DiskCritical,
		usage.Total/1024/1024/1024, usage.Used/1024/1024/1024, usage.Free/1024/1024/1024)

	err := SendMail(subject, body)
	return subject, err
}
