package api

// Scenario represents a test scenario with expected behavior
type Scenario struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	ResponseCode    string            `json:"response_code"`
	ResponseText    string            `json:"response_text"`
	Amount          int64             `json:"amount"`
	Currency        string            `json:"currency"`
	CardBrand       string            `json:"card_brand"`
	ExpectedOutcome string            `json:"expected_outcome"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// GetScenarios returns predefined test scenarios
func GetScenarios() []Scenario {
	return []Scenario{
		{
			ID:              "success_auth",
			Name:            "Successful Authorization",
			Description:     "Standard approved transaction",
			ResponseCode:    "00",
			ResponseText:    "Approved",
			Amount:          10000,
			Currency:        "986",
			CardBrand:       "visa",
			ExpectedOutcome: "Transaction approved, auth code generated",
		},
		{
			ID:              "declined_generic",
			Name:            "Generic Decline",
			Description:     "Transaction declined without specific reason",
			ResponseCode:    "05",
			ResponseText:    "Do not honor",
			Amount:          50000,
			Currency:        "986",
			CardBrand:       "mastercard",
			ExpectedOutcome: "Transaction declined by issuer",
		},
		{
			ID:              "insufficient_funds",
			Name:            "Insufficient Funds",
			Description:     "Card has insufficient balance",
			ResponseCode:    "51",
			ResponseText:    "Insufficient funds",
			Amount:          100000,
			Currency:        "986",
			CardBrand:       "visa",
			ExpectedOutcome: "Decline due to insufficient balance",
		},
		{
			ID:              "3ds_required",
			Name:            "3DS Authentication Required",
			Description:     "Transaction requires 3D Secure authentication",
			ResponseCode:    "00",
			ResponseText:    "Approved",
			Amount:          25000,
			Currency:        "986",
			CardBrand:       "visa",
			ExpectedOutcome: "Redirect to 3DS flow before completion",
			Metadata: map[string]string{
				"3ds_required": "true",
				"3ds_version":  "2.0",
			},
		},
		{
			ID:              "auth_only",
			Name:            "Authorization Only (Pre-Auth)",
			Description:     "Amount held but not captured",
			ResponseCode:    "00",
			ResponseText:    "Approved",
			Amount:          15000,
			Currency:        "986",
			CardBrand:       "mastercard",
			ExpectedOutcome: "Authorization successful, requires capture",
			Metadata: map[string]string{
				"type": "pre_auth",
			},
		},
		{
			ID:              "captured",
			Name:            "Captured Transaction",
			Description:     "Previously authorized transaction now captured",
			ResponseCode:    "00",
			ResponseText:    "Approved",
			Amount:          15000,
			Currency:        "986",
			CardBrand:       "mastercard",
			ExpectedOutcome: "Funds captured successfully",
			Metadata: map[string]string{
				"type":         "capture",
				"original_ref": "AUTH123456",
			},
		},
		{
			ID:              "refunded_partial",
			Name:            "Partial Refund",
			Description:     "Part of transaction amount refunded",
			ResponseCode:    "00",
			ResponseText:    "Approved",
			Amount:          5000,
			Currency:        "986",
			CardBrand:       "visa",
			ExpectedOutcome: "Partial refund processed",
			Metadata: map[string]string{
				"type":           "refund",
				"original_amount": "15000",
				"refund_amount":   "5000",
			},
		},
		{
			ID:              "chargeback_open",
			Name:            "Chargeback Initiated",
			Description:     "Customer disputed the transaction",
			ResponseCode:    "00",
			ResponseText:    "Approved",
			Amount:          20000,
			Currency:        "986",
			CardBrand:       "mastercard",
			ExpectedOutcome: "Chargeback case opened, merchant response required",
			Metadata: map[string]string{
				"type":           "chargeback",
				"reason_code":    "4853",
				"dispute_amount": "20000",
			},
		},
		{
			ID:              "pix_paid",
			Name:            "PIX Payment Successful",
			Description:     "Brazilian instant payment completed",
			ResponseCode:    "00",
			ResponseText:    "Approved",
			Amount:          35000,
			Currency:        "986",
			CardBrand:       "pix",
			ExpectedOutcome: "PIX payment confirmed instantly",
			Metadata: map[string]string{
				"payment_method": "pix",
				"pix_key":        "user@example.com",
			},
		},
		{
			ID:              "boleto_pending",
			Name:            "Boleto Payment Pending",
			Description:     "Brazilian boleto generated, awaiting payment",
			ResponseCode:    "00",
			ResponseText:    "Approved",
			Amount:          45000,
			Currency:        "986",
			CardBrand:       "boleto",
			ExpectedOutcome: "Boleto generated, pending customer payment",
			Metadata: map[string]string{
				"payment_method": "boleto",
				"due_date":       "2025-11-01",
				"barcode":        "34191.79001 01043.510047 91020.150008 1 96610000012345",
			},
		},
		{
			ID:              "subscription_recurring",
			Name:            "Recurring Subscription Payment",
			Description:     "Monthly subscription charge",
			ResponseCode:    "00",
			ResponseText:    "Approved",
			Amount:          4990,
			Currency:        "986",
			CardBrand:       "visa",
			ExpectedOutcome: "Recurring payment successful",
			Metadata: map[string]string{
				"type":              "subscription",
				"subscription_id":   "SUB_12345",
				"billing_cycle":     "monthly",
				"next_billing_date": "2025-11-20",
			},
		},
		{
			ID:              "tokenized_payment",
			Name:            "Tokenized Card Payment",
			Description:     "Payment using stored card token",
			ResponseCode:    "00",
			ResponseText:    "Approved",
			Amount:          12000,
			Currency:        "986",
			CardBrand:       "mastercard",
			ExpectedOutcome: "Token-based payment successful",
			Metadata: map[string]string{
				"type":       "token_payment",
				"token_id":   "tok_1A2B3C4D5E6F",
				"token_type": "card_on_file",
			},
		},
	}
}
