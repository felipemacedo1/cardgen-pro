package generator

import (
	"testing"
)

func TestValidateLuhn(t *testing.T) {
	tests := []struct {
		name  string
		pan   string
		valid bool
	}{
		{"Valid Visa 16", "4000000000000002", true},
		{"Valid Visa 13", "4000000000006", true},
		{"Valid Mastercard", "5100000000000016", true},
		{"Valid Amex", "340000000000009", true},
		{"Invalid checksum", "4000000000000001", false},
		{"Too short", "123", false},
		{"Too long", "12345678901234567890", false},
		{"Non-numeric", "400000000000000X", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateLuhn(tt.pan)
			if result != tt.valid {
				t.Errorf("ValidateLuhn(%s) = %v, want %v", tt.pan, result, tt.valid)
			}
		})
	}
}

func TestCalculateLuhnCheckDigit(t *testing.T) {
	tests := []struct {
		name       string
		partialPAN string
		expected   int
	}{
		{"Visa 16 ending in 2", "400000000000000", 2},
		{"Visa 13 ending in 6", "400000000000", 6},
		{"Mastercard ending in 6", "510000000000001", 6},
		{"Amex ending in 9", "34000000000000", 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateLuhnCheckDigit(tt.partialPAN)
			if result != tt.expected {
				t.Errorf("CalculateLuhnCheckDigit(%s) = %d, want %d", tt.partialPAN, result, tt.expected)
			}
		})
	}
}

func TestAppendLuhnCheckDigit(t *testing.T) {
	tests := []struct {
		name       string
		partialPAN string
		expected   string
	}{
		{"Visa 16", "400000000000000", "4000000000000002"},
		{"Visa 13", "400000000000", "4000000000006"},
		{"Mastercard", "510000000000001", "5100000000000016"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AppendLuhnCheckDigit(tt.partialPAN)
			if result != tt.expected {
				t.Errorf("AppendLuhnCheckDigit(%s) = %s, want %s", tt.partialPAN, result, tt.expected)
			}
			
			// Verify result passes Luhn
			if !ValidateLuhn(result) {
				t.Errorf("AppendLuhnCheckDigit(%s) produced invalid PAN: %s", tt.partialPAN, result)
			}
		})
	}
}

func TestGeneratePAN(t *testing.T) {
	tests := []struct {
		name      string
		bin       string
		length    int
		shouldErr bool
	}{
		{"Valid Visa 16", "400000", 16, false},
		{"Valid Visa 13", "400000", 13, false},
		{"Valid Visa 19", "400000", 19, false},
		{"Valid Mastercard", "510000", 16, false},
		{"Valid Amex", "340000", 15, false},
		{"BIN too short", "4000", 16, true},
		{"Length too short", "400000", 12, true},
		{"Length too long", "400000", 20, true},
		{"BIN longer than PAN", "40000000000000000", 16, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pan, err := GeneratePAN(tt.bin, tt.length)
			
			if tt.shouldErr {
				if err == nil {
					t.Errorf("GeneratePAN(%s, %d) expected error but got none", tt.bin, tt.length)
				}
				return
			}

			if err != nil {
				t.Errorf("GeneratePAN(%s, %d) unexpected error: %v", tt.bin, tt.length, err)
				return
			}

			// Verify length
			if len(pan) != tt.length {
				t.Errorf("GeneratePAN(%s, %d) length = %d, want %d", tt.bin, tt.length, len(pan), tt.length)
			}

			// Verify starts with BIN
			if pan[:len(tt.bin)] != tt.bin {
				t.Errorf("GeneratePAN(%s, %d) does not start with BIN: %s", tt.bin, tt.length, pan)
			}

			// Verify Luhn
			if !ValidateLuhn(pan) {
				t.Errorf("GeneratePAN(%s, %d) produced invalid PAN: %s", tt.bin, tt.length, pan)
			}
		})
	}
}

func TestGenerateExpiry(t *testing.T) {
	for i := 0; i < 100; i++ {
		month, year := GenerateExpiry()

		if month < 1 || month > 12 {
			t.Errorf("GenerateExpiry() month = %d, want 1-12", month)
		}

		// Should be in the future
		if year < 2025 || year > 2030 {
			t.Errorf("GenerateExpiry() year = %d, want 2025-2030", year)
		}
	}
}

func TestGenerateDeterministicCVC(t *testing.T) {
	tests := []struct {
		name     string
		pan      string
		month    string
		year     string
		secret   string
		cvcLen   int
		shouldErr bool
	}{
		{"Visa 3 digits", "4000000000000002", "12", "2027", "test-secret", 3, false},
		{"Mastercard 3 digits", "5100000000000016", "06", "2026", "test-secret", 3, false},
		{"Amex 4 digits", "340000000000009", "03", "2025", "test-secret", 4, false},
		{"No secret", "4000000000000002", "12", "2027", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cvc, err := GenerateDeterministicCVC(tt.pan, tt.month, tt.year, tt.secret)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("GenerateDeterministicCVC() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateDeterministicCVC() unexpected error: %v", err)
				return
			}

			// Verify length
			if len(cvc) != tt.cvcLen {
				t.Errorf("GenerateDeterministicCVC() length = %d, want %d", len(cvc), tt.cvcLen)
			}

			// Verify all digits
			for _, c := range cvc {
				if c < '0' || c > '9' {
					t.Errorf("GenerateDeterministicCVC() contains non-digit: %c", c)
				}
			}

			// Verify determinism (same input = same output)
			cvc2, _ := GenerateDeterministicCVC(tt.pan, tt.month, tt.year, tt.secret)
			if cvc != cvc2 {
				t.Errorf("GenerateDeterministicCVC() not deterministic: %s != %s", cvc, cvc2)
			}

			// Verify different secret = different CVC
			cvc3, _ := GenerateDeterministicCVC(tt.pan, tt.month, tt.year, "different-secret")
			if cvc == cvc3 {
				t.Errorf("GenerateDeterministicCVC() same CVC for different secrets")
			}
		})
	}
}

func TestGenerateTrack2(t *testing.T) {
	tests := []struct {
		name        string
		pan         string
		month       int
		year        int
		serviceCode string
	}{
		{"Visa", "4000000000000002", 12, 2027, "201"},
		{"Mastercard", "5100000000000016", 6, 2026, "201"},
		{"Amex", "340000000000009", 3, 2025, "201"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			track2 := GenerateTrack2(tt.pan, tt.month, tt.year, tt.serviceCode)

			// Verify format: PAN=YYMM<ServiceCode><Discretionary>
			if len(track2) < len(tt.pan)+1+4+3+4 {
				t.Errorf("GenerateTrack2() too short: %s", track2)
			}

			// Verify starts with PAN=
			expected := tt.pan + "="
			if track2[:len(expected)] != expected {
				t.Errorf("GenerateTrack2() does not start with '%s': %s", expected, track2)
			}

			// Verify contains service code
			if track2[len(tt.pan)+1+4:len(tt.pan)+1+4+3] != tt.serviceCode {
				t.Errorf("GenerateTrack2() service code mismatch in: %s", track2)
			}
		})
	}
}

func TestMaskPAN(t *testing.T) {
	tests := []struct {
		name     string
		pan      string
		expected string
	}{
		{"Visa 16", "4000000000000002", "400000******0002"},
		{"Visa 13", "4000000000006", "400000***0006"},
		{"Amex 15", "340000000000009", "340000*****0009"},
		{"Short", "123456", "******"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskPAN(tt.pan)
			if result != tt.expected {
				t.Errorf("MaskPAN(%s) = %s, want %s", tt.pan, result, tt.expected)
			}
		})
	}
}

func BenchmarkGeneratePAN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GeneratePAN("400000", 16)
	}
}

func BenchmarkValidateLuhn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateLuhn("4000000000000002")
	}
}

func BenchmarkGenerateDeterministicCVC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateDeterministicCVC("4000000000000002", "12", "2027", "test-secret")
	}
}
