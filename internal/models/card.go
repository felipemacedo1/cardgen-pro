package models

import "time"

// Card represents a generated payment card with all associated data
type Card struct {
	PAN          string            `json:"pan"`
	MaskedPAN    string            `json:"masked_pan"`
	Brand        string            `json:"brand"`
	ExpiryMonth  int               `json:"expiry_month"`
	ExpiryYear   int               `json:"expiry_year"`
	CVC          string            `json:"cvc,omitempty"`
	Track2       string            `json:"track2,omitempty"`
	ISOFields    map[string]string `json:"iso_fields,omitempty"`
	GeneratedAt  time.Time         `json:"generated_at"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// CardBrand represents a card brand configuration
type CardBrand struct {
	Name        string
	BINRanges   []BINRange
	PANLength   []int
	CVCLength   int
	ServiceCode string
}

// BINRange represents a valid BIN range for a brand
type BINRange struct {
	Start  string
	End    string
	Length int
}

// GenerateOptions contains options for card generation
type GenerateOptions struct {
	BIN         string
	Brand       string
	Count       int
	Secret      string
	IncludeISO  bool
	IncludeTrack2 bool
	Metadata    map[string]string
}

// TransformOptions contains options for transforming orders with CVCs
type TransformOptions struct {
	InputPath  string
	OutputPath string
	Secret     string
}

// Order represents a payment order/transaction
type Order struct {
	ID          string            `json:"id"`
	PAN         string            `json:"pan"`
	ExpiryMonth int               `json:"expiry_month"`
	ExpiryYear  int               `json:"expiry_year"`
	CVC         string            `json:"cvc,omitempty"`
	Amount      int64             `json:"amount"`
	Currency    string            `json:"currency"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}
