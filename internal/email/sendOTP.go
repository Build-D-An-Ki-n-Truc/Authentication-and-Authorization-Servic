package email

import (
	"math/rand/v2"
	"strconv"

	"github.com/Build-D-An-Ki-n-Truc/auth/internal/config"
	"gopkg.in/gomail.v2"
)

func randomOTP() int {
	random := rand.IntN(999999-100000) + 100000
	return random
}

var CFG = config.LoadConfig()

// send email with OTP using gmail smtp server
// this will return OTP and error if any
func SendEmail(email string) (string, error) {
	// Configure the SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "kiznlh", CFG.EmailPassword)

	m := gomail.NewMessage()

	// Set email headers
	m.SetHeader("From", "kiznlh@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Mobile Verification OTP")

	// Set email body
	otp := strconv.Itoa(randomOTP())
	body := "This is your OTP: " + otp
	m.SetBody("text/plain", body)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return "", err
	}

	return otp, nil
}
