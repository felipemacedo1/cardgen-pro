package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/felipemacedo/cardgen-pro/internal/api"
	"github.com/felipemacedo/cardgen-pro/internal/generator"
	"github.com/felipemacedo/cardgen-pro/internal/iso"
	"github.com/felipemacedo/cardgen-pro/internal/models"
	"github.com/felipemacedo/cardgen-pro/pkg/transformer"
)

const version = "1.0.0"

const banner = `
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│   ██████╗ █████╗ ██████╗ ██████╗  ██████╗ ███████╗███╗   ██╗│
│  ██╔════╝██╔══██╗██╔══██╗██╔══██╗██╔════╝ ██╔════╝████╗  ██║│
│  ██║     ███████║██████╔╝██║  ██║██║  ███╗█████╗  ██╔██╗ ██║│
│  ██║     ██╔══██║██╔══██╗██║  ██║██║   ██║██╔══╝  ██║╚██╗██║│
│  ╚██████╗██║  ██║██║  ██║██████╔╝╚██████╔╝███████╗██║ ╚████║│
│   ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝  ╚═════╝ ╚══════╝╚═╝  ╚═══╝│
│                                                             │
│          Card Data & ISO-8583 Test Suite - PRO             │
│                    Version %s                           │
│                                                             │
│  ⚠️  WARNING: TEST/SANDBOX USE ONLY - NOT FOR PRODUCTION   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
`

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "generate":
		handleGenerate()
	case "transform":
		handleTransform()
	case "serve":
		handleServe()
	case "validate":
		handleValidate()
	case "scenarios":
		handleScenarios()
	case "version":
		fmt.Printf("cardgen-pro version %s\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(banner, version)
	fmt.Println("\nUsage: cardgen-pro <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  generate    Generate test card data")
	fmt.Println("  transform   Transform orders by injecting CVCs")
	fmt.Println("  serve       Start HTTP API server for fixtures")
	fmt.Println("  validate    Validate card numbers using Luhn")
	fmt.Println("  scenarios   List predefined test scenarios")
	fmt.Println("  version     Print version information")
	fmt.Println("  help        Show this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  cardgen-pro generate --bin 400000 --brand visa --count 10 --out cards.json")
	fmt.Println("  cardgen-pro transform --input orders.json --output orders_cvc.json")
	fmt.Println("  cardgen-pro serve --port 8080 --token my-dev-token")
	fmt.Println("  cardgen-pro validate 4000000000000002")
	fmt.Println("\nEnvironment Variables:")
	fmt.Println("  CARDGEN_SECRET    Secret key for deterministic CVC generation")
	fmt.Println("\nFor detailed help on a command, run: cardgen-pro <command> --help")
}

func handleGenerate() {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	
	bin := fs.String("bin", "", "BIN (Bank Identification Number) - first 6 digits")
	brand := fs.String("brand", "visa", "Card brand (visa, mastercard, amex)")
	count := fs.Int("count", 10, "Number of cards to generate")
	output := fs.String("out", "", "Output file path (JSON)")
	format := fs.String("format", "json", "Output format (json, ndjson, csv)")
	includeISO := fs.Bool("iso", false, "Include ISO-8583 fields")
	includeTrack2 := fs.Bool("track2", false, "Include Track2 data")
	secret := fs.String("secret", "", "Secret for CVC generation (or use CARDGEN_SECRET env)")
	
	fs.Parse(os.Args[2:])

	// Get secret from env if not provided
	secretValue := *secret
	if secretValue == "" {
		secretValue = os.Getenv("CARDGEN_SECRET")
	}

	if secretValue == "" {
		log.Println("⚠️  Warning: No secret provided. CVCs will not be generated.")
		log.Println("   Set CARDGEN_SECRET environment variable or use --secret flag")
	}

	// Generate cards
	opts := models.GenerateOptions{
		BIN:           *bin,
		Brand:         strings.ToLower(*brand),
		Count:         *count,
		Secret:        secretValue,
		IncludeISO:    *includeISO,
		IncludeTrack2: *includeTrack2,
	}

	cards := []*models.Card{}
	for i := 0; i < *count; i++ {
		card, err := generator.GenerateCard(opts)
		if err != nil {
			log.Fatalf("Failed to generate card: %v", err)
		}

		// Add ISO fields if requested
		if *includeISO {
			isoFields := iso.GenerateISO8583Fields(card, 10000, "986")
			card.ISOFields = isoFields
		}

		cards = append(cards, card)
	}

	// Output
	if *output == "" {
		// Print to stdout
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(cards); err != nil {
			log.Fatalf("Failed to encode cards: %v", err)
		}
	} else {
		// Write to file
		var err error
		switch *format {
		case "json":
			err = transformer.WriteCardsJSON(*output, cards)
		case "ndjson":
			err = transformer.WriteCardsNDJSON(*output, cards)
		case "csv":
			err = transformer.WriteCardsCSV(*output, cards)
		default:
			log.Fatalf("Unknown format: %s", *format)
		}

		if err != nil {
			log.Fatalf("Failed to write output: %v", err)
		}

		log.Printf("✓ Generated %d cards and saved to %s", len(cards), *output)
		
		// Show sample
		if len(cards) > 0 {
			sample := cards[0]
			log.Printf("\nSample card:")
			log.Printf("  PAN (masked): %s", sample.MaskedPAN)
			log.Printf("  Brand: %s", sample.Brand)
			log.Printf("  Expiry: %02d/%d", sample.ExpiryMonth, sample.ExpiryYear)
			if sample.CVC != "" {
				log.Printf("  CVC: %s", sample.CVC)
			}
		}
	}
}

func handleTransform() {
	fs := flag.NewFlagSet("transform", flag.ExitOnError)
	
	input := fs.String("input", "", "Input file path (JSON)")
	output := fs.String("output", "", "Output file path (JSON)")
	secret := fs.String("secret", "", "Secret for CVC generation (or use CARDGEN_SECRET env)")
	
	fs.Parse(os.Args[2:])

	if *input == "" {
		log.Fatal("Error: --input is required")
	}
	if *output == "" {
		log.Fatal("Error: --output is required")
	}

	// Get secret from env if not provided
	secretValue := *secret
	if secretValue == "" {
		secretValue = os.Getenv("CARDGEN_SECRET")
	}

	if secretValue == "" {
		log.Fatal("Error: Secret is required. Set CARDGEN_SECRET or use --secret flag")
	}

	// Transform
	opts := models.TransformOptions{
		InputPath:  *input,
		OutputPath: *output,
		Secret:     secretValue,
	}

	if err := transformer.TransformOrders(opts); err != nil {
		log.Fatalf("Failed to transform orders: %v", err)
	}

	log.Printf("✓ Transformed orders and saved to %s", *output)
}

func handleServe() {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	
	port := fs.Int("port", 8080, "HTTP server port")
	token := fs.String("token", "", "Authentication token (required)")
	
	fs.Parse(os.Args[2:])

	if *token == "" {
		log.Fatal("Error: --token is required for API server")
	}

	log.Printf("Starting cardgen-pro API server v%s", version)
	log.Println("⚠️  WARNING: This server is for TEST/SANDBOX use only")
	log.Printf("\nAuthentication: Bearer %s\n", *token)

	server := api.NewServer(*token, *port)
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func handleValidate() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: cardgen-pro validate <PAN>")
		os.Exit(1)
	}

	pan := os.Args[2]

	if generator.ValidateLuhn(pan) {
		fmt.Printf("✓ Valid: %s is a valid PAN (Luhn check passed)\n", generator.MaskPAN(pan))
	} else {
		fmt.Printf("✗ Invalid: %s failed Luhn check\n", generator.MaskPAN(pan))
		os.Exit(1)
	}
}

func handleScenarios() {
	scenarios := api.GetScenarios()

	fmt.Printf("\n=== Test Scenarios (%d total) ===\n\n", len(scenarios))

	for _, scenario := range scenarios {
		fmt.Printf("ID: %s\n", scenario.ID)
		fmt.Printf("  Name: %s\n", scenario.Name)
		fmt.Printf("  Description: %s\n", scenario.Description)
		fmt.Printf("  Response: [%s] %s\n", scenario.ResponseCode, scenario.ResponseText)
		fmt.Printf("  Amount: %d %s\n", scenario.Amount, scenario.Currency)
		fmt.Printf("  Expected: %s\n", scenario.ExpectedOutcome)
		if len(scenario.Metadata) > 0 {
			fmt.Printf("  Metadata: %v\n", scenario.Metadata)
		}
		fmt.Println()
	}
}
