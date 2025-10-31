package watcher

import (
	"fmt"
	"time"

	"VPSWatchDog/intelligence"
	"VPSWatchDog/logger"
	"VPSWatchDog/mailer"
	"github.com/shirou/gopsutil/v4/cpu"
)

var cpuIntelligence *intelligence.CPUIntelligence

// StartIntelligentCPUPoller starts the intelligent CPU monitoring system
func StartIntelligentCPUPoller(
	interval time.Duration,
	baselinePath string,
	learningDuration time.Duration,
	maxAlertsPerHour int,
) {
	// Initialize intelligence engine
	cpuIntelligence = intelligence.NewCPUIntelligence(
		baselinePath,
		learningDuration,
		maxAlertsPerHour,
	)

	// Get initial CPU info
	percentages, err := cpu.Percent(200*time.Millisecond, false)
	if err != nil || len(percentages) == 0 {
		logger.TimedError("Error reading CPU usage:", err)
	}

	usage := 0.0
	if len(percentages) > 0 {
		usage = percentages[0]
	}
	cores, err := cpu.Counts(true)
	if err != nil {
		logger.TimedError("Error reading CPU core counts:", err)
	}

	// Check if we're in learning or monitoring mode
	isLearning, progress, remaining := cpuIntelligence.GetLearningProgress()

	var initMessage string
	if isLearning {
		initMessage = fmt.Sprintf(
			"🧠 Intelligent CPU Poller Initialized (Learning Mode)\n"+
				"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"+
				"Polling Interval: %s\n"+
				"CPU Cores: %d\n"+
				"Current Avg Usage: %.2f%%\n\n"+
				"📚 Learning Phase:\n"+
				"  Progress: %.1f%%\n"+
				"  Time Remaining: %s\n\n"+
				"During learning, the system is collecting baseline data about\n"+
				"your CPU usage patterns. Intelligent monitoring will begin\n"+
				"automatically when learning is complete.",
			interval,
			cores,
			usage,
			progress,
			remaining,
		)
	} else {
		summary := cpuIntelligence.GetBaselineSummary()
		initMessage = fmt.Sprintf(
			"🎯 Intelligent CPU Poller Initialized (Monitoring Mode)\n"+
				"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"+
				"Polling Interval: %s\n"+
				"CPU Cores: %d\n"+
				"Current Avg Usage: %.2f%%\n\n"+
				"📊 Baseline Status:\n"+
				"  %s\n\n"+
				"Max Alerts per Hour: %d\n"+
				"Adaptive monitoring is now active!",
			interval,
			cores,
			usage,
			summary,
			maxAlertsPerHour,
		)
	}

	err = mailer.SendMail("Intelligent CPU Poller Initialized", initMessage)
	if err != nil {
		logger.TimedError("Failed to send CPU init email:", err)
	}

	logger.TimedInfo("Intelligent CPU poller started")

	// Start the monitoring goroutine
	go monitorCPU(interval)
}

func monitorCPU(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		// Get CPU usage
		percentages, err := cpu.Percent(interval, false)
		if err != nil || len(percentages) == 0 {
			logger.TimedError("Error reading CPU usage:", err)
			continue
		}

		avgUsage := percentages[0]

		// Process through intelligence engine
		shouldAlert, message, err := cpuIntelligence.ProcessUsage(avgUsage)
		if err != nil {
			logger.TimedError("Error processing CPU usage:", err)
			continue
		}

		// Send alert if needed
		if shouldAlert {
			subject := "⚠️ CPU Usage Alert - Intelligent Monitor"

			// Check if this is the learning completion message
			if message == "Learning phase completed. Intelligent monitoring now active." {
				subject = "✅ CPU Learning Complete - Monitoring Active"
				message = fmt.Sprintf(
					"🎉 Learning Phase Complete!\n"+
						"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n"+
						"%s\n\n"+
						"The system has finished learning your CPU usage patterns.\n"+
						"Intelligent anomaly detection is now active and will alert you\n"+
						"of significant deviations from normal behavior.\n\n"+
						"Baseline Summary:\n"+
						"  %s",
					message,
					cpuIntelligence.GetBaselineSummary(),
				)
			}

			err = mailer.SendMail(subject, message)
			if err != nil {
				logger.TimedError("Failed to send CPU alert email:", err)
			} else {
				logger.TimedInfo("CPU alert email sent:", subject)
			}
		}

		// Periodic status logging (every 12 iterations, ~1 minute if interval is 5s)
		// This is commented out to reduce log spam, but you can enable it for debugging
		// if iteration % 12 == 0 {
		//     isLearning, progress, remaining := cpuIntelligence.GetLearningProgress()
		//     if isLearning {
		//         logger.TimedDebug(fmt.Sprintf("CPU Learning: %.1f%% complete, %s remaining", progress, remaining))
		//     }
		// }
	}
}

// GetCPUIntelligenceStatus returns current status (for health checks or API)
func GetCPUIntelligenceStatus() map[string]interface{} {
	if cpuIntelligence == nil {
		return map[string]interface{}{
			"initialized": false,
		}
	}

	isLearning, progress, remaining := cpuIntelligence.GetLearningProgress()
	alertCount, isEnabled, timeUntilReset := cpuIntelligence.AlertState.GetStatus()

	status := map[string]interface{}{
		"initialized":           true,
		"is_learning":           isLearning,
		"learning_progress":     progress,
		"learning_remaining":    remaining.String(),
		"alerts_sent_this_hour": alertCount,
		"monitoring_enabled":    isEnabled,
		"time_until_reset":      timeUntilReset.String(),
	}

	if !isLearning {
		status["baseline_summary"] = cpuIntelligence.GetBaselineSummary()
	}

	return status
}
