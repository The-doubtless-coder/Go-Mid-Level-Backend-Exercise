package clients

import (
	"errors"
	"fmt"
	"github.com/go-gomail/gomail"
	"log"
	"net/smtp"
	"os"
	"strconv"
)

func SendAdminEmail(subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	m.SetHeader("To", os.Getenv("EMAIL_ADMIN_RECEIVER"))
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))

	d := gomail.NewDialer(
		os.Getenv("EMAIL_HOST"),
		port,
		os.Getenv("EMAIL_USERNAME"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	if err := d.DialAndSend(m); err != nil {
		return errors.New(fmt.Sprintf("Failed to send email: %v", err))
	}
	return nil
}

func SendAdminEmailAsync(subject, body string) {
	go func() {
		from := os.Getenv("EMAIL_SENDER")
		password := os.Getenv("EMAIL_PASSWORD")
		to := os.Getenv("EMAIL_ADMIN_RECEIVER") // set in env

		smtpHost := os.Getenv("EMAIL_HOST") // e.g., "smtp.gmail.com"
		smtpPort := os.Getenv("EMAIL_PORT") // e.g., "587"

		auth := smtp.PlainAuth("", from, password, smtpHost)

		msg := []byte("To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-version: 1.0;\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
			body + "\r\n")

		addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

		err := smtp.SendMail(addr, auth, from, []string{to}, msg)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			return
		}

		log.Println("Admin email sent successfully")
	}()
}
