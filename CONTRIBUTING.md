# Contributing to cardgen-pro

Thank you for your interest in contributing to **cardgen-pro**! This document provides guidelines for contributing to the project.

## Code of Conduct

- Be respectful and professional
- Focus on constructive feedback
- Help maintain a welcoming environment
- Remember: this is a testing tool for sandbox environments

## How to Contribute

### Reporting Issues

Before creating an issue, please:

1. Check if the issue already exists
2. Provide a clear title and description
3. Include steps to reproduce (if bug)
4. Include expected vs actual behavior
5. Add relevant logs/screenshots

### Suggesting Features

Feature requests are welcome! Please:

1. Clearly describe the use case
2. Explain why it's valuable for testing
3. Provide examples of how it would be used
4. Consider backward compatibility

### Pull Requests

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/my-feature`
3. **Make** your changes
4. **Test** thoroughly (see Testing section)
5. **Commit** with conventional commits (see below)
6. **Push** to your fork
7. **Open** a pull request

## Development Setup

```bash
# Clone repository
git clone https://github.com/felipemacedo1/cardgen-pro.git
cd cardgen-pro

# Install dependencies
go mod download

# Run tests
go test ./... -v

# Run linter
golangci-lint run

# Build
go build -o cardgen-pro ./cmd/cardgen-pro
```

## Code Style

### Go Guidelines

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use meaningful variable names
- Add comments for exported functions
- Keep functions focused and small
- Use error wrapping: `fmt.Errorf("context: %w", err)`

### Comment Style

```go
// GeneratePAN generates a valid PAN using Luhn algorithm.
// BIN: Bank Identification Number (first 6 digits)
// length: total PAN length (13-19)
func GeneratePAN(bin string, length int) (string, error) {
    // Implementation
}
```

### Error Handling

```go
// ✅ Good
if err != nil {
    return fmt.Errorf("failed to generate PAN: %w", err)
}

// ❌ Bad
if err != nil {
    return err  // Lost context
}
```

## Testing

### Running Tests

```bash
# Unit tests
go test ./internal/... -v

# Integration tests
go test ./test/... -v

# All tests with coverage
go test ./... -cover -coverprofile=coverage.out

# View coverage
go tool cover -html=coverage.out
```

### Writing Tests

```go
func TestGeneratePAN(t *testing.T) {
    tests := []struct {
        name      string
        bin       string
        length    int
        shouldErr bool
    }{
        {"Valid Visa 16", "400000", 16, false},
        {"Invalid length", "400000", 10, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            pan, err := GeneratePAN(tt.bin, tt.length)
            
            if tt.shouldErr {
                if err == nil {
                    t.Error("Expected error but got none")
                }
                return
            }

            if err != nil {
                t.Errorf("Unexpected error: %v", err)
            }

            if !ValidateLuhn(pan) {
                t.Error("Generated invalid PAN")
            }
        })
    }
}
```

### Test Coverage Requirements

- **Minimum:** 80% coverage
- **Unit tests:** All public functions
- **Integration tests:** End-to-end workflows
- **Edge cases:** Boundary conditions, error paths

## Commit Messages

Use **Conventional Commits** format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test additions/changes
- `refactor`: Code refactoring
- `chore`: Maintenance tasks
- `ci`: CI/CD changes

### Examples

```bash
feat(generator): add support for Discover cards

- Add Discover BIN ranges
- Update validation logic
- Add tests for Discover generation

Closes #42

---

fix(api): correct rate limiting calculation

Rate limiter was not properly cleaning old requests,
causing false positives on limit exceeded.

---

docs(readme): add Docker deployment instructions

---

test(integration): add CVC determinism test
```

## Pull Request Process

### Before Submitting

- [ ] Code follows style guidelines
- [ ] Tests pass locally
- [ ] Coverage is ≥80%
- [ ] Documentation updated (if needed)
- [ ] Commit messages follow conventions
- [ ] No secrets or sensitive data committed

### PR Template

```markdown
## Description

Brief description of changes

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing

- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed

## Checklist

- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex logic
- [ ] Documentation updated
- [ ] Tests pass
- [ ] No breaking changes (or documented)
```

### Review Process

1. Automated checks run (CI)
2. Maintainer reviews code
3. Feedback addressed
4. Approved and merged

## Project Structure

```
cardgen-pro/
├── cmd/cardgen-pro/        # CLI entry point
│   └── main.go
├── internal/               # Private packages
│   ├── generator/          # Card generation logic
│   │   ├── generator.go
│   │   ├── luhn.go
│   │   └── generator_test.go
│   ├── iso/                # ISO-8583 helpers
│   ├── api/                # HTTP API
│   └── models/             # Data structures
├── pkg/transformer/        # Public packages
├── test/                   # Integration tests
├── fixtures/               # Test fixtures
├── docs/                   # Documentation
└── README.md
```

## Adding New Features

### Example: Adding a New Card Brand

1. **Update models** (`internal/models/card.go`):
   ```go
   // Add BIN range to CardBrands map
   ```

2. **Add tests** (`internal/generator/generator_test.go`):
   ```go
   func TestGenerateNewBrand(t *testing.T) { ... }
   ```

3. **Update documentation**:
   - README.md (supported brands)
   - SPEC.md (BIN ranges)

4. **Add fixtures** (`fixtures/new_brand_examples.json`)

5. **Submit PR** with description and examples

## Release Process

See [MAINTAINER_GUIDE.md](./docs/MAINTAINER_GUIDE.md) for release procedures.

## Getting Help

- **GitHub Issues:** Ask questions or report bugs
- **Discussions:** General questions and ideas
- **Email:** maintainer@example.com (if set up)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to cardgen-pro!** 🎉
