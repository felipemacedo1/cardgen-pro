package iso

import (
	"fmt"
	"time"

	"github.com/felipemacedo/cardgen-pro/internal/models"
)

// ISO8583Fields represents common ISO-8583 message fields
// This is a SIMPLIFIED representation for testing purposes
//
// DESIGN RATIONALE:
// - ISO-8583 is the international standard for financial transaction messages
// - Contains 128 fields (bitmap-based), but we implement commonly used ones
// - This is NOT a full ISO-8583 implementation; use proper libraries for production
// - Fields are represented as map[string]string for simplicity
//
// COMMON FIELDS:
// 2  - Primary Account Number (PAN)
// 3  - Processing Code
// 4  - Transaction Amount
// 7  - Transmission Date & Time
// 11 - System Trace Audit Number (STAN)
// 12 - Local Transaction Time
// 13 - Local Transaction Date
// 14 - Expiration Date (YYMM)
// 22 - Point of Service Entry Mode
// 35 - Track 2 Data
// 37 - Retrieval Reference Number
// 41 - Card Acceptor Terminal ID
// 42 - Card Acceptor ID Code
// 49 - Transaction Currency Code
// 55 - EMV/ICC Data (chip)
// 95 - Replacement Amounts
type ISO8583Fields map[string]string

// GenerateISO8583Fields generates a map of common ISO-8583 fields for a card
// This simulates an authorization request message (MTI 0100)
func GenerateISO8583Fields(card *models.Card, amount int64, currency string) ISO8583Fields {
	now := time.Now()

	fields := ISO8583Fields{
		"2":  card.PAN,                           // Primary Account Number
		"3":  "000000",                            // Processing Code (purchase)
		"4":  fmt.Sprintf("%012d", amount),        // Transaction Amount (12 digits, padded)
		"7":  now.Format("0102150405"),            // Transmission Date & Time (MMDDhhmmss)
		"11": generateSTAN(),                      // System Trace Audit Number
		"12": now.Format("150405"),                // Local Time (hhmmss)
		"13": now.Format("0102"),                  // Local Date (MMDD)
		"14": fmt.Sprintf("%02d%02d", card.ExpiryYear%100, card.ExpiryMonth), // Expiry (YYMM)
		"22": "051",                               // POS Entry Mode (chip with PIN)
		"37": generateRRN(),                       // Retrieval Reference Number
		"41": "TERM0001",                          // Terminal ID
		"42": "MERCHANT000001",                    // Merchant ID
		"49": currency,                            // Currency Code (e.g., "986" for BRL)
	}

	// Add Track2 if available
	if card.Track2 != "" {
		fields["35"] = card.Track2
	}

	return fields
}

// generateSTAN generates a System Trace Audit Number (6 digits)
func generateSTAN() string {
	return fmt.Sprintf("%06d", time.Now().Unix()%1000000)
}

// generateRRN generates a Retrieval Reference Number (12 chars: YYMMDDHHMMSS)
func generateRRN() string {
	return time.Now().Format("060102150405")
}

// FormatISO8583 formats ISO fields as a readable string (for debugging)
func FormatISO8583(fields ISO8583Fields) string {
	result := "ISO-8583 Fields:\n"
	
	fieldOrder := []string{"2", "3", "4", "7", "11", "12", "13", "14", "22", "35", "37", "41", "42", "49"}
	
	for _, field := range fieldOrder {
		if value, ok := fields[field]; ok {
			result += fmt.Sprintf("  Field %s: %s\n", field, value)
		}
	}
	
	return result
}

// AuthorizationRequest represents a mock authorization request
type AuthorizationRequest struct {
	MTI       string         `json:"mti"`        // Message Type Indicator (e.g., "0100")
	Fields    ISO8583Fields  `json:"fields"`
	Timestamp time.Time      `json:"timestamp"`
}

// AuthorizationResponse represents a mock authorization response
type AuthorizationResponse struct {
	MTI            string         `json:"mti"`              // Message Type Indicator (e.g., "0110")
	Fields         ISO8583Fields  `json:"fields"`
	ResponseCode   string         `json:"response_code"`    // "00" = approved
	ResponseText   string         `json:"response_text"`
	AuthCode       string         `json:"auth_code,omitempty"`
	Timestamp      time.Time      `json:"timestamp"`
}

// GenerateMockAuthRequest generates a mock ISO-8583 authorization request
func GenerateMockAuthRequest(card *models.Card, amount int64, currency string) *AuthorizationRequest {
	return &AuthorizationRequest{
		MTI:       "0100", // Authorization request
		Fields:    GenerateISO8583Fields(card, amount, currency),
		Timestamp: time.Now(),
	}
}

// GenerateMockAuthResponse generates a mock ISO-8583 authorization response
func GenerateMockAuthResponse(request *AuthorizationRequest, responseCode string, responseText string) *AuthorizationResponse {
	// Copy request fields
	responseFields := ISO8583Fields{}
	for k, v := range request.Fields {
		responseFields[k] = v
	}

	// Add response-specific fields
	responseFields["39"] = responseCode // Response Code

	// Generate auth code if approved
	var authCode string
	if responseCode == "00" {
		authCode = fmt.Sprintf("AUTH%06d", time.Now().Unix()%1000000)
	}

	return &AuthorizationResponse{
		MTI:          "0110", // Authorization response
		Fields:       responseFields,
		ResponseCode: responseCode,
		ResponseText: responseText,
		AuthCode:     authCode,
		Timestamp:    time.Now(),
	}
}

// ResponseCodes contains common ISO-8583 response codes
var ResponseCodes = map[string]string{
	"00": "Approved",
	"01": "Refer to card issuer",
	"03": "Invalid merchant",
	"04": "Capture card",
	"05": "Do not honor",
	"12": "Invalid transaction",
	"13": "Invalid amount",
	"14": "Invalid card number",
	"30": "Format error",
	"41": "Lost card",
	"43": "Stolen card",
	"51": "Insufficient funds",
	"54": "Expired card",
	"55": "Incorrect PIN",
	"57": "Transaction not permitted",
	"58": "Transaction not permitted to terminal",
	"61": "Exceeds withdrawal limit",
	"62": "Restricted card",
	"63": "Security violation",
	"65": "Exceeds withdrawal frequency",
	"75": "PIN tries exceeded",
	"91": "Issuer unavailable",
	"96": "System malfunction",
}
