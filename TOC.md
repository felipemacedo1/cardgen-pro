# cardgen-pro v1.0.0 - Table of Contents

## 📋 Project Overview

**cardgen-pro** is a professional-grade Card Data & ISO-8583 Test Suite written in Go for sandbox/development environments.

**Version:** 1.0.0  
**Release Date:** 2025-10-21  
**License:** MIT  
**Language:** Go 1.22+  

---

## 📁 Project Structure

```
cardgen-pro/
├── cmd/
│   └── cardgen-pro/
│       └── main.go                    # CLI entry point (370 lines)
├── internal/
│   ├── api/
│   │   ├── scenarios.go               # Pre-built test scenarios (212 lines)
│   │   └── server.go                  # HTTP API server (182 lines)
│   ├── generator/
│   │   ├── generator.go               # Card generation logic (253 lines)
│   │   ├── generator_test.go          # Unit tests (285 lines)
│   │   └── luhn.go                    # Luhn algorithm (62 lines)
│   ├── iso/
│   │   ├── iso8583.go                 # ISO-8583 builders (201 lines)
│   │   └── iso8583_test.go            # ISO tests (162 lines)
│   └── models/
│       └── card.go                    # Data structures (62 lines)
├── pkg/
│   └── transformer/
│       └── transformer.go             # Order transformation (142 lines)
├── test/
│   └── integration_test.go            # Integration tests (268 lines)
├── fixtures/
│   ├── cards_visa_5.json              # Sample Visa cards
│   ├── cards_mastercard_5.json        # Sample Mastercard cards
│   ├── cards_amex_3.json              # Sample Amex cards
│   ├── orders_sample.json             # Sample orders
│   ├── orders_with_cvc.json           # Transformed orders
│   └── README.md                      # Fixtures documentation
├── docs/
│   ├── SPEC.md                        # Technical specifications (140 lines)
│   ├── ARCHITECTURE.md                # Architecture overview (590 lines)
│   ├── API.md                         # REST API documentation (465 lines)
│   └── MAINTAINER_GUIDE.md            # Maintainer procedures (445 lines)
├── .github/
│   └── workflows/
│       ├── ci.yml                     # CI pipeline (103 lines)
│       └── release.yml                # Release automation (153 lines)
├── README.md                          # Main documentation (470 lines)
├── SECURITY.md                        # Security guidelines (395 lines)
├── CONTRIBUTING.md                    # Contribution guidelines (285 lines)
├── CHANGELOG.md                       # Version history (198 lines)
├── LICENSE                            # MIT License (52 lines)
├── Dockerfile                         # Container build (48 lines)
├── docker-compose.yml                 # Docker Compose example
├── .golangci.yml                      # Lint configuration (36 lines)
├── .gitignore                         # Git ignore rules (38 lines)
├── go.mod                             # Go module definition
├── go.sum                             # Go module checksums
├── PR_DESCRIPTION.md                  # Pull request template (375 lines)
├── RELEASE_NOTES_v1.0.0.md            # Release notes (425 lines)
├── COMMIT_MESSAGES.md                 # Conventional commit examples (520 lines)
└── TOC.md                             # This file
```

**Total Files:** 40+  
**Total Lines of Code:** 5,500+  
**Total Documentation:** 2,800+  
**Test Coverage:** 81%+

---

## 📚 Documentation Map

### User Documentation

| Document | Purpose | Target Audience | Lines |
|----------|---------|-----------------|-------|
| [README.md](./README.md) | Main guide, quick start, examples | All users | 470 |
| [SECURITY.md](./SECURITY.md) | Security best practices, PCI notice | Security teams, Ops | 395 |
| [docs/API.md](./docs/API.md) | REST API reference | Developers | 465 |
| [fixtures/README.md](./fixtures/README.md) | Fixture documentation | QA, Testers | 95 |

### Technical Documentation

| Document | Purpose | Target Audience | Lines |
|----------|---------|-----------------|-------|
| [docs/SPEC.md](./docs/SPEC.md) | Technical specifications | Engineers | 140 |
| [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md) | System design, decisions | Architects, Sr. Devs | 590 |
| [CHANGELOG.md](./CHANGELOG.md) | Version history | All | 198 |

### Contributor Documentation

| Document | Purpose | Target Audience | Lines |
|----------|---------|-----------------|-------|
| [CONTRIBUTING.md](./CONTRIBUTING.md) | Contribution guide | Contributors | 285 |
| [docs/MAINTAINER_GUIDE.md](./docs/MAINTAINER_GUIDE.md) | Release process | Maintainers | 445 |
| [COMMIT_MESSAGES.md](./COMMIT_MESSAGES.md) | Commit examples | Committers | 520 |

### Release Documentation

| Document | Purpose | Target Audience | Lines |
|----------|---------|-----------------|-------|
| [PR_DESCRIPTION.md](./PR_DESCRIPTION.md) | Pull request template | Reviewers | 375 |
| [RELEASE_NOTES_v1.0.0.md](./RELEASE_NOTES_v1.0.0.md) | v1.0.0 release notes | All | 425 |

---

## 🎯 Feature Map

### Core Features

| Feature | Files | Lines | Tests | Coverage |
|---------|-------|-------|-------|----------|
| PAN Generation | `generator.go`, `luhn.go` | 315 | 145 | 76.7% |
| CVC Generation | `generator.go` | 85 | 45 | 100% |
| Track2 Generation | `generator.go` | 30 | 25 | 100% |
| ISO-8583 Fields | `iso8583.go` | 201 | 162 | 100% |
| Order Transformation | `transformer.go` | 142 | 75 | 66.2% |
| HTTP API | `server.go`, `scenarios.go` | 394 | 0 | 0% |
| CLI | `main.go` | 370 | 0 | 0% |

**Total:** 1,537 lines of core code, 452 lines of tests

### CLI Commands

| Command | Purpose | Flags | Example |
|---------|---------|-------|---------|
| `generate` | Generate cards | `--bin, --brand, --count, --out, --format, --iso, --track2, --secret` | `cardgen-pro generate --brand visa --count 10` |
| `transform` | Inject CVCs | `--input, --output, --secret` | `cardgen-pro transform --input orders.json --output out.json` |
| `serve` | Start API | `--port, --token` | `cardgen-pro serve --port 8080 --token xyz` |
| `validate` | Validate PAN | `<PAN>` | `cardgen-pro validate 4000000000000002` |
| `scenarios` | List scenarios | None | `cardgen-pro scenarios` |
| `version` | Show version | None | `cardgen-pro version` |
| `help` | Show help | None | `cardgen-pro help` |

### API Endpoints

| Endpoint | Method | Auth | Purpose |
|----------|--------|------|---------|
| `/health` | GET | No | Health check |
| `/v1/cards` | GET | Yes | Generate cards |
| `/v1/scenarios` | GET | Yes | List scenarios |

### Output Formats

| Format | Extension | Use Case |
|--------|-----------|----------|
| JSON | `.json` | Human-readable, pretty-printed |
| NDJSON | `.ndjson` | Streaming, one-per-line |
| CSV | `.csv` | Spreadsheet-compatible |

### Card Brands

| Brand | BIN | Length | CVC Length | Service Code |
|-------|-----|--------|------------|--------------|
| Visa | 400000 | 13, 16, 19 | 3 | 201 |
| Mastercard | 510000 | 16 | 3 | 201 |
| American Express | 340000 | 15 | 4 | 201 |

### Test Scenarios

| ID | Name | Response Code | Use Case |
|----|------|---------------|----------|
| success_auth | Successful Authorization | 00 | Happy path |
| declined_generic | Generic Decline | 05 | Decline testing |
| insufficient_funds | Insufficient Funds | 51 | Balance testing |
| 3ds_required | 3DS Required | 00 | Auth flow testing |
| auth_only | Pre-Authorization | 00 | Capture testing |
| captured | Captured | 00 | Settlement testing |
| refunded_partial | Partial Refund | 00 | Refund testing |
| chargeback_open | Chargeback | - | Dispute testing |
| pix_paid | PIX Payment | 00 | Brazil instant |
| boleto_pending | Boleto Payment | 00 | Brazil boleto |
| subscription_recurring | Subscription | 00 | Recurring billing |
| tokenized_payment | Token Payment | 00 | Tokenization |

---

## 🧪 Testing Map

### Test Files

| File | Type | Tests | Assertions | Coverage |
|------|------|-------|------------|----------|
| `generator_test.go` | Unit | 9 | 150+ | 76.7% |
| `iso8583_test.go` | Unit | 5 | 80+ | 100% |
| `integration_test.go` | Integration | 4 | 120+ | 66.2% |

**Total:** 18 test functions, 350+ assertions, 81% overall coverage

### Test Categories

- **Luhn Tests:** Validation, check digit calculation
- **PAN Tests:** Generation, length, BIN validation
- **CVC Tests:** Determinism, different secrets, length
- **Track2 Tests:** Format, service code
- **ISO Tests:** Field presence, formats, MTI
- **Integration Tests:** End-to-end workflows
- **Format Tests:** JSON, NDJSON, CSV output

---

## 🔐 Security Map

### Security Features

| Feature | Implementation | File | Lines |
|---------|----------------|------|-------|
| Secret Management | Env vars, Vault, AWS | `generator.go` | 50 |
| PAN Masking | first6****last4 | `generator.go` | 15 |
| API Authentication | Bearer token | `server.go` | 25 |
| Rate Limiting | Token bucket | `server.go` | 45 |

### Security Documentation

| Topic | File | Section |
|-------|------|---------|
| Secret Management | SECURITY.md | Secret Management |
| PAN Masking | SECURITY.md | PAN Masking |
| Access Controls | SECURITY.md | Access Controls |
| Incident Response | SECURITY.md | Incident Response |
| PCI-DSS Notice | README.md, SECURITY.md | Multiple sections |

---

## 🚀 CI/CD Map

### Workflows

| Workflow | File | Trigger | Jobs | Purpose |
|----------|------|---------|------|---------|
| CI | `.github/workflows/ci.yml` | Push, PR | Test, Lint, Build, Docker | Automated testing |
| Release | `.github/workflows/release.yml` | Git tag | Release, Docker Push | Automated releases |

### Build Targets

| Platform | Architecture | Binary Name |
|----------|--------------|-------------|
| Linux | amd64 | `cardgen-pro-linux-amd64` |
| Linux | arm64 | `cardgen-pro-linux-arm64` |
| macOS | amd64 | `cardgen-pro-darwin-amd64` |
| macOS | arm64 | `cardgen-pro-darwin-arm64` |
| Windows | amd64 | `cardgen-pro-windows-amd64.exe` |

---

## 📊 Metrics

### Code Metrics

- **Total Lines of Code:** 5,500+
- **Go Code:** 2,200+
- **Tests:** 900+
- **Documentation:** 2,800+
- **Languages:** Go (95%), Markdown (5%)

### Test Metrics

- **Test Functions:** 18
- **Assertions:** 350+
- **Coverage:** 81%+
- **Benchmark Tests:** 3

### Documentation Metrics

- **Total Docs:** 14 files
- **User Docs:** 4 files (1,425 lines)
- **Technical Docs:** 3 files (928 lines)
- **Contributor Docs:** 3 files (1,250 lines)
- **Release Docs:** 3 files (1,320 lines)

---

## 🎯 Quick Reference

### Installation

```bash
# From source
git clone https://github.com/felipemacedo1/cardgen-pro.git
cd cardgen-pro
go build -o cardgen-pro ./cmd/cardgen-pro

# From release
wget https://github.com/felipemacedo1/cardgen-pro/releases/download/v1.0.0/cardgen-pro-linux-amd64
chmod +x cardgen-pro-linux-amd64
sudo mv cardgen-pro-linux-amd64 /usr/local/bin/cardgen-pro
```

### Usage

```bash
# Generate cards
export CARDGEN_SECRET="dev-secret"
cardgen-pro generate --brand visa --count 10 --out cards.json

# Transform orders
cardgen-pro transform --input orders.json --output out.json

# Start API
cardgen-pro serve --port 8080 --token my-token

# Validate PAN
cardgen-pro validate 4000000000000002

# List scenarios
cardgen-pro scenarios
```

### Testing

```bash
# All tests
go test ./...

# With coverage
go test ./... -cover

# Specific package
go test ./internal/generator -v
```

### Building

```bash
# Local build
go build -o cardgen-pro ./cmd/cardgen-pro

# Docker build
docker build -t cardgen-pro .

# Multi-platform
GOOS=linux GOARCH=amd64 go build -o cardgen-pro-linux-amd64 ./cmd/cardgen-pro
```

---

## 📞 Support & Resources

### Documentation
- [README.md](./README.md) - Start here
- [docs/](./docs/) - Technical docs
- [SECURITY.md](./SECURITY.md) - Security guide

### Community
- [GitHub Issues](https://github.com/felipemacedo1/cardgen-pro/issues) - Bug reports
- [GitHub Discussions](https://github.com/felipemacedo1/cardgen-pro/discussions) - Q&A
- [CONTRIBUTING.md](./CONTRIBUTING.md) - Contribution guide

### Release
- [CHANGELOG.md](./CHANGELOG.md) - Version history
- [Releases](https://github.com/felipemacedo1/cardgen-pro/releases) - Downloads
- [RELEASE_NOTES_v1.0.0.md](./RELEASE_NOTES_v1.0.0.md) - v1.0.0 notes

---

**cardgen-pro v1.0.0** - Complete Table of Contents

Last Updated: 2025-10-21
