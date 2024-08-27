package email

import (
	"fmt"
	"net/smtp"
)

// SendVerificationEmail sends a verification email to the user.
func SendVerificationEmail(email, token string) error {
	from := "your-email@example.com" // Replace with your email
	password := "your-email-password" // Replace with your email password
	to := email
	subject := "Verify Your Email Address"
	body := fmt.Sprintf(
		"Please click the following link to verify your email address: https://yourdomain.com/users/verify-email?token=%s&email=%s",
		token,
		email,
	)

	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	auth := smtp.PlainAuth("", from, password, "smtp.example.com") // Replace with your SMTP server

	return smtp.SendMail("smtp.example.com:587", auth, from, []string{to}, msg) // Replace with your SMTP server address
}
