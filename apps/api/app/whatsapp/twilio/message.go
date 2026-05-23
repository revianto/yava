package twilio

import "github.com/revianto/yava/api/helpers"

// =============================================================================
// TWILIO TEMPLATE CONTENT SIDS
// =============================================================================

// TemplateContentSID mapping untuk template yang sudah dibuat di Twilio
type TemplateConfig struct {
	ContentSID  string
	Description string
}

var Templates = map[string]TemplateConfig{
	"password_reset": {
		ContentSID:  "HX229f5a04fd0510ce1b071852155d3e75", // Old: 6-digit code only
		Description: "Password reset code",
	},
	"password_reset_url": {
		ContentSID:  "HX91ca0e965a050ed786172a7b6da34cc5",
		Description: "Password reset dengan URL link",
	},
	"otp_verification": {
		ContentSID:  "", // Update sesuai Twilio Anda
		Description: "OTP verification code",
	},
	"welcome": {
		ContentSID:  "", // Update sesuai Twilio Anda
		Description: "Welcome message",
	},
}

// =============================================================================
// MESSAGE BUILDERS
// =============================================================================

// PasswordResetMessage membuat password reset message dengan kode untuk Twilio
func PasswordResetMessage(phone, resetCode string) (string, string, map[string]string) {
	config := Templates["password_reset"]
	return phone, config.ContentSID, map[string]string{
		"1": resetCode,
	}
}

// PasswordResetURLMessage membuat password reset message dengan URL link untuk Twilio
// Template parameters:
// {{1}} = userName
// {{2}} = resetURL
// {{3}} = companyName
func PasswordResetURLMessage(phone, userName, resetURL, companyName string) (string, string, map[string]string) {
	config := Templates["password_reset_url"]
	return phone, config.ContentSID, map[string]string{
		"1": userName,
		"2": resetURL,
		"3": companyName,
	}
}

// OTPMessage membuat OTP message untuk Twilio
func OTPMessage(phone, otp string) (string, string, map[string]string) {
	config := Templates["otp_verification"]
	return phone, config.ContentSID, map[string]string{
		"1": otp,
	}
}

// WelcomeMessage membuat welcome message untuk Twilio
func WelcomeMessage(phone, userName string) (string, string, map[string]string) {
	config := Templates["welcome"]
	return phone, config.ContentSID, map[string]string{
		"1": userName,
	}
}

// =============================================================================
// HELPERS
// =============================================================================

// normalizePhone menormalisasi nomor telepon untuk WhatsApp format
func normalizePhone(phone string) string {
	normalized := helpers.NormalizePhone(phone)

	// Remove + prefix jika ada
	if len(normalized) > 0 && normalized[0:1] == "+" {
		normalized = normalized[1:]
	}

	return normalized
}
