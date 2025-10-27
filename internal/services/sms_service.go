package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// SMSService handles SMS sending operations
type SMSService struct {
	smsURL    string
	apiKey    string
	senderID  string
	enabled   bool
}

// NewSMSService creates a new SMS service
func NewSMSService(smsURL, apiKey, senderID string) *SMSService {
	enabled := smsURL != "" && apiKey != "" && senderID != ""
	
	if !enabled {
		log.Println("SMS service is disabled (missing configuration)")
	}

	return &SMSService{
		smsURL:   smsURL,
		apiKey:   apiKey,
		senderID: senderID,
		enabled:  enabled,
	}
}

// SMSPayload represents the JSON payload for SMS API (matches Spring Boot Constant.java)
// Field names match: OTP_SMS_API_KEY, OTP_SMS_SENDER_ID, OTP_SMS_CLIENT_NUMBER_KEY, OTP_SMS_CLIENT_MESSAGE_KEY
type SMSPayload struct {
	APIKey   string `json:"api_key"`    // Constant.OTP_SMS_API_KEY
	SenderID string `json:"senderid"`   // Constant.OTP_SMS_SENDER_ID
	Number   string `json:"number"`     // Constant.OTP_SMS_CLIENT_NUMBER_KEY (not "contacts")
	Message  string `json:"message"`    // Constant.OTP_SMS_CLIENT_MESSAGE_KEY (not "msg")
}

// SMSResponse represents the response from SMS API
type SMSResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SendSMS sends an SMS to the specified phone number (matches Spring Boot SmsOtpSenderStrategy)
func (s *SMSService) SendSMS(phoneNumber, message string) error {
	if !s.enabled {
		// In development mode, just log the OTP
		log.Printf("[SMS] Would send to %s: %s", phoneNumber, message)
		return nil
	}

	// Build JSON payload (matches Spring Boot's getSmsPayload)
	payload := SMSPayload{
		APIKey:   s.apiKey,
		SenderID: s.senderID,
		Number:   phoneNumber,
		Message:  message,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal SMS payload: %w", err)
	}

	// Create POST request with JSON content type (matches Spring Boot)
	req, err := http.NewRequest("POST", s.smsURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create SMS request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read SMS response: %w", err)
	}

	// Parse JSON response
	var smsResp SMSResponse
	if err := json.Unmarshal(body, &smsResp); err != nil {
		// If JSON parsing fails, check HTTP status as fallback
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to send SMS: HTTP %d - %s", resp.StatusCode, string(body))
		}
		log.Printf("Warning: Could not parse SMS response JSON: %v", err)
	}

	// Validate response (matches Spring Boot's ValidationUtils.otpResponseValidation)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send SMS: HTTP %d - %s", resp.StatusCode, smsResp.Message)
	}

	// Additional validation based on response status field if present
	if smsResp.Status != "" && smsResp.Status != "success" && smsResp.Status != "SUCCESS" {
		return fmt.Errorf("SMS API error: %s", smsResp.Message)
	}

	log.Printf("SMS sent successfully to %s", phoneNumber)
	return nil
}

// IsEnabled returns whether SMS service is enabled
func (s *SMSService) IsEnabled() bool {
	return s.enabled
}

// Ensure SMSService implements SMSServiceInterface
var _ SMSServiceInterface = (*SMSService)(nil)
