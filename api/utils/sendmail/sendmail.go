package sendmail

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// Mail : structure for mail
type Mail struct {
	ToName   string
	ToAddr   string
	FromName string
	FromAddr string
	Subject  string
	Body     string
}

// Mailer : mail sending function
func Mailer(m Mail) error {
	smtpUser := os.Getenv("EMAiL_USERNAME")
	if smtpUser == "" {
		log.Fatal("Deploy: Unable to retrieve SMTP user.")
	}

	smtpPwd := os.Getenv("EMAIL_PASSWORD")
	if smtpPwd == "" {
		log.Fatal("Deploy: Unable to retrieve SMTP password.")
	}

	mailer := os.Getenv("EMAIL_HOST")
	mailerPort := fmt.Sprintf(":%s", os.Getenv("EMAIL_PORT"))

	msg := []byte("From: " + m.FromName + " <" + m.FromAddr + ">" + "\r\n" +
		"To: " + m.ToAddr + "\r\n" +
		"Subject: " + m.Subject + "\r\n\r\n" +
		m.Body + "\r\n")

	auth := smtp.PlainAuth("", smtpUser, smtpPwd, mailer)
	err := smtp.SendMail(mailer+mailerPort, auth, m.FromAddr, []string{m.ToAddr}, msg)

	return err
}

// m := Mail{	ToAddr: "recipient@example.com",
// FromName: "Sender's Name",
// FromAddr: "sender@example.com",
// Subject:  "Subject line",
// Body:     "Test email. \r\n Test, test." }

// err := mailer(m)
// if err != nil {
// log.Fatal(err)
// }
