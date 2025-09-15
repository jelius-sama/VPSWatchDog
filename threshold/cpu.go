package threshold

import "sync"

type CPULevel int

const (
	CPUNormal CPULevel = iota
	CPUWarning
	CPUCritical
)

var (
	CPUThresholds = struct {
		CPUWarning  float64
		CPUCritical float64
	}{
		CPUWarning:  70.0,
		CPUCritical: 90.0,
	}

	LastCPULevel CPULevel = CPUNormal
	CPUMu        sync.Mutex
)
