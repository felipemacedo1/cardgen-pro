package transformer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/felipemacedo/cardgen-pro/internal/generator"
	"github.com/felipemacedo/cardgen-pro/internal/models"
)

// TransformOrders reads orders from input file, injects CVCs, and writes to output
// Supports JSON and NDJSON formats
func TransformOrders(opts models.TransformOptions) error {
	if opts.Secret == "" {
		return fmt.Errorf("secret is required for CVC generation")
	}

	// Read input file
	orders, err := ReadOrders(opts.InputPath)
	if err != nil {
		return fmt.Errorf("failed to read orders: %w", err)
	}

	// Transform: inject CVCs
	for i := range orders {
		order := &orders[i]
		
		// Skip if already has CVC
		if order.CVC != "" {
			continue
		}

		// Generate deterministic CVC
		cvc, err := generator.GenerateDeterministicCVC(
			order.PAN,
			fmt.Sprintf("%02d", order.ExpiryMonth),
			fmt.Sprintf("%d", order.ExpiryYear),
			opts.Secret,
		)
		if err != nil {
			return fmt.Errorf("failed to generate CVC for order %s: %w", order.ID, err)
		}

		order.CVC = cvc
	}

	// Write output file
	if err := WriteOrders(opts.OutputPath, orders); err != nil {
		return fmt.Errorf("failed to write orders: %w", err)
	}

	return nil
}

// ReadOrders reads orders from a JSON or NDJSON file
func ReadOrders(path string) ([]models.Order, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Try parsing as JSON array first
	var orders []models.Order
	if err := json.Unmarshal(data, &orders); err == nil {
		return orders, nil
	}

	// Try parsing as NDJSON (one JSON per line)
	decoder := json.NewDecoder(file)
	file.Seek(0, 0) // Reset file pointer

	orders = []models.Order{}
	for {
		var order models.Order
		if err := decoder.Decode(&order); err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to parse NDJSON: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// WriteOrders writes orders to a JSON file
func WriteOrders(path string, orders []models.Order) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	return encoder.Encode(orders)
}

// WriteCardsJSON writes cards to a JSON file (pretty-printed)
func WriteCardsJSON(path string, cards []*models.Card) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	return encoder.Encode(cards)
}

// WriteCardsNDJSON writes cards to an NDJSON file (one JSON per line)
func WriteCardsNDJSON(path string, cards []*models.Card) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	
	for _, card := range cards {
		if err := encoder.Encode(card); err != nil {
			return err
		}
	}

	return nil
}

// WriteCardsCSV writes cards to a CSV file
func WriteCardsCSV(path string, cards []*models.Card) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"PAN", "MaskedPAN", "Brand", "ExpiryMonth", "ExpiryYear", "CVC", "Track2"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, card := range cards {
		row := []string{
			card.PAN,
			card.MaskedPAN,
			card.Brand,
			fmt.Sprintf("%d", card.ExpiryMonth),
			fmt.Sprintf("%d", card.ExpiryYear),
			card.CVC,
			card.Track2,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
