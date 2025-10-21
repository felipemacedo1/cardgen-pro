# Architecture Overview

## System Design

**cardgen-pro** is designed as a modular, layered system with clear separation of concerns.

```
┌─────────────────────────────────────────────────────────────┐
│                        CLI Interface                         │
│                   (cmd/cardgen-pro/main.go)                  │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Generate   │  │  Transform   │  │    Serve     │      │
│  │   Command    │  │   Command    │  │   Command    │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
└─────────┼──────────────────┼──────────────────┼─────────────┘
          │                  │                  │
          ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────┐
│                      Core Business Logic                     │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              Card Generator (internal/generator/)       │ │
│  │  • PAN Generation        • Expiry Generation           │ │
│  │  • Luhn Validation       • Track2 Generation           │ │
│  │  • CVC Generation        • PAN Masking                 │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │           ISO-8583 Builder (internal/iso/)             │ │
│  │  • Field Generation      • Mock Request/Response       │ │
│  │  • Response Codes        • Field Formatting            │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │           Transformer (pkg/transformer/)               │ │
│  │  • Order Reading         • CVC Injection               │ │
│  │  • Format Conversion     • File I/O                    │ │
│  └────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              API Server (internal/api/)                │ │
│  │  • HTTP Endpoints        • Authentication              │ │
│  │  • Rate Limiting         • Scenario Fixtures           │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      Data Layer                              │
│  ┌────────────────────────────────────────────────────────┐ │
│  │              Models (internal/models/)                 │ │
│  │  • Card            • Order                             │ │
│  │  • CardBrand       • GenerateOptions                   │ │
│  │  • TransformOptions                                    │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      Output Layer                            │
│     JSON    │    NDJSON    │    CSV    │   HTTP API         │
└─────────────────────────────────────────────────────────────┘
```

## Component Responsibilities

### CLI Layer (`cmd/cardgen-pro/`)

**Responsibility:** User interaction and command routing

- Parse command-line arguments
- Validate user input
- Route to appropriate command handlers
- Handle output formatting
- Display help and version information

### Core Generators (`internal/generator/`)

**Responsibility:** Card data generation logic

**Key Components:**
- `luhn.go` - Luhn algorithm validation and check digit calculation
- `generator.go` - PAN, expiry, CVC, Track2 generation
- `generator_test.go` - Comprehensive unit tests

**Key Algorithms:**
1. **Luhn Check:** Validates card numbers (ISO/IEC 7812-1)
2. **PAN Generation:** Creates valid card numbers with BIN + random + check digit
3. **CVC Derivation:** HMAC-SHA256(BIN6|LAST4|EXPMM|EXPYYYY, secret)
4. **Track2 Construction:** Formats magnetic stripe data

### ISO-8583 Handler (`internal/iso/`)

**Responsibility:** Financial message simulation

**Key Components:**
- `iso8583.go` - Field generation and message building
- `iso8583_test.go` - ISO field validation tests

**Key Features:**
- Common field generation (2, 3, 4, 11, 14, 22, 35, 37, 39, 41, 42, 49)
- Mock authorization request/response
- Response code mapping
- Field formatting utilities

### Transformer (`pkg/transformer/`)

**Responsibility:** File I/O and data transformation

**Key Components:**
- `transformer.go` - Order transformation and format conversion

**Supported Formats:**
- JSON (pretty-printed)
- NDJSON (newline-delimited)
- CSV (comma-separated values)

**Operations:**
- Read orders from file
- Inject CVCs deterministically
- Write transformed data
- Format conversion

### API Server (`internal/api/`)

**Responsibility:** HTTP fixture serving (optional)

**Key Components:**
- `server.go` - HTTP server and middleware
- `scenarios.go` - Pre-built test scenarios

**Features:**
- RESTful endpoints
- Token-based authentication
- Rate limiting (100 req/min per IP)
- Health checks
- Scenario listing

**Endpoints:**
- `GET /health` - Health check (public)
- `GET /v1/cards` - Generate cards (protected)
- `GET /v1/scenarios` - List scenarios (protected)

### Models (`internal/models/`)

**Responsibility:** Data structures

**Key Types:**
- `Card` - Generated card data
- `Order` - Payment order/transaction
- `CardBrand` - Brand configuration
- `GenerateOptions` - Generation parameters
- `TransformOptions` - Transformation parameters

## Data Flow

### Generate Command Flow

```
User Input
    │
    ├─> Parse CLI args (--bin, --brand, --count, --secret)
    │
    ├─> Load brand configuration (CardBrands map)
    │
    ├─> FOR each card:
    │   │
    │   ├─> Generate PAN (BIN + random + Luhn check digit)
    │   │
    │   ├─> Generate expiry (random future date)
    │   │
    │   ├─> Generate CVC (HMAC-SHA256)
    │   │
    │   ├─> [Optional] Generate Track2
    │   │
    │   └─> [Optional] Generate ISO-8583 fields
    │
    └─> Output to file or stdout (JSON/NDJSON/CSV)
```

### Transform Command Flow

```
User Input (--input, --output, --secret)
    │
    ├─> Read orders from input file (JSON/NDJSON)
    │
    ├─> FOR each order:
    │   │
    │   ├─> Extract PAN, expiry
    │   │
    │   ├─> Generate deterministic CVC (HMAC-SHA256)
    │   │
    │   └─> Inject CVC into order
    │
    └─> Write transformed orders to output file (JSON)
```

### Serve Command Flow

```
User Input (--port, --token)
    │
    ├─> Initialize HTTP server
    │
    ├─> Register routes and middleware
    │   │
    │   ├─> Authentication middleware (Bearer token)
    │   │
    │   └─> Rate limiting middleware (100 req/min)
    │
    └─> Listen and serve
        │
        ├─> GET /health → return status
        │
        ├─> GET /v1/cards → generate cards on-the-fly
        │
        └─> GET /v1/scenarios → return fixture list
```

## Security Architecture

### Secret Management

```
┌─────────────────────────────────────────────────────────────┐
│                      Secret Sources                          │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Environment  │  │    Vault     │  │  AWS Secrets │      │
│  │   Variable   │  │  (HashiCorp) │  │   Manager    │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
└─────────┼──────────────────┼──────────────────┼─────────────┘
          │                  │                  │
          └──────────────────┴──────────────────┘
                            │
                            ▼
                   ┌─────────────────┐
                   │  CARDGEN_SECRET │
                   └────────┬────────┘
                            │
                            ▼
                   ┌─────────────────┐
                   │  HMAC-SHA256    │
                   │   Algorithm     │
                   └────────┬────────┘
                            │
                            ▼
                   ┌─────────────────┐
                   │  Deterministic  │
                   │      CVC        │
                   └─────────────────┘
```

### Authentication Flow (API Server)

```
Client Request
    │
    ├─> Extract Authorization header
    │
    ├─> Validate Bearer token
    │   │
    │   ├─> Match? → Continue
    │   │
    │   └─> No match? → 401 Unauthorized
    │
    ├─> Check rate limit (per IP)
    │   │
    │   ├─> Within limit? → Continue
    │   │
    │   └─> Exceeded? → 429 Too Many Requests
    │
    └─> Process request
```

## Testing Architecture

### Test Pyramid

```
                      ┌──────────────┐
                      │  End-to-End  │  (Minimal)
                      │  (CLI tests) │
                      └──────────────┘
                    ┌──────────────────┐
                    │   Integration    │  (Medium)
                    │   Tests (test/)  │
                    └──────────────────┘
              ┌────────────────────────────┐
              │     Unit Tests             │  (Majority)
              │  (internal/*_test.go)      │
              └────────────────────────────┘
```

### Test Coverage Strategy

1. **Unit Tests (76-100% per package)**
   - All public functions
   - Edge cases and error paths
   - Luhn validation
   - CVC generation determinism

2. **Integration Tests (66%+)**
   - End-to-end workflows
   - File I/O operations
   - Format conversions
   - CVC recalculation verification

3. **Manual Tests**
   - CLI usability
   - API endpoints
   - Docker deployment

## Performance Considerations

### Optimizations

1. **Cryptographic Efficiency**
   - HMAC-SHA256 is computationally efficient
   - Single hash per CVC (~1-5 microseconds)

2. **Randomness**
   - `crypto/rand` for secure randomness
   - Minimal overhead for PAN generation

3. **Memory**
   - Streaming NDJSON for large datasets
   - Batch generation with bounded memory

4. **Concurrency**
   - Currently sequential (sufficient for testing)
   - Future: Parallel generation for 10k+ cards

### Benchmarks

```
BenchmarkGeneratePAN-8              100000      10000 ns/op
BenchmarkValidateLuhn-8            1000000       1200 ns/op
BenchmarkGenerateDeterministicCVC-8  50000      25000 ns/op
```

## Deployment Architectures

### Standalone CLI

```
Developer Machine
    │
    ├─> Install binary
    │
    ├─> Set CARDGEN_SECRET
    │
    └─> Run commands (generate, transform, validate)
```

### Docker Container

```
Docker Host
    │
    ├─> docker run cardgen-pro generate ...
    │
    └─> Mount volumes for input/output files
```

### API Server (Kubernetes)

```
┌─────────────────────────────────────────┐
│           Kubernetes Cluster            │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │      cardgen-pro Deployment       │ │
│  │                                   │ │
│  │  ┌─────────┐     ┌─────────┐    │ │
│  │  │  Pod 1  │     │  Pod 2  │    │ │
│  │  │ (serve) │     │ (serve) │    │ │
│  │  └────┬────┘     └────┬────┘    │ │
│  └───────┼───────────────┼──────────┘ │
│          │               │            │
│  ┌───────┴───────────────┴──────────┐ │
│  │         Service (ClusterIP)      │ │
│  └───────────────┬──────────────────┘ │
│                  │                    │
│  ┌───────────────┴──────────────────┐ │
│  │          Ingress (HTTPS)         │ │
│  └───────────────┬──────────────────┘ │
└──────────────────┼────────────────────┘
                   │
                   ▼
              External Users
```

## Extensibility Points

### Adding New Card Brands

1. Update `CardBrands` map in `generator.go`
2. Add BIN ranges and configuration
3. Add tests in `generator_test.go`
4. Update documentation

### Adding New ISO-8583 Fields

1. Update `GenerateISO8583Fields` in `iso8583.go`
2. Add field descriptions in comments
3. Add validation tests
4. Update SPEC.md

### Adding New Output Formats

1. Implement writer function in `transformer.go`
2. Add format option to CLI
3. Add tests for new format
4. Update README

## Design Decisions & Trade-offs

### Why HMAC-SHA256 for CVC?

**Decision:** Use HMAC-SHA256 for deterministic CVC generation

**Rationale:**
- ✅ Deterministic (same input = same output)
- ✅ Cryptographically strong
- ✅ Secret-based (prevents trivial generation)
- ✅ Standard, well-audited algorithm

**Trade-offs:**
- Real CVVs use issuer HSMs with proprietary algorithms
- This is a simulation, not production-grade
- Sufficient for testing purposes

### Why Simplified ISO-8583?

**Decision:** Implement subset of ISO-8583 fields without bitmap encoding

**Rationale:**
- ✅ 90% of test cases use ~15 fields
- ✅ Simpler to understand and debug
- ✅ Faster development
- ✅ Sufficient for sandbox testing

**Trade-offs:**
- Not a full ISO-8583 implementation
- No binary bitmap encoding
- Cannot replace production ISO libraries

### Why Go?

**Decision:** Implement in Go language

**Rationale:**
- ✅ Fast compilation and execution
- ✅ Excellent standard library (crypto, HTTP)
- ✅ Single binary deployment
- ✅ Strong typing and error handling
- ✅ Great testing support

**Trade-offs:**
- Learning curve for non-Go developers
- Verbose error handling
- Not as expressive as some languages

## Future Enhancements

### Roadmap (v2.0)

- [ ] EMV/chip data generation (tag-length-value format)
- [ ] Additional card brands (Discover, Diners, JCB)
- [ ] 3D Secure data simulation
- [ ] Tokenization workflow support
- [ ] Database persistence (PostgreSQL)
- [ ] GraphQL API option
- [ ] Prometheus metrics
- [ ] OpenTelemetry tracing

### Contributions Welcome

See [CONTRIBUTING.md](../CONTRIBUTING.md) for how to contribute.

---

**Architecture Version:** 1.0.0  
**Last Updated:** 2025-10-21
