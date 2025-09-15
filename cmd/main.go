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

	// Start polling every 5 seconds
	watcher.StartCPUPoller(5 * time.Second)
	watcher.StartMemPoller(10 * time.Second)
	watcher.StartDiskPoller(30 * time.Second)
	watcher.StartNetPoller(15 * time.Second)
	watcher.StartSwapPoller(20 * time.Second)
	watcher.StartLoadPoller(20 * time.Second)

	// Keep program alive
	select {}
}
