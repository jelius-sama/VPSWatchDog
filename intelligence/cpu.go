package intelligence

import (
	"VPSWatchDog/baseline"
	"fmt"
	"sync"
	"time"
)

// CPUIntelligence combines baseline data with alerting logic
type CPUIntelligence struct {
	Baseline   *baseline.CPUBaseline
	AlertState *AlertState
	Storage    *baseline.Storage

	// Adaptive learning parameters
	AdaptiveWindowSize int // Number of samples to consider for baseline updates
	RecentUsages       []float64

	mu sync.RWMutex
}

// NewCPUIntelligence creates a new CPU intelligence engine
func NewCPUIntelligence(
	baselinePath string,
	learningDuration time.Duration,
	maxAlertsPerHour int,
) *CPUIntelligence {
	storage := baseline.NewStorage(baselinePath)

	// Try to load existing baseline
	existingBaseline, err := storage.LoadCPUBaseline()

	var cpuBaseline *baseline.CPUBaseline
	if err == nil && existingBaseline != nil && !existingBaseline.IsLearning {
		// Use existing baseline
		cpuBaseline = existingBaseline
	} else {
		// Create new baseline in learning mode
		cpuBaseline = baseline.NewCPUBaseline(learningDuration)
	}

	return &CPUIntelligence{
		Baseline:           cpuBaseline,
		AlertState:         NewAlertState(maxAlertsPerHour),
		Storage:            storage,
		AdaptiveWindowSize: 12, // 12 samples for adaptive learning (e.g., 1 hour if sampling every 5 min)
		RecentUsages:       make([]float64, 0, 12),
	}
}

// ProcessUsage processes a CPU usage reading
// Returns (shouldAlert, alertMessage, error)
func (c *CPUIntelligence) ProcessUsage(usage float64) (bool, string, error) {
	c.mu.Lock()

	// Check if we need to reset the alert window
	wasReEnabled := c.AlertState.CheckAndResetWindow()
	if wasReEnabled {
		// Log that monitoring was re-enabled
		c.mu.Unlock()
		return false, "", nil
	}

	// If in learning mode, just record the sample
	if c.Baseline.IsLearning {
		c.Baseline.RecordSample(usage)

		// Check if learning period is over
		elapsed := time.Since(c.Baseline.LearningStarted)
		if elapsed >= c.Baseline.LearningDuration {
			c.Baseline.FinalizeLearning()

			// Save baseline to disk
			if err := c.Storage.SaveCPUBaseline(c.Baseline); err != nil {
				c.mu.Unlock()
				return false, "", err
			}

			c.mu.Unlock()
			return true, "Learning phase completed. Intelligent monitoring now active.", nil
		}

		c.mu.Unlock()
		return false, "", nil
	}

	// Monitoring mode
	isAnomalous, deviation, description := c.Baseline.IsAnomalous(usage)

	// Record deviation for trend analysis
	c.AlertState.RecordDeviation(deviation)
	c.RecentUsages = append(c.RecentUsages, usage)
	if len(c.RecentUsages) > c.AdaptiveWindowSize {
		c.RecentUsages = c.RecentUsages[1:]
	}

	// Check for adaptive baseline update (sustained moderate changes)
	if len(c.RecentUsages) >= c.AdaptiveWindowSize {
		recentDeviations := make([]float64, 0, c.AdaptiveWindowSize)
		hour := time.Now().Hour()

		for _, u := range c.RecentUsages {
			dev := c.Baseline.CalculateDeviation(u, hour)
			recentDeviations = append(recentDeviations, dev)
		}

		if c.Baseline.ShouldUpdateBaseline(recentDeviations, c.AdaptiveWindowSize) {
			// Calculate new average for this hour
			sum := 0.0
			for _, u := range c.RecentUsages {
				sum += u
			}
			newAvg := sum / float64(len(c.RecentUsages))

			c.Baseline.UpdateHourlyBaseline(hour, newAvg)

			// Save updated baseline
			if err := c.Storage.SaveCPUBaseline(c.Baseline); err != nil {
				// Log error but don't fail
			}
		}
	}

	c.mu.Unlock()

	// If not anomalous, no alert needed
	if !isAnomalous {
		return false, "", nil
	}

	// Check if we should send an alert
	shouldAlert, reason := c.AlertState.ShouldAlert(
		deviation,
		c.Baseline.MinAlertGap,
	)

	if !shouldAlert {
		return false, "", nil
	}

	// Determine alert level
	level := CalculateAlertLevel(deviation)

	// Record the alert
	c.AlertState.RecordAlert(deviation, level)

	// Build alert message
	hour := time.Now().Hour()
	expected := c.Baseline.GetExpectedUsage(hour)
	isPeak := c.Baseline.IsPeakHour(hour)

	peakStr := "off-peak"
	if isPeak {
		peakStr = "peak"
	}

	message := fmt.Sprintf(
		"CPU Usage Anomaly Detected\n"+
			"----------------------------\n"+
			"%s Alert\n\n"+
			"%s\n\n"+
			"Details:\n"+
			"  Current Usage: %.2f%%\n"+
			"  Expected Usage: %.2f%%\n"+
			"  Deviation: %.2f%%\n"+
			"  Hour: %d:00 (%s)\n"+
			"  Alert Count: %d/%d this hour\n\n"+
			"Reason: %s",
		AlertLevelToString(level),
		description,
		usage,
		expected,
		deviation,
		hour,
		peakStr,
		c.AlertState.AlertCount,
		c.AlertState.MaxAlertsPerHour,
		reason,
	)

	return true, message, nil
}

// GetLearningProgress returns learning progress information
func (c *CPUIntelligence) GetLearningProgress() (isLearning bool, progress float64, remaining time.Duration) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.Baseline.IsLearning {
		return false, 100.0, 0
	}

	elapsed := time.Since(c.Baseline.LearningStarted)
	remaining = c.Baseline.LearningDuration - elapsed

	if remaining < 0 {
		remaining = 0
	}

	progress = (float64(elapsed) / float64(c.Baseline.LearningDuration)) * 100.0

	return true, progress, remaining
}

// GetBaselineSummary returns a summary of the current baseline
func (c *CPUIntelligence) GetBaselineSummary() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.Baseline.GetSummary()
}

// Helper functions for formatting
func formatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func formatInt(i int) string {
	return fmt.Sprintf("%d", i)
}
