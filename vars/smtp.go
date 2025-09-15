package vars

import "VPSWatchDog/logger"

type MailVars struct {
	SMTPHost string
	SMTPUser string
	SMTPPort int
	SMTPPass string
	MailFrom string
	MailTo   string
}

var MailVar MailVars

func InitMailVars(v MailVars) {
	MailVar.SMTPHost = v.SMTPHost
	MailVar.SMTPUser = v.SMTPUser
	MailVar.SMTPPort = v.SMTPPort
	MailVar.SMTPPass = v.SMTPPass
	MailVar.MailFrom = v.MailFrom
	MailVar.MailTo = v.MailTo

	logger.Info("SMPT credentials initialized successfully")
}
