# QA Checklist for cardgen-pro v1.0.0

## 📋 Overview

This checklist should be completed before approving the release of cardgen-pro v1.0.0. Each section contains verification steps with expected results.

**Reviewer:** _____________  
**Date:** _____________  
**Environment:** _____________  
**Go Version:** _____________  

---

## ✅ 1. Build & Installation

### 1.1 Build from Source

- [ ] Clone repository successfully
  ```bash
  git clone https://github.com/felipemacedo1/cardgen-pro.git
  cd cardgen-pro
  ```
  **Expected:** Clean checkout, no errors

- [ ] Dependencies download successfully
  ```bash
  go mod download
  ```
  **Expected:** All modules downloaded, no errors

- [ ] Build succeeds
  ```bash
  go build -o cardgen-pro ./cmd/cardgen-pro
  ```
  **Expected:** Binary created, exit code 0

- [ ] Binary executes
  ```bash
  ./cardgen-pro version
  ```
  **Expected:** Output "cardgen-pro version 1.0.0"

### 1.2 Docker Build

- [ ] Docker image builds
  ```bash
  docker build -t cardgen-pro:test .
  ```
  **Expected:** Image built successfully

- [ ] Docker container runs
  ```bash
  docker run --rm cardgen-pro:test version
  ```
  **Expected:** Output "cardgen-pro version 1.0.0"

---

## ✅ 2. Core Functionality Tests

### 2.1 Luhn Validation

- [ ] Valid PAN passes
  ```bash
  ./cardgen-pro validate 4000000000000002
  ```
  **Expected:** "✓ Valid: 400000******0002 is a valid PAN (Luhn check passed)"

- [ ] Invalid PAN fails
  ```bash
  ./cardgen-pro validate 4000000000000001
  ```
  **Expected:** "✗ Invalid: 400000******0001 failed Luhn check", exit code 1

- [ ] Masked PAN displayed
  **Expected:** PANs masked as first6****last4

### 2.2 Card Generation

- [ ] Generate Visa cards
  ```bash
  export CARDGEN_SECRET="test-qa-secret"
  ./cardgen-pro generate --brand visa --count 5 --out /tmp/visa.json
  ```
  **Expected:** 5 cards in `/tmp/visa.json`, all with Luhn-valid PANs

- [ ] Generate Mastercard cards
  ```bash
  ./cardgen-pro generate --brand mastercard --count 5 --out /tmp/mc.json
  ```
  **Expected:** 5 Mastercard cards, BIN starts with 51-55 or 2221-2720

- [ ] Generate Amex cards
  ```bash
  ./cardgen-pro generate --brand amex --count 3 --out /tmp/amex.json
  ```
  **Expected:** 3 Amex cards, 15 digits, 4-digit CVC

- [ ] Generate with custom BIN
  ```bash
  ./cardgen-pro generate --bin 400000 --count 3 --out /tmp/custom.json
  ```
  **Expected:** All PANs start with 400000

- [ ] Generate with ISO fields
  ```bash
  ./cardgen-pro generate --brand visa --count 2 --iso --out /tmp/iso.json
  ```
  **Expected:** Cards include `iso_fields` with field 2, 3, 4, etc.

- [ ] Generate with Track2
  ```bash
  ./cardgen-pro generate --brand visa --count 2 --track2 --out /tmp/track2.json
  ```
  **Expected:** Cards include `track2` field with format PAN=YYMM...

### 2.3 CVC Generation

- [ ] CVC generated with secret
  ```bash
  export CARDGEN_SECRET="test-secret"
  ./cardgen-pro generate --count 1
  ```
  **Expected:** Card has `cvc` field populated

- [ ] No CVC without secret
  ```bash
  unset CARDGEN_SECRET
  ./cardgen-pro generate --count 1
  ```
  **Expected:** Warning displayed, no CVC in output

- [ ] CVC is deterministic
  ```bash
  export CARDGEN_SECRET="test-secret"
  ./cardgen-pro generate --bin 400000 --count 1 --out /tmp/test1.json
  ./cardgen-pro generate --bin 400000 --count 1 --out /tmp/test2.json
  # Compare CVCs for same PAN
  ```
  **Expected:** Same PAN → Same CVC (if expiry matches)

- [ ] Different secret → Different CVC
  ```bash
  export CARDGEN_SECRET="secret1"
  ./cardgen-pro generate --count 1 > /tmp/cvc1.json
  export CARDGEN_SECRET="secret2"
  ./cardgen-pro generate --count 1 > /tmp/cvc2.json
  ```
  **Expected:** Different CVCs for same PAN

### 2.4 Output Formats

- [ ] JSON format (pretty-printed)
  ```bash
  ./cardgen-pro generate --count 2 --format json --out /tmp/cards.json
  cat /tmp/cards.json
  ```
  **Expected:** Pretty JSON array with indentation

- [ ] NDJSON format (newline-delimited)
  ```bash
  ./cardgen-pro generate --count 2 --format ndjson --out /tmp/cards.ndjson
  cat /tmp/cards.ndjson
  ```
  **Expected:** One JSON object per line, no array brackets

- [ ] CSV format
  ```bash
  ./cardgen-pro generate --count 2 --format csv --out /tmp/cards.csv
  cat /tmp/cards.csv
  ```
  **Expected:** Header row + data rows, comma-separated

### 2.5 Transform Command

- [ ] Transform adds CVCs
  ```bash
  export CARDGEN_SECRET="transform-secret"
  # Create sample orders (no CVCs)
  echo '[{"id":"ORD001","pan":"4000000000000002","expiry_month":12,"expiry_year":2027,"amount":10000,"currency":"986"}]' > /tmp/orders_in.json
  ./cardgen-pro transform --input /tmp/orders_in.json --output /tmp/orders_out.json
  ```
  **Expected:** Output file has CVCs added to each order

- [ ] Transform preserves data
  **Expected:** All original fields (ID, amount, etc.) preserved

- [ ] Transform requires secret
  ```bash
  unset CARDGEN_SECRET
  ./cardgen-pro transform --input /tmp/orders_in.json --output /tmp/out.json
  ```
  **Expected:** Error "Secret is required"

### 2.6 Scenarios Command

- [ ] Lists all scenarios
  ```bash
  ./cardgen-pro scenarios
  ```
  **Expected:** 12 scenarios listed with ID, name, description, response code

- [ ] Scenarios include metadata
  **Expected:** Each scenario has expected_outcome and relevant fields

---

## ✅ 3. HTTP API Tests

### 3.1 Server Start

- [ ] Server starts successfully
  ```bash
  export CARDGEN_SECRET="api-secret"
  ./cardgen-pro serve --port 8080 --token "test-token" &
  PID=$!
  sleep 2
  ```
  **Expected:** Server logs indicate startup on port 8080

- [ ] Health endpoint works (no auth)
  ```bash
  curl http://localhost:8080/health
  ```
  **Expected:** `{"status":"ok","time":"..."}`

### 3.2 Authentication

- [ ] Protected endpoint requires auth
  ```bash
  curl http://localhost:8080/v1/cards
  ```
  **Expected:** 401 Unauthorized

- [ ] Valid token works
  ```bash
  curl -H "Authorization: Bearer test-token" http://localhost:8080/v1/cards
  ```
  **Expected:** 200 OK with cards array

- [ ] Invalid token fails
  ```bash
  curl -H "Authorization: Bearer wrong-token" http://localhost:8080/v1/cards
  ```
  **Expected:** 401 Unauthorized

### 3.3 Card Generation API

- [ ] Generate cards via API
  ```bash
  curl -H "Authorization: Bearer test-token" \
    "http://localhost:8080/v1/cards?brand=visa&count=3"
  ```
  **Expected:** JSON with 3 Visa cards

- [ ] API respects count parameter
  ```bash
  curl -H "Authorization: Bearer test-token" \
    "http://localhost:8080/v1/cards?count=5" | jq '.count'
  ```
  **Expected:** Output "5"

- [ ] API respects brand parameter
  ```bash
  curl -H "Authorization: Bearer test-token" \
    "http://localhost:8080/v1/cards?brand=mastercard&count=2" | jq '.cards[0].brand'
  ```
  **Expected:** Output "Mastercard"

### 3.4 Scenarios API

- [ ] Scenarios endpoint works
  ```bash
  curl -H "Authorization: Bearer test-token" \
    http://localhost:8080/v1/scenarios | jq 'length'
  ```
  **Expected:** Output "12"

### 3.5 Rate Limiting

- [ ] Rate limit enforced (100 req/min)
  ```bash
  for i in {1..101}; do
    curl -s -H "Authorization: Bearer test-token" \
      http://localhost:8080/v1/cards > /dev/null
  done
  ```
  **Expected:** Last request returns 429 Too Many Requests

### 3.6 Server Cleanup

- [ ] Stop server
  ```bash
  kill $PID
  ```
  **Expected:** Server stops gracefully

---

## ✅ 4. Automated Tests

### 4.1 Unit Tests

- [ ] All unit tests pass
  ```bash
  go test ./internal/... -v
  ```
  **Expected:** All tests PASS, no failures

- [ ] Generator tests pass
  ```bash
  go test ./internal/generator -v
  ```
  **Expected:** All tests PASS (Luhn, PAN, CVC, Track2, etc.)

- [ ] ISO tests pass
  ```bash
  go test ./internal/iso -v
  ```
  **Expected:** All tests PASS (100% coverage)

### 4.2 Integration Tests

- [ ] Integration tests pass
  ```bash
  export CARDGEN_SECRET="integration-test-secret"
  go test ./test -v
  ```
  **Expected:** All integration tests PASS

- [ ] CVC determinism test passes
  **Expected:** Same inputs produce same CVCs over multiple runs

- [ ] Transform workflow test passes
  **Expected:** Generate → Transform → Verify workflow succeeds

### 4.3 Coverage

- [ ] Coverage meets threshold (≥80%)
  ```bash
  go test ./... -cover
  ```
  **Expected:** Overall coverage ≥80%

- [ ] Generate coverage report
  ```bash
  go test ./... -coverprofile=coverage.out
  go tool cover -html=coverage.out -o coverage.html
  open coverage.html
  ```
  **Expected:** HTML report shows good coverage

---

## ✅ 5. Code Quality

### 5.1 Formatting

- [ ] Code is gofmt'd
  ```bash
  gofmt -l .
  ```
  **Expected:** No output (all files formatted)

- [ ] No vet warnings
  ```bash
  go vet ./...
  ```
  **Expected:** No warnings

### 5.2 Linting

- [ ] golangci-lint passes
  ```bash
  golangci-lint run
  ```
  **Expected:** No errors or warnings

---

## ✅ 6. Security Verification

### 6.1 No Hardcoded Secrets

- [ ] No secrets in source code
  ```bash
  grep -r "secret.*=" . --exclude-dir=.git | grep -v "CARDGEN_SECRET"
  ```
  **Expected:** No hardcoded secrets found

- [ ] Secret read from environment only
  ```bash
  grep -r "os.Getenv.*SECRET" internal/
  ```
  **Expected:** All secrets from environment variables

### 6.2 PAN Masking

- [ ] PANs masked in logs
  ```bash
  # Run generate and check logs
  ./cardgen-pro generate --count 1 2>&1 | grep "PAN"
  ```
  **Expected:** Only masked PANs in logs (400000******0002)

- [ ] MaskPAN function works
  ```bash
  # Check in code or tests
  grep -A 5 "func MaskPAN" internal/generator/generator.go
  ```
  **Expected:** Masking logic present and correct

### 6.3 Security Documentation

- [ ] SECURITY.md exists and complete
  **Expected:** File present with 395+ lines

- [ ] README has security warnings
  **Expected:** Prominent "TEST/SANDBOX ONLY" warnings

- [ ] PCI-DSS notice present
  **Expected:** Disclaimer about PCI-DSS requirements

---

## ✅ 7. Documentation Verification

### 7.1 README

- [ ] README.md exists and complete
  **Expected:** 470+ lines

- [ ] Quick start examples work
  **Expected:** All commands in README execute successfully

- [ ] Installation instructions accurate
  **Expected:** Can install by following instructions

### 7.2 Technical Docs

- [ ] SPEC.md complete
  **Expected:** Technical specifications documented

- [ ] ARCHITECTURE.md complete
  **Expected:** Architecture diagrams and explanations

- [ ] API.md complete
  **Expected:** All endpoints documented with examples

### 7.3 Contributor Docs

- [ ] CONTRIBUTING.md complete
  **Expected:** Contribution guidelines clear

- [ ] MAINTAINER_GUIDE.md complete
  **Expected:** Release process documented

- [ ] CHANGELOG.md complete
  **Expected:** v1.0.0 entry with all changes

### 7.4 Documentation Links

- [ ] No broken links in docs
  ```bash
  # Check manually or use link checker
  grep -r "\[.*\](" *.md docs/*.md
  ```
  **Expected:** All links valid

---

## ✅ 8. CI/CD Verification

### 8.1 GitHub Actions

- [ ] CI workflow exists
  **Expected:** `.github/workflows/ci.yml` present

- [ ] Release workflow exists
  **Expected:** `.github/workflows/release.yml` present

- [ ] CI runs on push/PR
  **Expected:** Workflows configured correctly

### 8.2 Build Artifacts

- [ ] Multiple platform builds configured
  **Expected:** Linux, macOS, Windows (amd64, arm64)

- [ ] Docker multi-platform build configured
  **Expected:** linux/amd64, linux/arm64

---

## ✅ 9. Release Artifacts

### 9.1 Files Present

- [ ] README.md
- [ ] LICENSE (MIT)
- [ ] SECURITY.md
- [ ] CHANGELOG.md
- [ ] CONTRIBUTING.md
- [ ] go.mod / go.sum
- [ ] Dockerfile
- [ ] .gitignore
- [ ] .golangci.yml
- [ ] All source files
- [ ] All test files
- [ ] All documentation files
- [ ] Fixture samples

### 9.2 Version Consistency

- [ ] main.go version = "1.0.0"
- [ ] CHANGELOG.md has [1.0.0] entry
- [ ] Release notes match version
- [ ] Git tag will be v1.0.0

---

## ✅ 10. Final Verification

### 10.1 End-to-End Workflow

- [ ] Complete user workflow
  ```bash
  # 1. Set secret
  export CARDGEN_SECRET="e2e-test-secret"
  
  # 2. Generate cards
  ./cardgen-pro generate --brand visa --count 5 --iso --track2 --out /tmp/cards.json
  
  # 3. Validate a card
  PAN=$(jq -r '.[0].pan' /tmp/cards.json)
  ./cardgen-pro validate $PAN
  
  # 4. Create orders
  echo '[{"id":"ORD001","pan":"'$PAN'","expiry_month":12,"expiry_year":2027,"amount":10000,"currency":"986"}]' > /tmp/orders.json
  
  # 5. Transform orders
  ./cardgen-pro transform --input /tmp/orders.json --output /tmp/orders_cvc.json
  
  # 6. Verify CVC matches
  # Compare CVC in cards.json vs orders_cvc.json
  
  # 7. Start API
  ./cardgen-pro serve --port 8080 --token "e2e-token" &
  sleep 2
  
  # 8. Query API
  curl -H "Authorization: Bearer e2e-token" \
    "http://localhost:8080/v1/cards?count=3"
  
  # 9. Stop API
  killall cardgen-pro
  ```
  **Expected:** All steps succeed, CVCs match

### 10.2 Clean Up

- [ ] Remove test files
  ```bash
  rm -f /tmp/cards.json /tmp/orders*.json /tmp/visa.json /tmp/mc.json /tmp/amex.json
  ```

---

## 📊 QA Summary

### Test Results

| Category | Total | Passed | Failed | Coverage |
|----------|-------|--------|--------|----------|
| Build & Installation | 5 | ___ | ___ | N/A |
| Core Functionality | 25 | ___ | ___ | N/A |
| HTTP API | 13 | ___ | ___ | N/A |
| Automated Tests | 6 | ___ | ___ | ≥80% |
| Code Quality | 3 | ___ | ___ | N/A |
| Security | 8 | ___ | ___ | N/A |
| Documentation | 8 | ___ | ___ | N/A |
| CI/CD | 4 | ___ | ___ | N/A |
| Release Artifacts | 16 | ___ | ___ | N/A |
| Final Verification | 2 | ___ | ___ | N/A |
| **TOTAL** | **90** | ___ | ___ | ___ |

### Pass Criteria

- [ ] All critical tests pass (Build, Core, Security)
- [ ] ≥90% of total tests pass
- [ ] Test coverage ≥80%
- [ ] No hardcoded secrets
- [ ] Documentation complete
- [ ] Security warnings prominent

### Reviewer Sign-Off

**Approved for Release:** ☐ Yes ☐ No

**Reviewer Signature:** _________________________

**Date:** _________________________

**Notes:**
```
_________________________________________________________________
_________________________________________________________________
_________________________________________________________________
```

---

## 📝 Issue Tracking

If any tests fail, log them here:

| Test # | Category | Issue | Severity | Status |
|--------|----------|-------|----------|--------|
| | | | | |
| | | | | |
| | | | | |

**Severities:** Critical, High, Medium, Low

---

**QA Checklist for cardgen-pro v1.0.0**  
**Last Updated:** 2025-10-21
