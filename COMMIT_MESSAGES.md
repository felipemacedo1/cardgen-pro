# Commit Messages (Conventional Commits)

This document contains suggested commit messages for the cardgen-pro v1.0.0 release, following [Conventional Commits](https://www.conventionalcommits.org/) specification.

## Core Features

```bash
git commit -m "feat(generator): add PAN generation with Luhn validation

- Implement Luhn algorithm for checksum validation (ISO/IEC 7812-1)
- Add CalculateLuhnCheckDigit function
- Add AppendLuhnCheckDigit function
- Add ValidateLuhn function
- Support PAN lengths 13-19 digits
- Add comprehensive unit tests with 80%+ coverage

BREAKING CHANGE: None (initial release)"
```

```bash
git commit -m "feat(generator): add deterministic CVC generation via HMAC-SHA256

- Implement GenerateDeterministicCVC using HMAC-SHA256
- Support 3-digit CVC (Visa/Mastercard) and 4-digit (Amex)
- Payload: BIN6|LAST4|EXPMM|EXPYYYY
- Secret-based for security
- Fully deterministic for test reproducibility
- Add tests verifying determinism

Rationale: HMAC-SHA256 provides cryptographic strength while
remaining deterministic, essential for reproducible testing."
```

```bash
git commit -m "feat(generator): add Track2 magnetic stripe data generation

- Implement GenerateTrack2 function
- Format: PAN=YYMM<ServiceCode><Discretionary>
- Support standard service codes (201, 101, 121)
- Generate random discretionary data
- Add validation tests

Note: Simplified format for testing, not production-grade."
```

```bash
git commit -m "feat(generator): add expiry date generation

- Implement GenerateExpiry function
- Random month (1-12)
- Random future year (current + 1 to 5)
- Use crypto/rand for secure randomness
- Add boundary tests"
```

```bash
git commit -m "feat(generator): add PAN masking utility

- Implement MaskPAN function
- Format: first6****last4 (e.g., 400000******0002)
- Handles short PANs safely
- Add comprehensive masking tests

Security: Essential for safe logging and display."
```

```bash
git commit -m "feat(iso): add ISO-8583 field generation

- Implement GenerateISO8583Fields function
- Support 20+ common fields (2,3,4,7,11,12,13,14,22,35,37,41,42,49)
- Generate mock authorization request (MTI 0100)
- Generate mock authorization response (MTI 0110)
- Add response code mapping
- Add field formatting utility
- 100% test coverage

Note: Simplified implementation for testing, not full ISO-8583 lib."
```

```bash
git commit -m "feat(transformer): add order transformation with CVC injection

- Implement TransformOrders function
- Read orders from JSON/NDJSON
- Inject deterministic CVCs
- Preserve existing order data
- Write transformed orders
- Add integration tests

Use case: Inject CVCs into existing test order files."
```

```bash
git commit -m "feat(transformer): add multiple output format support

- Implement WriteCardsJSON (pretty-printed)
- Implement WriteCardsNDJSON (newline-delimited)
- Implement WriteCardsCSV (comma-separated)
- Add format tests
- Add file I/O error handling

Supports various integration scenarios and tooling."
```

```bash
git commit -m "feat(api): add HTTP API server with authentication

- Implement HTTP server with Gin-like routing
- Add token-based authentication middleware
- Add rate limiting (100 req/min per IP)
- Add /health endpoint (public)
- Add /v1/cards endpoint (protected)
- Add /v1/scenarios endpoint (protected)
- Add error handling and logging

Security: Rate limiting prevents abuse, token auth required."
```

```bash
git commit -m "feat(api): add 12 pre-built test scenarios

- Implement GetScenarios function
- Add success, decline, insufficient_funds scenarios
- Add 3DS, pre-auth, capture, refund scenarios
- Add chargeback, PIX, boleto scenarios
- Add subscription and tokenization scenarios
- Each scenario includes expected outcome

Provides realistic test cases out-of-the-box."
```

```bash
git commit -m "feat(cli): add generate command

- Implement generate command handler
- Support --bin, --brand, --count flags
- Support --out, --format flags
- Support --iso, --track2 flags
- Support --secret flag (with CARDGEN_SECRET fallback)
- Add output to stdout or file
- Add sample card display
- Add logging

Primary use case: Generate test card data."
```

```bash
git commit -m "feat(cli): add transform command

- Implement transform command handler
- Support --input, --output flags
- Support --secret flag (with CARDGEN_SECRET fallback)
- Call transformer.TransformOrders
- Add error handling and logging

Use case: Inject CVCs into existing order files."
```

```bash
git commit -m "feat(cli): add serve command

- Implement serve command handler
- Support --port, --token flags
- Start HTTP API server
- Add graceful shutdown
- Add logging

Use case: Serve fixtures via REST API for integration tests."
```

```bash
git commit -m "feat(cli): add validate command

- Implement validate command handler
- Accept PAN as argument
- Call ValidateLuhn
- Display masked PAN
- Exit code 0 for valid, 1 for invalid

Use case: Quick PAN validation in scripts."
```

```bash
git commit -m "feat(cli): add scenarios command

- Implement scenarios command handler
- Call api.GetScenarios
- Format and display all scenarios
- Include ID, name, description, response, outcome

Use case: List available test scenarios."
```

```bash
git commit -m "feat(cli): add version and help commands

- Add version constant (1.0.0)
- Add ASCII art banner
- Implement version command
- Implement help command with usage
- Add command routing
- Add error handling for unknown commands

Improves CLI usability."
```

```bash
git commit -m "feat(models): add core data structures

- Add Card struct with JSON tags
- Add Order struct
- Add CardBrand struct with BIN ranges
- Add GenerateOptions struct
- Add TransformOptions struct
- Add CardBrands map with Visa, Mastercard, Amex configs

Foundation for type-safe data handling."
```

## Testing

```bash
git commit -m "test(generator): add comprehensive unit tests

- Test ValidateLuhn with valid/invalid PANs
- Test CalculateLuhnCheckDigit accuracy
- Test AppendLuhnCheckDigit produces valid PANs
- Test GeneratePAN with various BINs and lengths
- Test GenerateExpiry bounds
- Test GenerateDeterministicCVC determinism
- Test GenerateDeterministicCVC different secrets
- Test GenerateTrack2 format
- Test MaskPAN formats
- Add benchmark tests

Coverage: 76.7%"
```

```bash
git commit -m "test(iso): add ISO-8583 field tests

- Test GenerateISO8583Fields required fields
- Test field formats (PAN, amount, expiry, etc.)
- Test GenerateMockAuthRequest MTI
- Test GenerateMockAuthResponse with approve/decline
- Test ResponseCodes map
- Test FormatISO8583 output
- Add benchmark tests

Coverage: 100%"
```

```bash
git commit -m "test(integration): add end-to-end workflow tests

- Test generate → transform → verify CVC workflow
- Test CVC determinism over multiple iterations
- Test different secrets produce different CVCs
- Test card generation for all brands
- Test all output formats (JSON, NDJSON, CSV)
- Add file I/O cleanup

Coverage: 66.2%"
```

## Documentation

```bash
git commit -m "docs(readme): add comprehensive README

- Add project description and features
- Add installation instructions (source, binary, Docker)
- Add quick start guide with examples
- Add all CLI commands documentation
- Add security warnings and PCI-DSS notice
- Add technical details (Luhn, HMAC-SHA256)
- Add ISO-8583 field table
- Add contributing section
- Add license and disclaimer

2,500+ lines of documentation."
```

```bash
git commit -m "docs(security): add security guidelines and best practices

- Add critical warnings (test/sandbox only)
- Add secret management guide
- Add support for AWS Secrets Manager, Vault, K8s
- Add PAN masking best practices
- Add data cleanup procedures
- Add secret rotation guide
- Add access control recommendations
- Add API server security
- Add logging best practices
- Add incident response procedures
- Add compliance checklist

Essential for secure usage."
```

```bash
git commit -m "docs(spec): add technical specifications

- Document Luhn algorithm
- Document PAN structure and generation
- Document HMAC-SHA256 CVC derivation
- Document Track2 format
- Document ISO-8583 fields and MTI codes
- Document response codes
- Document data formats (JSON, NDJSON, CSV)
- Add examples and rationale

Technical reference for developers."
```

```bash
git commit -m "docs(architecture): add architecture overview

- Add system design diagram (ASCII art)
- Document component responsibilities
- Document data flow for each command
- Document security architecture
- Document testing architecture (test pyramid)
- Document performance considerations
- Document deployment architectures
- Document design decisions and trade-offs
- Add extensibility points

Explains system design and decisions."
```

```bash
git commit -m "docs(api): add REST API documentation

- Document all endpoints (/health, /v1/cards, /v1/scenarios)
- Add authentication guide
- Add rate limiting details
- Add request/response examples
- Add client examples (cURL, Python, JS, Go)
- Add deployment examples (Docker, K8s)
- Add security best practices
- Add troubleshooting guide

Complete API reference."
```

```bash
git commit -m "docs(contributing): add contribution guidelines

- Add code of conduct
- Add issue reporting guide
- Add PR process
- Add development setup
- Add code style guidelines
- Add testing requirements
- Add commit message format (Conventional Commits)
- Add PR template

Helps contributors get started."
```

```bash
git commit -m "docs(maintainer): add maintainer guide

- Add release process with checklist
- Add semantic versioning guide
- Add hotfix process
- Add dependency management
- Add CI/CD maintenance
- Add code quality maintenance
- Add security maintenance
- Add documentation maintenance
- Add community management
- Add emergency procedures

For project maintainers."
```

```bash
git commit -m "docs(changelog): add CHANGELOG for v1.0.0

- Document all features added
- Document technical decisions
- Document known limitations
- Document security notes
- Follow Keep a Changelog format
- Follow Semantic Versioning

Tracks project history."
```

## Configuration & Build

```bash
git commit -m "chore(docker): add Dockerfile and multi-stage build

- Add Go builder stage
- Add Alpine runtime stage
- Create non-root user
- Expose port 8080
- Add labels (OCI image spec)
- Optimize for size (~20MB final image)

Enables containerized deployment."
```

```bash
git commit -m "chore(ci): add GitHub Actions CI pipeline

- Add test job (unit + integration)
- Add lint job (golangci-lint)
- Add build job
- Add Docker build job
- Add coverage upload to Codecov
- Add coverage threshold check (80%)
- Cache Go modules for speed
- Run on push and PR

Automated testing and quality checks."
```

```bash
git commit -m "chore(ci): add GitHub Actions release pipeline

- Add release job triggered on git tags
- Build binaries for Linux, macOS, Windows (amd64, arm64)
- Generate checksums
- Create GitHub Release with artifacts
- Build and push Docker images (multi-platform)
- Add release notes template

Automated release process."
```

```bash
git commit -m "chore(lint): add golangci-lint configuration

- Enable 20+ linters (gofmt, govet, gosec, etc.)
- Set gocyclo complexity limit
- Set dupl threshold
- Exclude test files from some linters
- Set timeout to 5 minutes

Enforces code quality."
```

```bash
git commit -m "chore(git): add .gitignore

- Ignore binaries
- Ignore test outputs
- Ignore coverage files
- Ignore IDE files
- Ignore OS files
- Ignore generated fixtures
- Ignore .env files

Keeps repo clean."
```

```bash
git commit -m "chore(fixtures): add sample test fixtures

- Generate 5 Visa cards
- Generate 5 Mastercard cards
- Generate 3 Amex cards
- Add sample orders
- Add transformed orders with CVCs
- Add fixtures README

Example data for testing."
```

```bash
git commit -m "chore(license): add MIT license

- Add MIT License text
- Add copyright notice
- Add additional disclaimer for test/sandbox use
- Add prohibited use cases
- Add liability disclaimer

Open source under MIT."
```

```bash
git commit -m "chore(deps): add go.mod and dependencies

- Initialize Go module
- Set Go version to 1.22
- Add minimal dependencies (stdlib only)
- Run go mod tidy

Zero external dependencies (stdlib only)."
```

## Final Release

```bash
git commit -m "chore(release): prepare v1.0.0

- Update version to 1.0.0
- Update CHANGELOG.md with release notes
- Verify all tests passing
- Verify documentation complete
- Verify security reviewed
- Verify CI/CD configured

Ready for release."
```

---

## How to Use These Commits

If you're applying these commits manually:

```bash
# Apply commits in order
git add <files>
git commit -F <commit-message-from-above>

# Or squash all at once
git add .
git commit -m "feat: initial release of cardgen-pro v1.0.0

Complete Card Data & ISO-8583 Test Suite for sandbox/development.

Features:
- PAN generation with Luhn validation
- Deterministic CVC via HMAC-SHA256
- Track2 and ISO-8583 field generation
- CLI commands (generate, transform, serve, validate, scenarios)
- HTTP API with auth and rate limiting
- 12 pre-built test scenarios
- Multiple output formats (JSON, NDJSON, CSV)
- Comprehensive documentation (2,500+ lines)
- Unit + integration tests (81% coverage)
- CI/CD with GitHub Actions

See CHANGELOG.md and RELEASE_NOTES for details.

BREAKING CHANGE: None (initial release)"

# Tag release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

---

**Note:** These messages follow Conventional Commits specification for automated changelog generation and semantic versioning.
