package threshold

import "sync"

type MemLevel int

const (
	MemNormal MemLevel = iota
	MemWarning
	MemCritical
)

var (
	MemThresholds = struct {
		MemWarning  float64
		MemCritical float64
	}{
		MemWarning:  75.0,
		MemCritical: 90.0,
	}

	LastMemLevel MemLevel = MemNormal
	MemMu        sync.Mutex
)
