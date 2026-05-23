package meta

import (
	"fmt"
	"log"
)

// =============================================================================
// HIGH-LEVEL SERVICE FUNCTIONS
// =============================================================================

// SendPasswordReset mengirim password reset code via Meta WhatsApp
func SendPasswordReset(phone, resetCode string) error {
	client := NewClient()

	if !client.IsConfigured() {
		// Skip silently jika tidak ada credentials
		log.Println("Meta WhatsApp credentials not configured, skipping password reset message")
		return nil
	}

	msg := NewPasswordResetMessage(phone, resetCode)
	resp, err := client.Send(msg)

	if err != nil {
		return fmt.Errorf("failed to send password reset message: %w", err)
	}

	if resp == nil || len(resp.Messages) == 0 {
		return fmt.Errorf("no messages in response")
	}

	log.Printf("Password reset message sent to %s via Meta (message_id: %s)", phone, resp.Messages[0].ID)
	return nil
}

// SendOTP mengirim OTP code via Meta WhatsApp
func SendOTP(phone, otp string) error {
	client := NewClient()

	if !client.IsConfigured() {
		log.Println("Meta WhatsApp credentials not configured, skipping OTP message")
		return nil
	}

	msg := NewOTPMessage(phone, otp)
	resp, err := client.Send(msg)

	if err != nil {
		return fmt.Errorf("failed to send OTP message: %w", err)
	}

	if resp == nil || len(resp.Messages) == 0 {
		return fmt.Errorf("no messages in response")
	}

	log.Printf("OTP message sent to %s via Meta (message_id: %s)", phone, resp.Messages[0].ID)
	return nil
}

// SendWelcome mengirim welcome message via Meta WhatsApp
func SendWelcome(phone, userName string) error {
	client := NewClient()

	if !client.IsConfigured() {
		return nil
	}

	msg := NewWelcomeMessage(phone, userName)
	resp, err := client.Send(msg)

	if err != nil {
		return fmt.Errorf("failed to send welcome message: %w", err)
	}

	if resp == nil || len(resp.Messages) == 0 {
		return fmt.Errorf("no messages in response")
	}

	log.Printf("Welcome message sent to %s via Meta (message_id: %s)", phone, resp.Messages[0].ID)
	return nil
}

// SendCustomText mengirim custom text message via Meta WhatsApp
func SendCustomText(phone, text string) error {
	client := NewClient()

	if !client.IsConfigured() {
		return fmt.Errorf("meta whatsapp credentials not configured")
	}

	msg := NewTextMessage(phone, text)
	resp, err := client.Send(msg)

	if err != nil {
		return fmt.Errorf("failed to send text message: %w", err)
	}

	if resp == nil || len(resp.Messages) == 0 {
		return fmt.Errorf("no messages in response")
	}

	log.Printf("Text message sent to %s via Meta (message_id: %s)", phone, resp.Messages[0].ID)
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
