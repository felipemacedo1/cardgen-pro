package generator

import (
	"strconv"
)

// ValidateLuhn validates a PAN using the Luhn algorithm (mod-10)
// Luhn algorithm: used by all major card networks (ISO/IEC 7812)
// Returns true if the PAN is valid according to Luhn checksum
func ValidateLuhn(pan string) bool {
	if len(pan) < 13 || len(pan) > 19 {
		return false
	}

	sum := 0
	isSecond := false

	// Traverse from right to left
	for i := len(pan) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(pan[i]))
		if err != nil {
			return false
		}

		if isSecond {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isSecond = !isSecond
	}

	return sum%10 == 0
}

// CalculateLuhnCheckDigit calculates the Luhn check digit for a partial PAN
// Input: partial PAN without check digit (e.g., "40000000000000")
// Output: check digit (0-9)
func CalculateLuhnCheckDigit(partialPAN string) int {
	sum := 0
	isSecond := true // Since we're adding a digit, the last digit becomes second

	// Traverse from right to left
	for i := len(partialPAN) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(partialPAN[i]))
		if err != nil {
			return 0
		}

		if isSecond {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isSecond = !isSecond
	}

	// Calculate check digit: (10 - (sum % 10)) % 10
	checkDigit := (10 - (sum % 10)) % 10
	return checkDigit
}

// AppendLuhnCheckDigit appends the correct Luhn check digit to a partial PAN
func AppendLuhnCheckDigit(partialPAN string) string {
	checkDigit := CalculateLuhnCheckDigit(partialPAN)
	return partialPAN + strconv.Itoa(checkDigit)
}
