# Technical Specification

## Overview

This document provides detailed technical specifications for the **cardgen-pro** Card Data & ISO-8583 Test Suite, including algorithms, data formats, and implementation details.

## Table of Contents

1. [Luhn Algorithm](#luhn-algorithm)
2. [PAN Generation](#pan-generation)
3. [Deterministic CVC Generation](#deterministic-cvc-generation)
4. [Track2 Format](#track2-format)
5. [ISO-8583 Fields](#iso-8583-fields)
6. [Response Codes](#response-codes)
7. [Data Formats](#data-formats)

## Luhn Algorithm

### Purpose

The Luhn algorithm (also known as mod-10 checksum) is used to validate credit card numbers and detect simple data entry errors.

### Algorithm Steps

1. Starting from the **rightmost** digit (check digit), double every **second** digit
2. If doubling results in a two-digit number, subtract 9 (or sum the digits)
3. Sum all the digits
4. If the total modulo 10 equals 0, the number is valid

### Implementation

```go
func ValidateLuhn(pan string) bool {
    sum := 0
    isSecond := false
    
    for i := len(pan) - 1; i >= 0; i-- {
        digit := int(pan[i] - '0')
        
        if isSecond {
            digit *= 2
            if digit > 9 {
                digit -= 9
            }
        }
        
        sum += digit
        isSecond = !isSecond
    }
    
    return sum % 10 == 0
}
```

### Example

Validate PAN: `4000000000000002`

```
Positions (right to left):
1: 2 (no double) = 2
2: 0 (double) = 0
...
16: 4 (double) = 8

Sum = 2 + 0 + 0 + 0 + 0 + 0 + 0 + 0 + 0 + 0 + 0 + 0 + 0 + 0 + 0 + 8 = 10
10 % 10 = 0 ✓ Valid
```

## PAN Generation

### Structure

```
[IIN/BIN][Account Number][Check Digit]
  (6)      (variable)        (1)
```

### Card Brand Ranges

| Brand | BIN Range | Length | CVC Length |
|-------|-----------|--------|------------|
| Visa | 4xxxxx | 13, 16, 19 | 3 |
| Mastercard | 51-55xxxx, 2221-2720xx | 16 | 3 |
| American Express | 34xxxx, 37xxxx | 15 | 4 |

## Deterministic CVC Generation

### Algorithm: HMAC-SHA256

**Payload:**
```
Payload = BIN6 | LAST4 | EXPMM | EXPYYYY
Example: "400000|0002|12|2027"
```

**Why HMAC-SHA256?**
- Determinism: Same input → same output
- Security: Cryptographically strong, one-way
- Secret-based: Prevents trivial generation

See full specification in source code.

## Track2 Format

```
[PAN]=[Expiry][Service Code][Discretionary Data]
Example: 4000000000000002=27122011234
```

## ISO-8583 Fields

### Common Fields

| Field | Name | Example |
|-------|------|---------|
| 2 | PAN | "4000000000000002" |
| 3 | Processing Code | "000000" |
| 4 | Amount | "000000010000" |
| 11 | STAN | "123456" |
| 14 | Expiry | "2712" |
| 39 | Response Code | "00" |
| 49 | Currency | "986" |

## Response Codes

| Code | Description |
|------|-------------|
| 00 | Approved |
| 05 | Do not honor |
| 51 | Insufficient funds |
| 54 | Expired card |
| 91 | Issuer unavailable |

See full specification for details.
