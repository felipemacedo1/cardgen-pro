package test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/felipemacedo/cardgen-pro/internal/generator"
	"github.com/felipemacedo/cardgen-pro/internal/models"
	"github.com/felipemacedo/cardgen-pro/pkg/transformer"
)

func TestIntegrationGenerateAndTransform(t *testing.T) {
	secret := "integration-test-secret"

	// Step 1: Generate orders without CVC
	orders := []models.Order{
		{
			ID:          "ORD001",
			PAN:         "4000000000000002",
			ExpiryMonth: 12,
			ExpiryYear:  2027,
			Amount:      10000,
			Currency:    "986",
		},
		{
			ID:          "ORD002",
			PAN:         "5100000000000016",
			ExpiryMonth: 6,
			ExpiryYear:  2026,
			Amount:      25000,
			Currency:    "986",
		},
	}

	// Step 2: Write orders to temp file
	inputPath := "/tmp/cardgen_test_orders_input.json"
	outputPath := "/tmp/cardgen_test_orders_output.json"
	
	defer os.Remove(inputPath)
	defer os.Remove(outputPath)

	if err := writeOrdersToFile(inputPath, orders); err != nil {
		t.Fatalf("Failed to write input orders: %v", err)
	}

	// Step 3: Transform (inject CVCs)
	opts := models.TransformOptions{
		InputPath:  inputPath,
		OutputPath: outputPath,
		Secret:     secret,
	}

	if err := transformer.TransformOrders(opts); err != nil {
		t.Fatalf("Transform failed: %v", err)
	}

	// Step 4: Read transformed orders
	transformedOrders, err := readOrdersFromFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output orders: %v", err)
	}

	// Step 5: Verify CVCs were injected
	if len(transformedOrders) != len(orders) {
		t.Errorf("Order count mismatch: got %d, want %d", len(transformedOrders), len(orders))
	}

	for _, order := range transformedOrders {
		if order.CVC == "" {
			t.Errorf("Order %s missing CVC", order.ID)
		}

		// Verify CVC is correct length
		expectedLen := 3
		if len(order.PAN) == 15 {
			expectedLen = 4
		}
		if len(order.CVC) != expectedLen {
			t.Errorf("Order %s CVC length = %d, want %d", order.ID, len(order.CVC), expectedLen)
		}

		// Step 6: Recalculate CVC and verify determinism
		recalculatedCVC, err := generator.GenerateDeterministicCVC(
			order.PAN,
			fmt.Sprintf("%02d", order.ExpiryMonth),
			fmt.Sprintf("%d", order.ExpiryYear),
			secret,
		)
		if err != nil {
			t.Errorf("Failed to recalculate CVC for order %s: %v", order.ID, err)
			continue
		}

		if order.CVC != recalculatedCVC {
			t.Errorf("Order %s CVC mismatch: got %s, recalculated %s", order.ID, order.CVC, recalculatedCVC)
		}

		t.Logf("✓ Order %s: CVC=%s (verified deterministic)", order.ID, order.CVC)
	}
}

func TestIntegrationGenerateCards(t *testing.T) {
	secret := "integration-test-secret"

	brands := []string{"visa", "mastercard", "amex"}

	for _, brand := range brands {
		t.Run(brand, func(t *testing.T) {
			opts := models.GenerateOptions{
				Brand:         brand,
				Count:         5,
				Secret:        secret,
				IncludeISO:    true,
				IncludeTrack2: true,
			}

			for i := 0; i < opts.Count; i++ {
				card, err := generator.GenerateCard(opts)
				if err != nil {
					t.Fatalf("Failed to generate %s card: %v", brand, err)
				}

				// Verify brand
				if card.Brand == "" {
					t.Error("Card brand is empty")
				}

				// Verify PAN is valid
				if !generator.ValidateLuhn(card.PAN) {
					t.Errorf("Generated PAN %s failed Luhn check", card.MaskedPAN)
				}

				// Verify CVC
				if card.CVC == "" {
					t.Error("Card missing CVC")
				}

				expectedCVCLen := 3
				if brand == "amex" {
					expectedCVCLen = 4
				}
				if len(card.CVC) != expectedCVCLen {
					t.Errorf("CVC length = %d, want %d", len(card.CVC), expectedCVCLen)
				}

				// Verify Track2
				if card.Track2 == "" {
					t.Error("Card missing Track2")
				}

				// Verify expiry is in future
				if card.ExpiryYear < 2025 {
					t.Errorf("ExpiryYear %d is in the past", card.ExpiryYear)
				}

				t.Logf("✓ Generated %s card: %s", brand, card.MaskedPAN)
			}
		})
	}
}

func TestIntegrationCVCDeterminism(t *testing.T) {
	secret := "determinism-test-secret"
	pan := "4000000000000002"
	month := "12"
	year := "2027"

	// Generate CVC multiple times
	iterations := 10
	firstCVC := ""

	for i := 0; i < iterations; i++ {
		cvc, err := generator.GenerateDeterministicCVC(pan, month, year, secret)
		if err != nil {
			t.Fatalf("Failed to generate CVC: %v", err)
		}

		if i == 0 {
			firstCVC = cvc
		} else {
			if cvc != firstCVC {
				t.Errorf("CVC not deterministic: iteration %d got %s, want %s", i, cvc, firstCVC)
			}
		}
	}

	t.Logf("✓ CVC determinism verified over %d iterations: %s", iterations, firstCVC)

	// Verify different secrets produce different CVCs
	differentSecret := "different-secret"
	differentCVC, err := generator.GenerateDeterministicCVC(pan, month, year, differentSecret)
	if err != nil {
		t.Fatalf("Failed to generate CVC with different secret: %v", err)
	}

	if differentCVC == firstCVC {
		t.Error("Different secrets produced same CVC")
	}

	t.Logf("✓ Different secret produced different CVC: %s vs %s", firstCVC, differentCVC)
}

func TestIntegrationOutputFormats(t *testing.T) {
	secret := "format-test-secret"
	
	opts := models.GenerateOptions{
		Brand:  "visa",
		Count:  3,
		Secret: secret,
	}

	cards := []*models.Card{}
	for i := 0; i < opts.Count; i++ {
		card, err := generator.GenerateCard(opts)
		if err != nil {
			t.Fatalf("Failed to generate card: %v", err)
		}
		cards = append(cards, card)
	}

	// Test JSON output
	t.Run("JSON", func(t *testing.T) {
		path := "/tmp/cardgen_test_cards.json"
		defer os.Remove(path)

		if err := transformer.WriteCardsJSON(path, cards); err != nil {
			t.Fatalf("Failed to write JSON: %v", err)
		}

		// Verify file exists and is readable
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("Failed to read JSON: %v", err)
		}

		var readCards []*models.Card
		if err := json.Unmarshal(data, &readCards); err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		if len(readCards) != len(cards) {
			t.Errorf("Card count mismatch: got %d, want %d", len(readCards), len(cards))
		}

		t.Logf("✓ JSON format verified: %d cards", len(readCards))
	})

	// Test NDJSON output
	t.Run("NDJSON", func(t *testing.T) {
		path := "/tmp/cardgen_test_cards.ndjson"
		defer os.Remove(path)

		if err := transformer.WriteCardsNDJSON(path, cards); err != nil {
			t.Fatalf("Failed to write NDJSON: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("NDJSON file not created: %v", err)
		}

		t.Logf("✓ NDJSON format verified")
	})

	// Test CSV output
	t.Run("CSV", func(t *testing.T) {
		path := "/tmp/cardgen_test_cards.csv"
		defer os.Remove(path)

		if err := transformer.WriteCardsCSV(path, cards); err != nil {
			t.Fatalf("Failed to write CSV: %v", err)
		}

		// Verify file exists
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("Failed to read CSV: %v", err)
		}

		// Verify has header
		if len(data) == 0 {
			t.Fatal("CSV file is empty")
		}

		t.Logf("✓ CSV format verified")
	})
}

// Helper functions

func writeOrdersToFile(path string, orders []models.Order) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(orders)
}

func readOrdersFromFile(path string) ([]models.Order, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var orders []models.Order
	if err := json.Unmarshal(data, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
