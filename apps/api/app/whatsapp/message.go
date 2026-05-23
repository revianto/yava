package whatsapp

import "github.com/revianto/yava/api/helpers"

// Message struktur untuk WhatsApp message
type Message struct {
	MessagingProduct string       `json:"messaging_product"`
	RecipientType    string       `json:"recipient_type"`
	To               string       `json:"to"`
	Type             string       `json:"type"`
	Template         *Template    `json:"template,omitempty"`
	Text             *TextMessage `json:"text,omitempty"`
}

type Template struct {
	Name       string      `json:"name"`
	Language   Language    `json:"language"`
	Components []Component `json:"components,omitempty"`
}

type Language struct {
	Code string `json:"code"`
}

type Component struct {
	Type       string      `json:"type"`
	Parameters []Parameter `json:"parameters,omitempty"`
}

type Parameter struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type TextMessage struct {
	Body string `json:"body"`
}

// =============================================================================
// PASSWORD RESET TEMPLATE
// =============================================================================

// NewPasswordResetMessage membuat message untuk password reset
func NewPasswordResetMessage(phone, resetCode string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               normalizePhone(phone),
		Type:             "template",
		Template: &Template{
			Name: "password_reset",
			Language: Language{
				Code: "id",
			},
			Components: []Component{
				{
					Type: "body",
					Parameters: []Parameter{
						{
							Type: "text",
							Text: resetCode,
						},
					},
				},
			},
		},
	}
}

// =============================================================================
// OTP VERIFICATION TEMPLATE
// =============================================================================

// NewOTPMessage membuat message untuk OTP verification
func NewOTPMessage(phone, otp string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               normalizePhone(phone),
		Type:             "template",
		Template: &Template{
			Name: "otp_verification",
			Language: Language{
				Code: "id",
			},
			Components: []Component{
				{
					Type: "body",
					Parameters: []Parameter{
						{
							Type: "text",
							Text: otp,
						},
					},
				},
			},
		},
	}
}

// =============================================================================
// WELCOME MESSAGE
// =============================================================================

// NewWelcomeMessage membuat welcome message
func NewWelcomeMessage(phone, userName string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               normalizePhone(phone),
		Type:             "template",
		Template: &Template{
			Name: "hello_world",
			Language: Language{
				Code: "id",
			},
		},
	}
}

// =============================================================================
// CUSTOM TEXT MESSAGE
// =============================================================================

// NewTextMessage membuat simple text message
func NewTextMessage(phone, text string) *Message {
	return &Message{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               normalizePhone(phone),
		Type:             "text",
		Text: &TextMessage{
			Body: text,
		},
	}
}

// =============================================================================
// HELPERS
// =============================================================================

// normalizePhone menormalisasi nomor telepon untuk WhatsApp format
// Input bisa: +6282134497226, 6282134497226, atau 082134497226
// Output: 6282134497226
func normalizePhone(phone string) string {
	normalized := helpers.NormalizePhone(phone)

	// Remove + prefix jika ada
	if len(normalized) > 0 && normalized[0:1] == "+" {
		normalized = normalized[1:]
	}

	return normalized
}
