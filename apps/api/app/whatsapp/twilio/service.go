package twilio

import (
	"fmt"
	"log"
)

// =============================================================================
// HIGH-LEVEL SERVICE FUNCTIONS
// =============================================================================

// SendPasswordReset mengirim password reset code via Twilio WhatsApp
func SendPasswordReset(phone, resetCode string) error {
	client := NewClient()

	if !client.IsConfigured() {
		log.Println("Twilio credentials not configured, skipping password reset message")
		return nil
	}

	to, contentSID, vars := PasswordResetMessage(phone, resetCode)

	resp, err := client.Send(to, contentSID, vars)
	if err != nil {
		return fmt.Errorf("failed to send password reset message: %w", err)
	}

	if resp == nil || resp.Sid == "" {
		return fmt.Errorf("empty response from Twilio")
	}

	log.Printf("Password reset code message sent to %s (sid: %s)", phone, resp.Sid)
	return nil
}

// SendPasswordResetURL mengirim password reset link via Twilio WhatsApp
func SendPasswordResetURL(phone, userName, resetURL, companyName string) error {
	client := NewClient()

	if !client.IsConfigured() {
		log.Println("Twilio credentials not configured, skipping password reset URL message")
		return nil
	}

	to, contentSID, vars := PasswordResetURLMessage(phone, userName, resetURL, companyName)

	resp, err := client.Send(to, contentSID, vars)
	if err != nil {
		return fmt.Errorf("failed to send password reset URL message: %w", err)
	}

	if resp == nil || resp.Sid == "" {
		return fmt.Errorf("empty response from Twilio")
	}

	log.Printf("Password reset URL message sent to %s (sid: %s)", phone, resp.Sid)
	return nil
}

// SendOTP mengirim OTP code via Twilio WhatsApp
func SendOTP(phone, otp string) error {
	client := NewClient()

	if !client.IsConfigured() {
		log.Println("Twilio credentials not configured, skipping OTP message")
		return nil
	}

	to, contentSID, vars := OTPMessage(phone, otp)

	resp, err := client.Send(to, contentSID, vars)
	if err != nil {
		return fmt.Errorf("failed to send OTP message: %w", err)
	}

	if resp == nil || resp.Sid == "" {
		return fmt.Errorf("empty response from Twilio")
	}

	log.Printf("OTP message sent to %s (sid: %s)", phone, resp.Sid)
	return nil
}

// SendWelcome mengirim welcome message via Twilio WhatsApp
func SendWelcome(phone, userName string) error {
	client := NewClient()

	if !client.IsConfigured() {
		return nil
	}

	to, contentSID, vars := WelcomeMessage(phone, userName)

	resp, err := client.Send(to, contentSID, vars)
	if err != nil {
		return fmt.Errorf("failed to send welcome message: %w", err)
	}

	if resp == nil || resp.Sid == "" {
		return fmt.Errorf("empty response from Twilio")
	}

	log.Printf("Welcome message sent to %s (sid: %s)", phone, resp.Sid)
	return nil
}

// SendText mengirim custom text message via Twilio WhatsApp
func SendText(phone, text string) error {
	client := NewClient()

	if !client.IsConfigured() {
		return fmt.Errorf("twilio credentials not configured")
	}

	resp, err := client.SendText(phone, text)
	if err != nil {
		return fmt.Errorf("failed to send text message: %w", err)
	}

	if resp == nil || resp.Sid == "" {
		return fmt.Errorf("empty response from Twilio")
	}

	log.Printf("Text message sent to %s (sid: %s)", phone, resp.Sid)
	return nil
}

// SendAsync mengirim message secara asynchronous (non-blocking)
func SendAsync(fn func() error) {
	go func() {
		if err := fn(); err != nil {
			log.Printf("Error sending async message: %v", err)
		}
	}()
}
