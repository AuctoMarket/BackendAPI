package utils

import (
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendOtpMail(email string, otp string) error {
	from := mail.NewEmail("Aucto Admin", "admin@aucto.io")
	subject := "Verify your email address."
	to := mail.NewEmail("New collector", email)
	plainTextContent := "Please authenticate your email address. Your authentication OTP is: " + otp
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, "")
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)

	if err != nil {
		LogError(err, "Error in sending mail")
		return err
	}

	return nil
}
