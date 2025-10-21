package generator

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/felipemacedo/cardgen-pro/internal/models"
)

// CardBrands defines well-known card brand configurations
// These are TEST BINs and ranges - DO NOT USE IN PRODUCTION
var CardBrands = map[string]models.CardBrand{
	"visa": {
		Name: "Visa",
		BINRanges: []models.BINRange{
			{Start: "400000", End: "499999", Length: 16},
		},
		PANLength:   []int{13, 16, 19},
		CVCLength:   3,
		ServiceCode: "201", // Common test service code
	},
	"mastercard": {
		Name: "Mastercard",
		BINRanges: []models.BINRange{
			{Start: "510000", End: "559999", Length: 16},
			{Start: "222100", End: "272099", Length: 16},
		},
		PANLength:   []int{16},
		CVCLength:   3,
		ServiceCode: "201",
	},
	"amex": {
		Name: "American Express",
		BINRanges: []models.BINRange{
			{Start: "340000", End: "349999", Length: 15},
			{Start: "370000", End: "379999", Length: 15},
		},
		PANLength:   []int{15},
		CVCLength:   4,
		ServiceCode: "201",
	},
}

// GeneratePAN generates a valid PAN using Luhn algorithm
// BIN: Bank Identification Number (first 6 digits)
// length: total PAN length (13-19 for most cards, 15 for Amex)
func GeneratePAN(bin string, length int) (string, error) {
	if len(bin) < 6 {
		return "", fmt.Errorf("BIN must be at least 6 digits")
	}

	if length < 13 || length > 19 {
		return "", fmt.Errorf("PAN length must be between 13 and 19")
	}

	// Calculate how many random digits we need (excluding check digit)
	randomDigitsNeeded := length - len(bin) - 1

	if randomDigitsNeeded < 0 {
		return "", fmt.Errorf("BIN too long for requested PAN length")
	}

	// Generate random middle digits
	randomPart := generateRandomDigits(randomDigitsNeeded)

	// Construct partial PAN (BIN + random digits)
	partialPAN := bin + randomPart

	// Append Luhn check digit
	fullPAN := AppendLuhnCheckDigit(partialPAN)

	return fullPAN, nil
}

// generateRandomDigits generates n random digits
func generateRandomDigits(n int) string {
	if n <= 0 {
		return ""
	}

	digits := make([]byte, n)
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		digits[i] = byte('0' + num.Int64())
	}
	return string(digits)
}

// GenerateExpiry generates a plausible expiry date
// Returns month (1-12) and year (current year + 1 to current year + 5)
func GenerateExpiry() (month int, year int) {
	now := time.Now()
	
	// Random month
	monthBig, _ := rand.Int(rand.Reader, big.NewInt(12))
	month = int(monthBig.Int64()) + 1

	// Random year offset (1-5 years in the future)
	yearOffsetBig, _ := rand.Int(rand.Reader, big.NewInt(5))
	year = now.Year() + int(yearOffsetBig.Int64()) + 1

	return month, year
}

// GenerateDeterministicCVC generates a deterministic CVC using HMAC-SHA256
// 
// DESIGN RATIONALE:
// - HMAC-SHA256 provides cryptographic strength and determinism
// - Same input always produces same output (reproducible for tests)
// - Secret key prevents reverse engineering
// - Payload includes all relevant card data to ensure uniqueness
//
// SECURITY:
// - Secret MUST be provided via environment variable or secure vault
// - Never hardcode secrets in source code
// - CVC is derived, not a real CVV/CVV2 from card issuer
// - FOR TEST/SANDBOX USE ONLY
func GenerateDeterministicCVC(pan, expMonth, expYear, secret string) (string, error) {
	if secret == "" {
		return "", fmt.Errorf("secret is required for CVC generation")
	}

	// Extract BIN6 and last4 for payload
	bin6 := pan
	if len(pan) >= 6 {
		bin6 = pan[:6]
	}
	
	last4 := pan
	if len(pan) >= 4 {
		last4 = pan[len(pan)-4:]
	}

	// Construct payload: BIN6|LAST4|EXPMM|EXPYYYY
	payload := fmt.Sprintf("%s|%s|%s|%s", bin6, last4, expMonth, expYear)

	// Generate HMAC-SHA256
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	hash := h.Sum(nil)

	// Convert hash to hex and take first characters
	hexHash := hex.EncodeToString(hash)

	// Determine CVC length based on PAN (Amex = 4, others = 3)
	cvcLength := 3
	if len(pan) == 15 && (strings.HasPrefix(pan, "34") || strings.HasPrefix(pan, "37")) {
		cvcLength = 4
	}

	// Extract numeric digits from hex hash
	numericDigits := ""
	for _, char := range hexHash {
		if char >= '0' && char <= '9' {
			numericDigits += string(char)
			if len(numericDigits) == cvcLength {
				break
			}
		}
	}

	// Fallback: use first N digits converted from hex
	if len(numericDigits) < cvcLength {
		for i := 0; i < cvcLength; i++ {
			digit := int(hash[i]) % 10
			numericDigits += strconv.Itoa(digit)
		}
	}

	return numericDigits[:cvcLength], nil
}

// GenerateTrack2 generates a Track2-like string
// Format: PAN=YYMM<ServiceCode><DiscretionaryData>
//
// DESIGN RATIONALE:
// - Track2 is used in magnetic stripe and chip cards
// - Contains PAN, expiry, and service code
// - Discretionary data is random (used by issuers for various purposes)
// - This is a SIMPLIFIED version for testing; real Track2 has more fields
func GenerateTrack2(pan string, month, year int, serviceCode string) string {
	// Format expiry as YYMM
	yy := year % 100
	expiry := fmt.Sprintf("%02d%02d", yy, month)

	// Generate random discretionary data (3-5 digits)
	discretionaryLength := 4
	discretionary := generateRandomDigits(discretionaryLength)

	// Track2 format: PAN=YYMM<ServiceCode><Discretionary>
	track2 := fmt.Sprintf("%s=%s%s%s", pan, expiry, serviceCode, discretionary)

	return track2
}

// MaskPAN masks a PAN for safe logging/display
// Format: first6 **** last4 (e.g., "400000******1234")
func MaskPAN(pan string) string {
	if len(pan) < 10 {
		return strings.Repeat("*", len(pan))
	}

	first6 := pan[:6]
	last4 := pan[len(pan)-4:]
	masked := first6 + strings.Repeat("*", len(pan)-10) + last4

	return masked
}

// GenerateCard generates a complete card with all data
func GenerateCard(opts models.GenerateOptions) (*models.Card, error) {
	// Determine brand config
	brandConfig, ok := CardBrands[strings.ToLower(opts.Brand)]
	if !ok {
		return nil, fmt.Errorf("unknown brand: %s", opts.Brand)
	}

	// Use provided BIN or generate from brand
	bin := opts.BIN
	if bin == "" {
		// Use first BIN range of brand
		bin = brandConfig.BINRanges[0].Start
	}

	// Determine PAN length
	panLength := brandConfig.PANLength[0]
	if len(brandConfig.BINRanges) > 0 {
		panLength = brandConfig.BINRanges[0].Length
	}

	// Generate PAN
	pan, err := GeneratePAN(bin, panLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PAN: %w", err)
	}

	// Generate expiry
	month, year := GenerateExpiry()

	// Generate CVC if secret provided
	var cvc string
	if opts.Secret != "" {
		cvc, err = GenerateDeterministicCVC(
			pan,
			fmt.Sprintf("%02d", month),
			fmt.Sprintf("%d", year),
			opts.Secret,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to generate CVC: %w", err)
		}
	}

	card := &models.Card{
		PAN:         pan,
		MaskedPAN:   MaskPAN(pan),
		Brand:       brandConfig.Name,
		ExpiryMonth: month,
		ExpiryYear:  year,
		CVC:         cvc,
		GeneratedAt: time.Now(),
		Metadata:    opts.Metadata,
	}

	// Generate Track2 if requested
	if opts.IncludeTrack2 {
		card.Track2 = GenerateTrack2(pan, month, year, brandConfig.ServiceCode)
	}

	return card, nil
}
