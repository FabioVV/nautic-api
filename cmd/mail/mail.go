package mail

import (
	"fmt"
	"net/smtp"
	"os"

	"gopkg.in/gomail.v2"
)

func SendSystemEmail(to string, subject string, body string) bool {
	message := gomail.NewMessage()

	message.SetHeader("From", os.Getenv("MAIL_USERNAME"))
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	message.SetBody("text/html", fmt.Sprintf(`
        <html>
            <body>
                <h1>Nautic Systems</h1>
                <p><b>Olá!</b></p>
				<br/>
                %s
            </body>
        </html>
    `, body))

	message.AddAlternative("text/plain", fmt.Sprintf("Nautic Systems. \n\n %s \n\n %s", subject, body))

	dialer := gomail.NewDialer(os.Getenv("MAIL_HOST"), 465, os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD"))

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		return false
	}

	return true
}

func SendSystemSimpleEmail(from string, to []string, subject string, body string) bool {
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")

	auth := smtp.PlainAuth("", username, password, host)

	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"From: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n",
		from, to, subject, body,
	))

	if err := smtp.SendMail(host+":"+port, auth, from, to, msg); err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		return false
	}

	return true
}
