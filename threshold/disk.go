package threshold

import "sync"

type DiskLevel int

const (
	DiskNormal DiskLevel = iota
	DiskWarning
	DiskCritical
)

var (
	DiskThresholds = struct {
		DiskWarning  float64
		DiskCritical float64
	}{
		DiskWarning:  80.0,
		DiskCritical: 90.0,
	}

	LastDiskLevels = make(map[string]DiskLevel) // per mountpoint
	DiskMu         sync.Mutex
)
