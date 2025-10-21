package iso

import (
	"testing"
	"time"

	"github.com/felipemacedo/cardgen-pro/internal/models"
)

func TestGenerateISO8583Fields(t *testing.T) {
	card := &models.Card{
		PAN:         "4000000000000002",
		MaskedPAN:   "400000******0002",
		Brand:       "Visa",
		ExpiryMonth: 12,
		ExpiryYear:  2027,
		Track2:      "4000000000000002=2712201001234",
	}

	fields := GenerateISO8583Fields(card, 10000, "986")

	// Verify required fields
	requiredFields := []string{"2", "3", "4", "7", "11", "12", "13", "14", "22", "37", "41", "42", "49"}
	
	for _, field := range requiredFields {
		if _, ok := fields[field]; !ok {
			t.Errorf("GenerateISO8583Fields() missing field %s", field)
		}
	}

	// Verify PAN (field 2)
	if fields["2"] != card.PAN {
		t.Errorf("Field 2 (PAN) = %s, want %s", fields["2"], card.PAN)
	}

	// Verify processing code (field 3)
	if fields["3"] != "000000" {
		t.Errorf("Field 3 (Processing Code) = %s, want 000000", fields["3"])
	}

	// Verify amount (field 4) - should be 12 digits
	if len(fields["4"]) != 12 {
		t.Errorf("Field 4 (Amount) length = %d, want 12", len(fields["4"]))
	}

	// Verify expiry date (field 14) - YYMM format
	expectedExpiry := "2712"
	if fields["14"] != expectedExpiry {
		t.Errorf("Field 14 (Expiry) = %s, want %s", fields["14"], expectedExpiry)
	}

	// Verify currency (field 49)
	if fields["49"] != "986" {
		t.Errorf("Field 49 (Currency) = %s, want 986", fields["49"])
	}

	// Verify Track2 (field 35) is present
	if _, ok := fields["35"]; !ok {
		t.Errorf("Field 35 (Track2) missing")
	}
}

func TestGenerateMockAuthRequest(t *testing.T) {
	card := &models.Card{
		PAN:         "4000000000000002",
		Brand:       "Visa",
		ExpiryMonth: 12,
		ExpiryYear:  2027,
	}

	request := GenerateMockAuthRequest(card, 10000, "986")

	// Verify MTI
	if request.MTI != "0100" {
		t.Errorf("MTI = %s, want 0100", request.MTI)
	}

	// Verify fields exist
	if len(request.Fields) == 0 {
		t.Error("Request has no fields")
	}

	// Verify timestamp is recent
	if time.Since(request.Timestamp) > time.Minute {
		t.Error("Request timestamp is too old")
	}
}

func TestGenerateMockAuthResponse(t *testing.T) {
	card := &models.Card{
		PAN:         "4000000000000002",
		Brand:       "Visa",
		ExpiryMonth: 12,
		ExpiryYear:  2027,
	}

	request := GenerateMockAuthRequest(card, 10000, "986")
	response := GenerateMockAuthResponse(request, "00", "Approved")

	// Verify MTI
	if response.MTI != "0110" {
		t.Errorf("MTI = %s, want 0110", response.MTI)
	}

	// Verify response code
	if response.ResponseCode != "00" {
		t.Errorf("ResponseCode = %s, want 00", response.ResponseCode)
	}

	// Verify response text
	if response.ResponseText != "Approved" {
		t.Errorf("ResponseText = %s, want Approved", response.ResponseText)
	}

	// Verify auth code is generated for approved transactions
	if response.AuthCode == "" {
		t.Error("AuthCode should be generated for approved transaction")
	}

	// Verify response contains field 39 (response code)
	if _, ok := response.Fields["39"]; !ok {
		t.Error("Response missing field 39 (response code)")
	}

	// Test declined transaction (no auth code)
	declineResponse := GenerateMockAuthResponse(request, "05", "Do not honor")
	if declineResponse.AuthCode != "" {
		t.Error("AuthCode should be empty for declined transaction")
	}
}

func TestResponseCodes(t *testing.T) {
	// Verify common response codes exist
	commonCodes := []string{"00", "05", "51", "54", "91"}

	for _, code := range commonCodes {
		if _, ok := ResponseCodes[code]; !ok {
			t.Errorf("ResponseCodes missing code %s", code)
		}
	}

	// Verify specific codes
	if ResponseCodes["00"] != "Approved" {
		t.Errorf("ResponseCodes[00] = %s, want Approved", ResponseCodes["00"])
	}

	if ResponseCodes["51"] != "Insufficient funds" {
		t.Errorf("ResponseCodes[51] = %s, want Insufficient funds", ResponseCodes["51"])
	}
}

func TestFormatISO8583(t *testing.T) {
	fields := ISO8583Fields{
		"2":  "4000000000000002",
		"3":  "000000",
		"4":  "000000010000",
		"49": "986",
	}

	formatted := FormatISO8583(fields)

	// Verify output contains field information
	if formatted == "" {
		t.Error("FormatISO8583() returned empty string")
	}

	// Verify contains field values
	expectedSubstrings := []string{"Field 2", "4000000000000002", "Field 49", "986"}
	for _, substr := range expectedSubstrings {
		if !contains(formatted, substr) {
			t.Errorf("FormatISO8583() missing substring: %s", substr)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func BenchmarkGenerateISO8583Fields(b *testing.B) {
	card := &models.Card{
		PAN:         "4000000000000002",
		Brand:       "Visa",
		ExpiryMonth: 12,
		ExpiryYear:  2027,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateISO8583Fields(card, 10000, "986")
	}
}
