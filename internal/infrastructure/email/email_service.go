package email

import (
	"fmt"
	"net/smtp"

	"loan-tracker-api/config"
)

// SendVerificationEmail sends a verification email to the user.
func SendVerificationEmail(email, token string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	from := cfg.EmailUser         // Sender email address
	password := cfg.EmailPassword // Sender email password
	smtpServer := cfg.EmailHost   // SMTP server address
	smtpPort := cfg.EmailPort     // SMTP server port

	to := email
	subject := "Verify Your Email Address"
	body := fmt.Sprintf(
		"Please click the following link to verify your email address: %s/users/verify-email?token=%s&email=%s",
		cfg.AppName, token, email, // You can replace cfg.AppName with the domain if you have a domain variable.
	)

	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	auth := smtp.PlainAuth("", from, password, smtpServer)

	err = smtp.SendMail(fmt.Sprintf("%s:%d", smtpServer, smtpPort), auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %v", err)
	}

	return nil
}

// SendPasswordResetEmail sends a password reset email to the user.
func SendPasswordResetEmail(to, resetLink string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	from := cfg.EmailUser         // Sender email address
	password := cfg.EmailPassword // Sender email password
	smtpServer := cfg.EmailHost   // SMTP server address
	smtpPort := cfg.EmailPort     // SMTP server port

	subject := "Password Reset Request"
	body := fmt.Sprintf("Please use the following link to reset your password: %s", resetLink)

	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	auth := smtp.PlainAuth("", from, password, smtpServer)

	err = smtp.SendMail(fmt.Sprintf("%s:%d", smtpServer, smtpPort), auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send password reset email: %v", err)
	}

	return nil
}
