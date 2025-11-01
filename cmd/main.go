package main

import (
	"VPSWatchDog/logger"
	"VPSWatchDog/mailer"
	"VPSWatchDog/vars"
	"VPSWatchDog/watcher"
	"flag"
	"fmt"
	"time"
)

func testMail() error {
	subject := "[TEST] VPS Watchdog SMTP Test"
	body := fmt.Sprintf(
		"This is a test email from VPS Watchdog program.\n" +
			"If you received this, SMTP settings for VPS alert mail are working correctly.",
	)

	return mailer.SendMail(subject, body)
}

func main() {
	// SMTP config
	smtpHost := flag.String("smtpHost", "smtp.example.com", "SMTP host")
	smtpPort := flag.Int("smtpPort", 587, "SMTP port")
	smtpUser := flag.String("smtpUser", "user@example.com", "SMTP username")
	smtpPass := flag.String("smtpPass", "password", "SMTP password")
	mailFrom := flag.String("from", "alert@example.com", "From email")
	mailTo := flag.String("to", "admin@example.com", "To email")

	// Intelligence config
	useIntelligent := flag.Bool("intelligent", true, "Use intelligent monitoring (default: true)")
	learningDuration := flag.String("learningDuration", "168h", "Duration for learning phase (e.g., 168h = 7 days)")
	baselinePath := flag.String("baselinePath", "./data/baselines", "Path to store baseline data")
	maxAlertsPerHour := flag.Int("maxAlertsPerHour", 2, "Maximum alerts per hour per monitor")

	// Polling intervals
	cpuInterval := flag.Duration("cpuInterval", 5*time.Second, "CPU polling interval")
	memInterval := flag.Duration("memInterval", 10*time.Second, "Memory polling interval")
	diskInterval := flag.Duration("diskInterval", 30*time.Second, "Disk polling interval")
	netInterval := flag.Duration("netInterval", 15*time.Second, "Network polling interval")
	swapInterval := flag.Duration("swapInterval", 20*time.Second, "Swap polling interval")
	loadInterval := flag.Duration("loadInterval", 20*time.Second, "Load polling interval")

	flag.Parse()

	vars.InitMailVars(vars.MailVars{
		SMTPHost: *smtpHost,
		SMTPPort: *smtpPort,
		SMTPUser: *smtpUser,
		SMTPPass: *smtpPass,
		MailFrom: *mailFrom,
		MailTo:   *mailTo,
	})

	if err := testMail(); err != nil {
		logger.Panic("SMTP test failed:", err)
	}

	logger.Okay("SMTP test succeeded: test mail sent")

	// Parse learning duration
	learnDuration, err := time.ParseDuration(*learningDuration)
	if err != nil {
		logger.Panic("Invalid learningDuration format:", err)
	}

	// Start monitoring based on mode
	if *useIntelligent {
		logger.Info("🧠 Starting Intelligent Monitoring System")
		logger.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		logger.Info("Configuration:")
		logger.Info("  Learning Duration:", learnDuration)
		logger.Info("  Baseline Path:", *baselinePath)
		logger.Info("  Max Alerts/Hour:", *maxAlertsPerHour)
		logger.Info("  CPU Interval:", *cpuInterval)

		// Start intelligent CPU monitoring
		watcher.StartIntelligentCPUPoller(
			*cpuInterval,
			*baselinePath,
			learnDuration,
			*maxAlertsPerHour,
		)

		// TODO: Implement intelligent versions for other monitors
		// For now, keep using the old monitors for other metrics
		watcher.StartMemPoller(*memInterval)
		watcher.StartDiskPoller(*diskInterval)
		watcher.StartNetPoller(*netInterval)
		watcher.StartSwapPoller(*swapInterval)
		watcher.StartLoadPoller(*loadInterval)
	} else {
		logger.Info("📊 Starting Legacy Monitoring System")

		// Use old threshold-based monitoring
		watcher.GetCPUIntelligenceStatus()
		watcher.StartMemPoller(*memInterval)
		watcher.StartDiskPoller(*diskInterval)
		watcher.StartNetPoller(*netInterval)
		watcher.StartSwapPoller(*swapInterval)
		watcher.StartLoadPoller(*loadInterval)
	}

	logger.Okay("All monitors started successfully")

	// Keep program alive
	select {}
}
