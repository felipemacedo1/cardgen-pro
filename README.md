# cardgen-pro

> **⚠️ CRITICAL WARNING: TEST/SANDBOX USE ONLY — NEVER USE IN PRODUCTION**
> 
> This tool generates **synthetic test data** for payment card testing and ISO-8583 message simulation. Generated PANs are NOT real card numbers and MUST NOT be used in production payment networks.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue)](https://golang.org)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen)](#testing)
[![Coverage](https://img.shields.io/badge/coverage-85%25-green)](#testing)

**cardgen-pro** is a professional-grade Card Data & ISO-8583 Test Suite written in Go. It generates valid test payment card numbers, deterministic CVCs, Track2 data, and ISO-8583 message fields for QA, development, and infrastructure testing in sandbox environments.

## 🎯 Features

- ✅ **Valid PANs**: Generate Luhn-valid PANs for Visa, Mastercard, American Express
- 🔐 **Deterministic CVCs**: HMAC-SHA256-based CVC generation (reproducible, secure)
- 📊 **ISO-8583 Fields**: Generate common authorization message fields
- 🎫 **Track2 Data**: Magnetic stripe-like track data simulation
- 📦 **Multiple Formats**: JSON, NDJSON, CSV output
- 🔄 **Transform Mode**: Inject CVCs into existing order files
- 🌐 **HTTP API**: Optional sandbox fixture server with auth + rate limiting
- 🧪 **12 Test Scenarios**: Pre-built fixtures for common payment flows
- ✨ **CLI & Library**: Use as command-line tool or Go package

## 🚀 Quick Start

### Installation

```bash
# From source
git clone https://github.com/felipemacedo1/cardgen-pro.git
cd cardgen-pro
go build -o cardgen-pro ./cmd/cardgen-pro

# Or install directly
go install github.com/felipemacedo/cardgen-pro/cmd/cardgen-pro@latest
```

### Basic Usage

```bash
# Set your secret key (REQUIRED for CVC generation)
export CARDGEN_SECRET="your-dev-secret-key"

# Generate 10 Visa cards
cardgen-pro generate --bin 400000 --brand visa --count 10 --out cards.json

# Generate with ISO-8583 fields and Track2
cardgen-pro generate --brand mastercard --count 5 --iso --track2 --out cards.json

# Transform orders (inject CVCs)
cardgen-pro transform --input orders.json --output orders_with_cvc.json

# Start API server
cardgen-pro serve --port 8080 --token my-dev-token

# Validate a PAN
cardgen-pro validate 4000000000000002

# List test scenarios
cardgen-pro scenarios
```

## 📖 Documentation

### Generate Command

Generate test card data with Luhn-valid PANs and deterministic CVCs.

```bash
cardgen-pro generate [options]
```

**Options:**
- `--bin <string>`: BIN (Bank Identification Number) - first 6 digits (optional)
- `--brand <string>`: Card brand: `visa`, `mastercard`, `amex` (default: `visa`)
- `--count <int>`: Number of cards to generate (default: 10)
- `--out <path>`: Output file path (prints to stdout if not specified)
- `--format <string>`: Output format: `json`, `ndjson`, `csv` (default: `json`)
- `--iso`: Include ISO-8583 fields
- `--track2`: Include Track2 data
- `--secret <string>`: Secret for CVC generation (or use `CARDGEN_SECRET` env var)

**Example Output:**

```json
[
  {
    "pan": "4000000000000002",
    "masked_pan": "400000******0002",
    "brand": "Visa",
    "expiry_month": 12,
    "expiry_year": 2027,
    "cvc": "842",
    "track2": "4000000000000002=27122011234",
    "generated_at": "2025-10-21T00:00:00Z"
  }
]
```

### Transform Command

Inject deterministic CVCs into existing order files.

```bash
cardgen-pro transform --input orders.json --output orders_with_cvc.json
```

**Input Format:**

```json
[
  {
    "id": "ORD001",
    "pan": "4000000000000002",
    "expiry_month": 12,
    "expiry_year": 2027,
    "amount": 10000,
    "currency": "986"
  }
]
```

**Output:** Same structure with `cvc` field populated.

### Serve Command

Start an HTTP API server for fixture serving (sandbox only).

```bash
cardgen-pro serve --port 8080 --token <auth-token>
```

**Endpoints:**
- `GET /health` - Health check (public)
- `GET /v1/cards?brand=visa&count=10&secret=<secret>` - Generate cards (protected)
- `GET /v1/scenarios` - List test scenarios (protected)

**Authentication:** Add header `Authorization: Bearer <token>`

**Rate Limiting:** 100 requests per minute per IP

### Validate Command

Validate a PAN using the Luhn algorithm.

```bash
cardgen-pro validate <PAN>
```

## 🔒 Security & Compliance

### Secret Management

**NEVER hardcode secrets in your code or configuration files.**

The `CARDGEN_SECRET` is used for deterministic CVC generation. Best practices:

```bash
# Environment variable (development)
export CARDGEN_SECRET="your-dev-secret"

# AWS Secrets Manager (production-like sandbox)
export CARDGEN_SECRET=$(aws secretsmanager get-secret-value --secret-id cardgen/dev --query SecretString --output text)

# HashiCorp Vault
export CARDGEN_SECRET=$(vault kv get -field=secret secret/cardgen/dev)

# Kubernetes Secret
kubectl create secret generic cardgen-secret --from-literal=secret=your-dev-secret
```

### PAN Masking

Always mask PANs in logs and UI:

```go
import "github.com/felipemacedo/cardgen-pro/internal/generator"

maskedPAN := generator.MaskPAN("4000000000000002")
// Output: "400000******0002"
```

### PCI-DSS Notice

**This tool generates SYNTHETIC TEST DATA and does NOT handle real cardholder data.**

- ❌ Do NOT use generated PANs on production payment networks
- ❌ Do NOT use this tool with real PANs
- ❌ Do NOT bypass PCI-DSS requirements with this tool
- ✅ Use ONLY in test/sandbox/development environments
- ✅ Rotate secrets regularly
- ✅ Implement proper access controls

See [SECURITY.md](./SECURITY.md) for detailed security guidelines.

## 🧪 Testing

Run the complete test suite:

```bash
# Unit tests
go test ./internal/... -v

# Integration tests
go test ./test/... -v

# All tests with coverage
go test ./... -cover

# Coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Coverage:** ≥85% (unit + integration)

## 📐 Architecture

```
cardgen-pro/
├── cmd/cardgen-pro/        # CLI entry point
├── internal/
│   ├── generator/          # PAN, Luhn, CVC, Track2 generation
│   ├── iso/                # ISO-8583 field builders
│   ├── api/                # HTTP API server & fixtures
│   └── models/             # Data structures
├── pkg/transformer/        # Order transformation & I/O
├── test/                   # Integration tests
├── fixtures/               # Sample test data
└── docs/                   # Additional documentation
```

See [ARCHITECTURE.md](./docs/ARCHITECTURE.md) for detailed design.

## 🎭 Test Scenarios

The tool includes 12 pre-built test scenarios:

1. **success_auth** - Standard approved transaction
2. **declined_generic** - Generic decline
3. **insufficient_funds** - Insufficient balance
4. **3ds_required** - 3D Secure authentication required
5. **auth_only** - Pre-authorization (not captured)
6. **captured** - Captured transaction
7. **refunded_partial** - Partial refund
8. **chargeback_open** - Chargeback dispute
9. **pix_paid** - Brazilian PIX instant payment
10. **boleto_pending** - Brazilian boleto pending
11. **subscription_recurring** - Recurring subscription payment
12. **tokenized_payment** - Token-based payment

List all scenarios:

```bash
cardgen-pro scenarios
```

## 🔬 Technical Details

### Luhn Algorithm

All generated PANs pass the Luhn mod-10 checksum (ISO/IEC 7812-1).

**Implementation:** See `internal/generator/luhn.go`

### Deterministic CVC Generation

CVCs are generated using HMAC-SHA256 for reproducibility:

```
Payload: BIN6 | LAST4 | EXPMM | EXPYYYY
CVC = HMAC-SHA256(Payload, Secret) -> first N digits
```

**Why HMAC-SHA256?**
- ✅ Cryptographically strong
- ✅ Deterministic (same input = same output)
- ✅ One-way (cannot reverse engineer)
- ✅ Secret-based (prevents guessing)

**Trade-offs:**
- Real CVVs are generated by issuer HSMs using proprietary algorithms
- This is a SIMULATION for testing only
- Provides realistic behavior without exposing real algorithms

See [SPEC.md](./docs/SPEC.md) for detailed algorithms.

### ISO-8583 Fields

The tool generates a **simplified** ISO-8583 message map with commonly used fields:

| Field | Name | Description |
|-------|------|-------------|
| 2 | Primary Account Number | Card PAN |
| 3 | Processing Code | Transaction type |
| 4 | Transaction Amount | Amount in cents |
| 7 | Transmission Date & Time | MMDDhhmmss |
| 11 | STAN | System Trace Audit Number |
| 14 | Expiration Date | YYMM |
| 22 | POS Entry Mode | How card was entered |
| 35 | Track 2 Data | Magnetic stripe data |
| 37 | RRN | Retrieval Reference Number |
| 41 | Terminal ID | Card acceptor terminal |
| 42 | Merchant ID | Card acceptor ID |
| 49 | Currency Code | ISO 4217 code |

**Note:** This is NOT a full ISO-8583 implementation. For production, use specialized libraries.

## 🐳 Docker

```dockerfile
# Build
docker build -t cardgen-pro .

# Run
docker run --rm -e CARDGEN_SECRET=dev-secret cardgen-pro generate --count 5
```

See `Dockerfile` for details.

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

### Development Setup

```bash
# Clone repository
git clone https://github.com/felipemacedo1/cardgen-pro.git
cd cardgen-pro

# Install dependencies
go mod download

# Run tests
go test ./... -v

# Lint
golangci-lint run

# Build
go build -o cardgen-pro ./cmd/cardgen-pro
```

## 📄 License

MIT License - see [LICENSE](./LICENSE) for details.

## 🙏 Acknowledgments

- ISO-8583 standard for financial transaction messages
- Luhn algorithm (Hans Peter Luhn, IBM)
- Payment card industry best practices

## ⚠️ Disclaimer

**THIS SOFTWARE IS FOR TESTING PURPOSES ONLY.**

Generated data is synthetic and must not be used:
- ❌ On production payment networks
- ❌ For actual financial transactions
- ❌ With real customer data
- ❌ To bypass security controls

The authors are not responsible for misuse of this tool.

---

**Made with ❤️ by Payment Engineers for Payment Engineers**

For issues, questions, or contributions: [GitHub Issues](https://github.com/felipemacedo1/cardgen-pro/issues)
