package mail

import (
	"crypto/tls"
	"fmt"
	"go-notify/dto"
	"net/smtp"
	"os"
)

const port = "587" //25 or 587
const smtpDomain = "smtp.gmail.com"

var M *Mail

type Mail struct {
	Auth     smtp.Auth
	hostname string
	From     string
}

func NewMail() *Mail {
	return &Mail{
		hostname: smtpDomain,
	}
}

func (m *Mail) SetAuth() {
	username := os.Getenv("GOOGLE_MAIL")
	pass := os.Getenv("GOOGLE_APP_PASSWORD")
	m.From = username
	m.Auth = smtp.PlainAuth("", username, pass, m.hostname)
}

func (m *Mail) SendMail(content dto.Body) error {
	addr := m.hostname + ":" + port
	// Create a TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // Set to true only if needed for testing
		ServerName:         "smtp.gmail.com",
	}
	// client instance
	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()
	// Initiate TLS handshake
	if err = client.StartTLS(tlsConfig); err != nil {
		return err
	}
	// Authenticate with the server
	if err = client.Auth(m.Auth); err != nil {
		return ErrAuth
	}
	// Set the sender and recipient
	if err = client.Mail(m.From); err != nil {
		return ErrSender
	}
	if err = client.Rcpt(content.Receiver); err != nil {
		return ErrRecipient
	}
	// Write the email body
	writer, err := client.Data()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(writer, "Subject: Test Email\r\n\r\n %s.", content.Message)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	// Send the QUIT command and close the connection
	if err = client.Quit(); err != nil {
		return err
	}
	return nil
}
