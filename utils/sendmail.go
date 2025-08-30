package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_FROM_EMAIL")
	fromName := os.Getenv("SMTP_FROM_NAME")
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	msg := fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		fromName, from, to, subject, body)

	cSMTP, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer cSMTP.Close()

	if err = cSMTP.Mail(from); err != nil {
		return err
	}
	if err = cSMTP.Rcpt(to); err != nil {
		return err
	}
	wc, err := cSMTP.Data()
	if err != nil {
		return err
	}
	_, err = wc.Write([]byte(msg))
	if err != nil {
		return err
	}
	return wc.Close()
}
