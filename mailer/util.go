package mailer

import (
	"VPSWatchDog/vars"
	gomail "gopkg.in/gomail.v2"
)

func SendMail(subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", vars.MailVar.MailFrom)
	m.SetHeader("To", vars.MailVar.MailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(vars.MailVar.SMTPHost, vars.MailVar.SMTPPort, vars.MailVar.SMTPUser, vars.MailVar.SMTPPass)
	return d.DialAndSend(m)
}
