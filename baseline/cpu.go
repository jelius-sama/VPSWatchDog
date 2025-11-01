package baseline

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"
)

// CPUBaseline stores learned CPU usage patterns
type CPUBaseline struct {
	// Learning phase
	IsLearning       bool          `json:"is_learning"`
	LearningStarted  time.Time     `json:"learning_started"`
	LearningDuration time.Duration `json:"learning_duration"`

	// Hourly patterns (0-23 hours)
	HourlyAvg     [24]float64 `json:"hourly_avg"`     // Average usage per hour
	HourlyMax     [24]float64 `json:"hourly_max"`     // Max usage seen per hour
	HourlySamples [24]int     `json:"hourly_samples"` // Number of samples per hour

	// Global statistics
	GlobalAvg  float64 `json:"global_avg"`   // Overall average usage
	GlobalMax  float64 `json:"global_max"`   // Maximum usage ever seen
	OffPeakAvg float64 `json:"off_peak_avg"` // Average during off-peak hours
	PeakHours  []int   `json:"peak_hours"`   // Hours considered peak (e.g., [18, 19, 20])

	// Adaptive thresholds (calculated from patterns)
	DeviationThreshold float64 `json:"deviation_threshold"` // % deviation to trigger alert
	MinAlertGap        float64 `json:"min_alert_gap"`       // Minimum % increase for 2nd alert

	// Metadata
	LastUpdated  time.Time `json:"last_updated"`
	TotalSamples int       `json:"total_samples"`
	DataVersion  int       `json:"data_version"` // For future compatibility

	mu sync.RWMutex `json:"-"`
}

// NewCPUBaseline creates a new baseline with default values
func NewCPUBaseline(learningDuration time.Duration) *CPUBaseline {
	return &CPUBaseline{
		IsLearning:         true,
		LearningStarted:    time.Now(),
		LearningDuration:   learningDuration,
		DeviationThreshold: 25.0, // Default: alert if 25% above expected
		MinAlertGap:        15.0, // Default: 2nd alert needs 15% more deviation
		DataVersion:        1,
		PeakHours:          []int{},
	}
}

// RecordSample adds a CPU usage sample during learning phase
func (b *CPUBaseline) RecordSample(usage float64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.IsLearning {
		return
	}

	now := time.Now()
	hour := now.Hour()

	// Update hourly stats
	oldAvg := b.HourlyAvg[hour]
	n := float64(b.HourlySamples[hour])
	b.HourlyAvg[hour] = (oldAvg*n + usage) / (n + 1)
	b.HourlySamples[hour]++

	if usage > b.HourlyMax[hour] {
		b.HourlyMax[hour] = usage
	}

	// Update global stats
	oldGlobalAvg := b.GlobalAvg
	totalN := float64(b.TotalSamples)
	b.GlobalAvg = (oldGlobalAvg*totalN + usage) / (totalN + 1)
	b.TotalSamples++

	if usage > b.GlobalMax {
		b.GlobalMax = usage
	}

	b.LastUpdated = now
}

// FinalizeLearning calculates final statistics and identifies peak hours
func (b *CPUBaseline) FinalizeLearning() {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Identify peak hours (hours with avg > global avg + 15%)
	peakThreshold := b.GlobalAvg * 1.15
	b.PeakHours = []int{}
	offPeakSum := 0.0
	offPeakCount := 0

	for hour := 0; hour < 24; hour++ {
		if b.HourlySamples[hour] == 0 {
			continue
		}

		if b.HourlyAvg[hour] > peakThreshold {
			b.PeakHours = append(b.PeakHours, hour)
		} else {
			offPeakSum += b.HourlyAvg[hour]
			offPeakCount++
		}
	}

	if offPeakCount > 0 {
		b.OffPeakAvg = offPeakSum / float64(offPeakCount)
	} else {
		b.OffPeakAvg = b.GlobalAvg
	}

	b.IsLearning = false
	b.LastUpdated = time.Now()
}

// IsPeakHour checks if given hour is a peak hour
func (b *CPUBaseline) IsPeakHour(hour int) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ph := range b.PeakHours {
		if ph == hour {
			return true
		}
	}
	return false
}

// GetExpectedUsage returns expected CPU usage for given hour
func (b *CPUBaseline) GetExpectedUsage(hour int) float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.HourlySamples[hour] > 0 {
		return b.HourlyAvg[hour]
	}

	// If no data for this hour, return global average
	return b.GlobalAvg
}

// CalculateDeviation calculates % deviation from expected usage
func (b *CPUBaseline) CalculateDeviation(currentUsage float64, hour int) float64 {
	expected := b.GetExpectedUsage(hour)

	if expected == 0 {
		return 0
	}

	return ((currentUsage - expected) / expected) * 100
}

// IsAnomalous checks if current usage is anomalous
// Returns (isAnomalous, deviationPercent, description)
func (b *CPUBaseline) IsAnomalous(currentUsage float64) (bool, float64, string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	now := time.Now()
	hour := now.Hour()
	// expected := b.GetExpectedUsage(hour)
	deviation := b.CalculateDeviation(currentUsage, hour)
	isPeak := b.IsPeakHour(hour)

	// Check for significant increase
	if deviation > b.DeviationThreshold {
		if isPeak {
			return true, deviation, "Usage significantly above peak hour baseline"
		}
		return true, deviation, "Usage significantly above normal baseline"
	}

	// Check for significant decrease during peak hours
	if isPeak && deviation < -30.0 {
		return true, deviation, "Usage significantly below expected peak hour usage"
	}

	// Check for unexpected peak (high usage during off-peak hours)
	if !isPeak && currentUsage > b.GlobalAvg*1.5 && currentUsage > 50.0 {
		return true, deviation, "Unexpected high usage during off-peak hours"
	}

	return false, deviation, ""
}

// ShouldUpdateBaseline checks if we should adapt baseline for sustained changes
// Call this with a rolling window of recent deviations
func (b *CPUBaseline) ShouldUpdateBaseline(recentDeviations []float64, minSamples int) bool {
	if len(recentDeviations) < minSamples {
		return false
	}

	// Calculate average deviation
	sum := 0.0
	for _, dev := range recentDeviations {
		sum += math.Abs(dev)
	}
	avgDev := sum / float64(len(recentDeviations))

	// If sustained deviation is moderate (10-20%), update baseline
	return avgDev > 10.0 && avgDev < 20.0
}

// UpdateHourlyBaseline updates the baseline for a specific hour (adaptive learning)
func (b *CPUBaseline) UpdateHourlyBaseline(hour int, newAverage float64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Weighted average: 70% old data, 30% new observation
	b.HourlyAvg[hour] = b.HourlyAvg[hour]*0.7 + newAverage*0.3
	b.LastUpdated = time.Now()
}

// Marshal converts baseline to JSON
func (b *CPUBaseline) Marshal() ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return json.MarshalIndent(b, "", "  ")
}

// Unmarshal loads baseline from JSON
func (b *CPUBaseline) Unmarshal(data []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	return json.Unmarshal(data, b)
}

// GetSummary returns a human-readable summary of the baseline
func (b *CPUBaseline) GetSummary() string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.IsLearning {
		elapsed := time.Since(b.LearningStarted)
		remaining := b.LearningDuration - elapsed
		return fmt.Sprintf("Learning mode - %s remaining", remaining)
	}

	return fmt.Sprintf("Monitoring active - %d peak hours identified", len(b.PeakHours))
}
