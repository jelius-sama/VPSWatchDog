package intelligence

import (
	"sync"
	"time"
)

// AlertLevel represents the severity of an alert
type AlertLevel int

const (
	AlertNone AlertLevel = iota
	AlertMinor
	AlertModerate
	AlertSevere
)

// AlertState tracks the alerting state for a specific monitor
type AlertState struct {
	// Rate limiting
	MaxAlertsPerHour int
	AlertCount       int
	WindowStart      time.Time

	// Monitoring state
	IsEnabled      bool
	LastAlertTime  time.Time
	LastAlertLevel AlertLevel
	LastDeviation  float64

	// For progressive alerting
	DeviationHistory []float64 // Recent deviations to track trends
	MaxHistorySize   int

	mu sync.Mutex
}

// NewAlertState creates a new alert state
func NewAlertState(maxAlertsPerHour int) *AlertState {
	return &AlertState{
		MaxAlertsPerHour: maxAlertsPerHour,
		WindowStart:      time.Now(),
		IsEnabled:        true,
		MaxHistorySize:   20, // Keep last 20 measurements
		DeviationHistory: make([]float64, 0, 20),
	}
}

// RecordDeviation adds a deviation measurement to history
func (a *AlertState) RecordDeviation(deviation float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.DeviationHistory = append(a.DeviationHistory, deviation)
	if len(a.DeviationHistory) > a.MaxHistorySize {
		// Remove oldest
		a.DeviationHistory = a.DeviationHistory[1:]
	}
}

// ShouldAlert determines if an alert should be sent
// Returns (shouldAlert, reason)
func (a *AlertState) ShouldAlert(deviation float64, minProgressiveGap float64) (bool, string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if monitoring is enabled
	if !a.IsEnabled {
		return false, "monitoring disabled"
	}

	// Check if we're in a new hour window
	now := time.Now()
	if now.Sub(a.WindowStart) >= time.Hour {
		// Reset for new window
		a.resetWindow(now)
	}

	// Check if we've hit the alert limit
	if a.AlertCount >= a.MaxAlertsPerHour {
		// Disable monitoring until next window
		a.IsEnabled = false
		return false, "alert limit reached for this hour"
	}

	// First alert in the window - send if anomalous
	if a.AlertCount == 0 {
		return true, "first alert in window"
	}

	// For subsequent alerts, check if deviation is significantly worse
	deviationIncrease := deviation - a.LastDeviation
	if deviationIncrease < minProgressiveGap {
		return false, "deviation not significantly worse than previous alert"
	}

	return true, "progressive alert - situation worsening"
}

// RecordAlert records that an alert was sent
func (a *AlertState) RecordAlert(deviation float64, level AlertLevel) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.AlertCount++
	a.LastAlertTime = time.Now()
	a.LastAlertLevel = level
	a.LastDeviation = deviation

	// If we hit the limit, disable monitoring
	if a.AlertCount >= a.MaxAlertsPerHour {
		a.IsEnabled = false
	}
}

// resetWindow resets the hourly window
func (a *AlertState) resetWindow(now time.Time) {
	a.AlertCount = 0
	a.WindowStart = now
	a.IsEnabled = true // Re-enable monitoring if it was disabled
	a.LastDeviation = 0
}

// CheckAndResetWindow checks if an hour has passed and resets if needed
// Call this periodically in your monitoring loop
func (a *AlertState) CheckAndResetWindow() bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now()
	if now.Sub(a.WindowStart) >= time.Hour {
		wasDisabled := !a.IsEnabled
		a.resetWindow(now)
		return wasDisabled // Returns true if monitoring was re-enabled
	}

	return false
}

// GetStatus returns current state information
func (a *AlertState) GetStatus() (alertCount int, isEnabled bool, timeUntilReset time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()

	timeUntilReset = time.Hour - time.Since(a.WindowStart)
	if timeUntilReset < 0 {
		timeUntilReset = 0
	}

	return a.AlertCount, a.IsEnabled, timeUntilReset
}

// GetAverageRecentDeviation calculates average deviation from recent history
func (a *AlertState) GetAverageRecentDeviation(samples int) float64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(a.DeviationHistory) == 0 {
		return 0
	}

	// Get last N samples
	startIdx := len(a.DeviationHistory) - samples
	if startIdx < 0 {
		startIdx = 0
	}

	sum := 0.0
	count := 0
	for i := startIdx; i < len(a.DeviationHistory); i++ {
		sum += a.DeviationHistory[i]
		count++
	}

	if count == 0 {
		return 0
	}

	return sum / float64(count)
}

// CalculateAlertLevel determines severity based on deviation
func CalculateAlertLevel(deviation float64) AlertLevel {
	absDeviation := deviation
	if absDeviation < 0 {
		absDeviation = -absDeviation
	}

	if absDeviation >= 50.0 {
		return AlertSevere
	} else if absDeviation >= 35.0 {
		return AlertModerate
	} else if absDeviation >= 20.0 {
		return AlertMinor
	}

	return AlertNone
}

// AlertLevelToString converts alert level to string
func AlertLevelToString(level AlertLevel) string {
	switch level {
	case AlertNone:
		return "Normal"
	case AlertMinor:
		return "Minor"
	case AlertModerate:
		return "Moderate"
	case AlertSevere:
		return "Severe"
	default:
		return "Unknown"
	}
}
