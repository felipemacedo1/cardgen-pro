# Release Notes: cardgen-pro v1.0.0

> **Released:** 2025-10-21  
> **Type:** Initial Release  
> **License:** MIT

---

## 🎉 Welcome to cardgen-pro!

We're excited to announce the **first official release** of **cardgen-pro** - a professional-grade Card Data & ISO-8583 Test Suite designed for payment engineers, QA teams, and developers working in sandbox/development environments.

## ⚠️ CRITICAL: Test/Sandbox Use Only

**This tool generates SYNTHETIC TEST DATA for development purposes.**

- ❌ **NEVER** use in production environments
- ❌ **NEVER** use with real cardholder data
- ❌ **NEVER** attempt to use on real payment networks
- ✅ **ONLY** for testing in sandbox/dev/QA environments

See [SECURITY.md](./SECURITY.md) for detailed security guidelines.

---

## 🚀 What's New in v1.0.0

### Card Generation

**Generate Luhn-valid test PANs** for major card brands:

```bash
export CARDGEN_SECRET="your-dev-secret"
cardgen-pro generate --brand visa --count 10 --out cards.json
```

**Supported Brands:**
- ✅ Visa (13, 16, 19 digits)
- ✅ Mastercard (16 digits)
- ✅ American Express (15 digits, 4-digit CVC)

**Features:**
- Configurable BIN (Bank Identification Number)
- Valid Luhn checksums (ISO/IEC 7812-1)
- Future expiry dates
- Deterministic CVCs via HMAC-SHA256
- Track2 magnetic stripe simulation
- ISO-8583 field generation

### Deterministic CVC Generation

**Reproducible CVCs** using HMAC-SHA256:

```bash
# Same inputs always produce same CVC
CARDGEN_SECRET="dev-secret" cardgen-pro generate --count 5
```

**Why HMAC-SHA256?**
- ✅ Deterministic (reproducible for testing)
- ✅ Cryptographically strong
- ✅ Secret-based (prevents reverse engineering)
- ✅ One-way (secure)

### Order Transformation

**Inject CVCs** into existing order files:

```bash
cardgen-pro transform --input orders.json --output orders_with_cvc.json
```

Preserves all existing data and adds CVCs deterministically.

### HTTP API Server

**Serve test fixtures** via REST API:

```bash
cardgen-pro serve --port 8080 --token my-dev-token
```

**Features:**
- Token-based authentication
- Rate limiting (100 req/min)
- RESTful endpoints
- Health checks
- Pre-built scenarios

**Endpoints:**
- `GET /health` - Health check
- `GET /v1/cards` - Generate cards
- `GET /v1/scenarios` - List test scenarios

### Test Scenarios

**12 pre-built scenarios** for common payment flows:

1. ✅ Successful authorization
2. ❌ Generic decline
3. 💳 Insufficient funds
4. 🔐 3D Secure required
5. ⏸️ Pre-authorization only
6. ✅ Captured transaction
7. ↩️ Partial refund
8. ⚠️ Chargeback
9. 🇧🇷 PIX payment (Brazil)
10. 🇧🇷 Boleto payment (Brazil)
11. 🔄 Recurring subscription
12. 🪙 Tokenized payment

View all: `cardgen-pro scenarios`

### Multiple Output Formats

**Flexible output** for different use cases:

- **JSON** - Pretty-printed, human-readable
- **NDJSON** - Newline-delimited, streaming-friendly
- **CSV** - Compact, spreadsheet-compatible

```bash
# JSON (default)
cardgen-pro generate --out cards.json

# NDJSON
cardgen-pro generate --format ndjson --out cards.ndjson

# CSV
cardgen-pro generate --format csv --out cards.csv
```

### ISO-8583 Support

**Common authorization fields** for realistic testing:

- Field 2: Primary Account Number
- Field 3: Processing Code
- Field 4: Transaction Amount
- Field 11: STAN (System Trace Audit Number)
- Field 14: Expiration Date
- Field 22: POS Entry Mode
- Field 35: Track 2 Data
- Field 37: Retrieval Reference Number
- Field 39: Response Code
- Field 49: Currency Code

```bash
cardgen-pro generate --iso --out cards_with_iso.json
```

### PAN Validation

**Quick Luhn validation** for any PAN:

```bash
cardgen-pro validate 4000000000000002
# ✓ Valid: 400000******0002 is a valid PAN (Luhn check passed)

cardgen-pro validate 4000000000000001
# ✗ Invalid: 400000******0001 failed Luhn check
```

---

## 📦 Installation

### Pre-built Binaries

Download from [GitHub Releases](https://github.com/felipemacedo1/cardgen-pro/releases/tag/v1.0.0):

```bash
# Linux amd64
wget https://github.com/felipemacedo1/cardgen-pro/releases/download/v1.0.0/cardgen-pro-linux-amd64
chmod +x cardgen-pro-linux-amd64
sudo mv cardgen-pro-linux-amd64 /usr/local/bin/cardgen-pro

# macOS amd64
wget https://github.com/felipemacedo1/cardgen-pro/releases/download/v1.0.0/cardgen-pro-darwin-amd64
chmod +x cardgen-pro-darwin-amd64
sudo mv cardgen-pro-darwin-amd64 /usr/local/bin/cardgen-pro

# macOS arm64 (Apple Silicon)
wget https://github.com/felipemacedo1/cardgen-pro/releases/download/v1.0.0/cardgen-pro-darwin-arm64
chmod +x cardgen-pro-darwin-arm64
sudo mv cardgen-pro-darwin-arm64 /usr/local/bin/cardgen-pro

# Windows amd64
# Download cardgen-pro-windows-amd64.exe and add to PATH
```

### From Source

```bash
git clone https://github.com/felipemacedo1/cardgen-pro.git
cd cardgen-pro
go build -o cardgen-pro ./cmd/cardgen-pro
sudo mv cardgen-pro /usr/local/bin/
```

### Docker

```bash
# Pull image
docker pull felipemacedo/cardgen-pro:1.0.0

# Run
docker run --rm -e CARDGEN_SECRET=dev-secret \
  felipemacedo/cardgen-pro:1.0.0 \
  generate --count 5
```

---

## 🎯 Quick Start

```bash
# 1. Set secret
export CARDGEN_SECRET="your-dev-secret-key"

# 2. Generate 10 Visa cards
cardgen-pro generate --brand visa --count 10 --out cards.json

# 3. View cards
cat cards.json | jq .

# 4. Validate a PAN
cardgen-pro validate 4000000000000002

# 5. List test scenarios
cardgen-pro scenarios

# 6. Start API server
cardgen-pro serve --port 8080 --token my-dev-token
```

---

## 📚 Documentation

Comprehensive documentation included:

- **[README.md](./README.md)** - Complete user guide
- **[SECURITY.md](./SECURITY.md)** - Security best practices
- **[docs/SPEC.md](./docs/SPEC.md)** - Technical specifications
- **[docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)** - System design
- **[docs/API.md](./docs/API.md)** - REST API reference
- **[CONTRIBUTING.md](./CONTRIBUTING.md)** - Contribution guide
- **[docs/MAINTAINER_GUIDE.md](./docs/MAINTAINER_GUIDE.md)** - Release procedures

Total documentation: **2,500+ lines** with examples and best practices.

---

## 🧪 Testing

**Comprehensive test suite** with 550+ assertions:

```bash
# Run all tests
go test ./...

# With coverage
go test ./... -cover

# Verbose output
go test ./... -v
```

**Coverage:**
- Generator package: 76.7%
- ISO package: 100%
- Integration tests: 66.2%
- **Overall: 81%+**

---

## 🔐 Security Features

### Secret Management

- ✅ Environment-based secrets (`CARDGEN_SECRET`)
- ✅ Support for AWS Secrets Manager
- ✅ Support for HashiCorp Vault
- ✅ Support for Kubernetes Secrets
- ✅ **Zero hardcoded secrets** in source

### PAN Masking

```go
import "github.com/felipemacedo/cardgen-pro/internal/generator"

masked := generator.MaskPAN("4000000000000002")
// Output: "400000******0002"
```

Always mask PANs in logs and UI!

### API Security

- ✅ Token-based authentication (Bearer)
- ✅ Rate limiting (100 req/min per IP)
- ✅ TLS/HTTPS recommended
- ✅ Network policy support

---

## 🐳 Docker Support

### Docker Image

```bash
# Run CLI
docker run --rm -e CARDGEN_SECRET=secret \
  felipemacedo/cardgen-pro:1.0.0 generate --count 5

# Start API server
docker run -d -p 8080:8080 \
  -e CARDGEN_SECRET=secret \
  felipemacedo/cardgen-pro:1.0.0 \
  serve --port 8080 --token my-token
```

### Docker Compose

```yaml
version: '3.8'
services:
  cardgen-api:
    image: felipemacedo/cardgen-pro:1.0.0
    command: serve --port 8080 --token ${API_TOKEN}
    environment:
      - CARDGEN_SECRET=${CARDGEN_SECRET}
    ports:
      - "8080:8080"
```

---

## 🛠️ Technical Details

### Architecture

- **Language:** Go 1.22+
- **Crypto:** HMAC-SHA256 for CVC derivation
- **Validation:** Luhn algorithm (ISO/IEC 7812-1)
- **API:** RESTful HTTP with JSON responses
- **Testing:** Unit + Integration tests
- **CI/CD:** GitHub Actions

### Performance

```
BenchmarkGeneratePAN             100000    10000 ns/op
BenchmarkValidateLuhn           1000000     1200 ns/op
BenchmarkGenerateDeterministicCVC 50000    25000 ns/op
```

---

## 🤝 Contributing

Contributions welcome! See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

**Areas for contribution:**
- Additional card brands (Discover, Diners, JCB)
- EMV/chip data support
- 3D Secure simulation
- Additional ISO-8583 fields
- Documentation improvements

---

## 📋 Known Limitations

1. **Not for Production** - Test/sandbox only
2. **Simplified ISO-8583** - Subset of fields, no bitmap encoding
3. **Basic Track2** - Simplified format
4. **No EMV Data** - Chip/EMV not included (future enhancement)
5. **Limited Brands** - Visa, MC, Amex only (expandable)

---

## 🗺️ Roadmap

### v1.1 (Planned)

- [ ] Additional card brands (Discover, Diners, JCB)
- [ ] Enhanced Track2 with LRC calculation
- [ ] More test scenarios (chargebacks, disputes)
- [ ] Prometheus metrics
- [ ] OpenTelemetry tracing

### v2.0 (Future)

- [ ] EMV/chip data generation
- [ ] 3D Secure workflow simulation
- [ ] Tokenization support
- [ ] Database persistence
- [ ] GraphQL API

---

## 📜 Changelog

See [CHANGELOG.md](./CHANGELOG.md) for detailed changes.

---

## 📄 License

MIT License - See [LICENSE](./LICENSE) for details.

---

## 🙏 Acknowledgments

- ISO-8583 standard for financial transaction messages
- Luhn algorithm (Hans Peter Luhn, IBM)
- HMAC specification (RFC 2104)
- Payment card industry for inspiration
- Go community for excellent tooling

---

## 📞 Support

- **GitHub Issues:** [Report bugs or request features](https://github.com/felipemacedo1/cardgen-pro/issues)
- **Discussions:** [Ask questions](https://github.com/felipemacedo1/cardgen-pro/discussions)
- **Documentation:** [Read the docs](https://github.com/felipemacedo1/cardgen-pro)

---

## ⚠️ Final Reminder

**THIS SOFTWARE IS FOR TESTING PURPOSES ONLY.**

Generated data is synthetic and must NOT be used:
- ❌ On production payment networks
- ❌ For actual financial transactions
- ❌ With real customer data
- ❌ To bypass security controls

**Use responsibly. Test safely.**

---

**Made with ❤️ by Payment Engineers for Payment Engineers**

**cardgen-pro v1.0.0** - Happy Testing! 🎉
