package meta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/revianto/yava/api/helpers"
)

// Client untuk Meta WhatsApp Business API
type Client struct {
	Token   string
	PhoneID string
	APIUrl  string
}

// Response dari Meta API
type APIResponse struct {
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
	Contacts []struct {
		Input string `json:"input"`
		WaID  string `json:"wa_id"`
	} `json:"contacts"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    int    `json:"code"`
	} `json:"error"`
}

// NewClient membuat instance Meta WhatsApp client
func NewClient() *Client {
	token := helpers.GetEnv("WHATSAPP_API_TOKEN")
	phoneID := helpers.GetEnv("WHATSAPP_PHONE_ID")

	return &Client{
		Token:   token,
		PhoneID: phoneID,
		APIUrl:  "https://graph.instagram.com/v22.0",
	}
}

// IsConfigured check apakah credentials sudah di-setup
func (c *Client) IsConfigured() bool {
	return c.Token != "" && c.PhoneID != ""
}

// Send mengirim message ke WhatsApp
func (c *Client) Send(msg *Message) (*APIResponse, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("meta whatsapp credentials not configured")
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	url := fmt.Sprintf("%s/%s/messages", c.APIUrl, c.PhoneID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")

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

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check untuk error dari API
	if apiResp.Error != nil {
		return nil, fmt.Errorf("meta api error: %s (code: %d)", apiResp.Error.Message, apiResp.Error.Code)
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("meta api error: status %d - %s", resp.StatusCode, string(body))
	}

	return &apiResp, nil
}
