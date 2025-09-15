package threshold

import "sync"

type LoadLevel int

const (
	LoadNormal LoadLevel = iota
	LoadWarning
	LoadCritical
)

var (
	LoadThresholds = struct {
		WarningPerCore  float64
		CriticalPerCore float64
	}{
		WarningPerCore:  0.7,
		CriticalPerCore: 1.0,
	}

	LastLoadLevel LoadLevel = LoadNormal
	LoadMu        sync.Mutex
)
