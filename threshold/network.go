package threshold

import (
	"github.com/shirou/gopsutil/v4/net"
	"sync"
)

type NetLevel int

const (
	NetNormal NetLevel = iota
	NetWarning
	NetCritical
)

var (
	NetThresholds = struct {
		WarningMBps  float64
		CriticalMBps float64
	}{
		WarningMBps:  50.0,
		CriticalMBps: 200.0,
	}

	LastNetLevels   = make(map[string]NetLevel)           // per interface
	LastNetCounters = make(map[string]net.IOCountersStat) // to calc deltas
	NetMu           sync.Mutex
)
