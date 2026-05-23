package twilio

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/revianto/yava/api/helpers"
)

// Client untuk Twilio WhatsApp API
type Client struct {
	AccountSID  string
	AuthToken   string
	PhoneNumber string
	APIUrl      string
}

// Response dari Twilio API
type SendResponse struct {
	Sid                 string `json:"sid"`
	DateCreated         string `json:"date_created"`
	DateUpdated         string `json:"date_updated"`
	DateSent            string `json:"date_sent"`
	AccountSid          string `json:"account_sid"`
	To                  string `json:"to"`
	From                string `json:"from"`
	MessagingServiceSid string `json:"messaging_service_sid"`
	Body                string `json:"body"`
	Status              string `json:"status"`
	NumMedia            string `json:"num_media"`
	NumSegments         string `json:"num_segments"`
	Direction           string `json:"direction"`
	ApiVersion          string `json:"api_version"`
	Price               string `json:"price"`
	PriceUnit           string `json:"price_unit"`
	ErrorCode           *int   `json:"error_code"`
	ErrorMessage        string `json:"error_message"`
}

// NewClient membuat instance Twilio client
func NewClient() *Client {
	accountSID := helpers.GetEnv("TWILIO_ACCOUNT_SID")
	authToken := helpers.GetEnv("TWILIO_AUTH_TOKEN")
	phoneNumber := helpers.GetEnv("TWILIO_WHATSAPP_NUMBER")

	return &Client{
		AccountSID:  accountSID,
		AuthToken:   authToken,
		PhoneNumber: phoneNumber,
		APIUrl:      "https://api.twilio.com/2010-04-01",
	}
}

// IsConfigured check apakah credentials sudah di-setup
func (c *Client) IsConfigured() bool {
	return c.AccountSID != "" && c.AuthToken != "" && c.PhoneNumber != ""
}

// Send mengirim message via Twilio WhatsApp
func (c *Client) Send(to, contentSid string, contentVariables map[string]string) (*SendResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("twilio credentials not configured")
	}

	// Prepare data
	data := url.Values{}
	data.Set("To", fmt.Sprintf("whatsapp:%s", normalizePhone(to)))
	data.Set("From", c.PhoneNumber)
	data.Set("ContentSid", contentSid)

	// Add content variables if provided
	if len(contentVariables) > 0 {
		varsJSON := `{`
		i := 1
		for _, val := range contentVariables {
			if i > 1 {
				varsJSON += `,`
			}
			varsJSON += fmt.Sprintf(`"%d":"%s"`, i, val)
			i++
		}
		varsJSON += `}`
		data.Set("ContentVariables", varsJSON)
	}

	// Create request
	url := fmt.Sprintf("%s/Accounts/%s/Messages.json", c.APIUrl, c.AccountSID)
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add auth header (Basic Auth)
	auth := base64.StdEncoding.EncodeToString([]byte(c.AccountSID + ":" + c.AuthToken))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response (basic parsing)
	var respData SendResponse
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("twilio api error: status %d - %s", resp.StatusCode, string(body))
	}

	return &respData, nil
}

// SendText mengirim plain text message via Twilio
func (c *Client) SendText(to, text string) (*SendResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("twilio credentials not configured")
	}

	data := url.Values{}
	data.Set("To", fmt.Sprintf("whatsapp:%s", normalizePhone(to)))
	data.Set("From", c.PhoneNumber)
	data.Set("Body", text)

	url := fmt.Sprintf("%s/Accounts/%s/Messages.json", c.APIUrl, c.AccountSID)
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(c.AccountSID + ":" + c.AuthToken))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("twilio api error: status %d - %s", resp.StatusCode, string(body))
	}

	return &SendResponse{}, nil
}
