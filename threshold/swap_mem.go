package threshold

import "sync"

type SwapLevel int

const (
	SwapNormal SwapLevel = iota
	SwapWarning
	SwapCritical
)

var (
	SwapThresholds = struct {
		SwapWarning  float64
		SwapCritical float64
	}{
		SwapWarning:  50.0,
		SwapCritical: 80.0,
	}

	LastSwapLevel SwapLevel = SwapNormal
	SwapMu        sync.Mutex
)
