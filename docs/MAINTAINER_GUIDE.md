# Maintainer Guide

This guide is for project maintainers responsible for releases, dependency updates, and project governance.

## Release Process

### Semantic Versioning

We follow [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR** version: Incompatible API changes
- **MINOR** version: Backward-compatible functionality additions
- **PATCH** version: Backward-compatible bug fixes

### Pre-Release Checklist

Before creating a new release:

- [ ] All tests passing (`go test ./...`)
- [ ] Code coverage ≥80%
- [ ] golangci-lint passes with no warnings
- [ ] CHANGELOG.md updated with version and changes
- [ ] README.md version badges updated (if needed)
- [ ] Documentation reviewed and updated
- [ ] Security vulnerabilities addressed
- [ ] Dependencies updated to stable versions

### Release Steps

#### 1. Update Version

Update version string in `cmd/cardgen-pro/main.go`:

```go
const version = "1.1.0"  // Update this
```

#### 2. Update CHANGELOG.md

Add new version section:

```markdown
## [1.1.0] - 2025-11-01

### Added
- New feature description

### Changed
- Changed functionality description

### Fixed
- Bug fix description

### Security
- Security fix description
```

#### 3. Commit Changes

```bash
git add .
git commit -m "chore(release): prepare v1.1.0"
git push origin main
```

#### 4. Create Git Tag

```bash
# Create annotated tag
git tag -a v1.1.0 -m "Release v1.1.0

- Feature 1
- Feature 2
- Bug fix 3
"

# Push tag
git push origin v1.1.0
```

#### 5. GitHub Release

GitHub Actions will automatically:
- Build binaries for multiple platforms
- Run tests
- Create GitHub Release
- Upload artifacts
- Build and push Docker images

Verify release at: https://github.com/felipemacedo1/cardgen-pro/releases

#### 6. Announce Release

- Update project website (if exists)
- Announce on social media
- Notify users via mailing list
- Update package managers (if applicable)

### Hotfix Process

For critical bugs in production releases:

1. Create hotfix branch: `git checkout -b hotfix/v1.0.1`
2. Fix bug and add tests
3. Update CHANGELOG.md
4. Update version to patch (e.g., 1.0.0 → 1.0.1)
5. Commit: `git commit -m "fix: critical bug description"`
6. Merge to main: `git checkout main && git merge hotfix/v1.0.1`
7. Tag and release: `git tag v1.0.1 && git push origin v1.0.1`
8. Delete hotfix branch: `git branch -d hotfix/v1.0.1`

## Dependency Management

### Updating Dependencies

```bash
# Check for outdated dependencies
go list -u -m all

# Update specific dependency
go get -u github.com/package/name@version

# Update all dependencies
go get -u ./...

# Tidy and verify
go mod tidy
go mod verify

# Run tests
go test ./...
```

### Security Audits

```bash
# Check for known vulnerabilities
go list -json -m all | nancy sleuth

# Or use govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

### Dependency Review

Before accepting dependency updates:

- [ ] Review changelog for breaking changes
- [ ] Check security advisories
- [ ] Verify license compatibility (MIT)
- [ ] Run full test suite
- [ ] Test CLI functionality manually

## CI/CD Maintenance

### GitHub Actions

**Workflows:**
- `.github/workflows/ci.yml` - CI pipeline (tests, lint, build)
- `.github/workflows/release.yml` - Release automation

**Maintenance Tasks:**

1. **Update Go version** (every 6 months):
   ```yaml
   - uses: actions/setup-go@v5
     with:
       go-version: '1.23'  # Update this
   ```

2. **Update action versions** (quarterly):
   ```yaml
   - uses: actions/checkout@v4  # Check for v5
   - uses: actions/cache@v4
   ```

3. **Review security alerts** (weekly):
   - Check GitHub Security tab
   - Update vulnerable dependencies

4. **Monitor CI performance** (monthly):
   - Check build times
   - Optimize slow tests
   - Review cache effectiveness

## Code Quality Maintenance

### Linter Configuration

Update `.golangci.yml` as needed:

```yaml
linters:
  enable:
    - gofmt
    - govet
    - staticcheck
    # Add new linters here
```

### Test Coverage

Maintain ≥80% coverage:

```bash
# Generate coverage report
go test ./... -coverprofile=coverage.out

# View coverage
go tool cover -html=coverage.out

# Check total coverage
go tool cover -func=coverage.out | grep total
```

### Code Review Guidelines

For all PRs:

- [ ] Tests included and passing
- [ ] Documentation updated
- [ ] No hardcoded secrets
- [ ] Error handling appropriate
- [ ] Code follows style guidelines
- [ ] Commit messages follow conventions

## Security Maintenance

### Vulnerability Disclosure

If vulnerability reported:

1. **Acknowledge** within 24 hours
2. **Assess** severity (Critical/High/Medium/Low)
3. **Fix** in private fork
4. **Test** thoroughly
5. **Release** hotfix
6. **Disclose** with credit to reporter
7. **Update** SECURITY.md if needed

### Secret Rotation

For project secrets:

```bash
# GitHub Actions secrets
- DOCKERHUB_USERNAME
- DOCKERHUB_TOKEN
- GITHUB_TOKEN (automatic)

# Rotate quarterly or on compromise
```

### Security Audits

**Schedule:**
- Weekly: Dependency scan
- Monthly: Code review for security issues
- Quarterly: Full security audit

## Documentation Maintenance

### Keep Updated

- [ ] README.md (features, examples, installation)
- [ ] CHANGELOG.md (all changes)
- [ ] SECURITY.md (security guidelines)
- [ ] CONTRIBUTING.md (contribution process)
- [ ] docs/SPEC.md (technical specifications)
- [ ] docs/ARCHITECTURE.md (design decisions)

### Documentation Review Checklist

- [ ] All commands documented
- [ ] Examples working and tested
- [ ] Links not broken
- [ ] Screenshots current
- [ ] API docs match implementation

## Community Management

### Issue Triage

**Labels:**
- `bug` - Something isn't working
- `enhancement` - New feature request
- `documentation` - Documentation improvements
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention needed
- `security` - Security-related
- `wontfix` - Will not be addressed

**Triage Process:**
1. Label appropriately
2. Ask for clarification if needed
3. Close duplicates
4. Assign to milestone (if accepted)
5. Respond within 48 hours

### Pull Request Review

**Review Checklist:**
- [ ] CI passing
- [ ] Tests included
- [ ] Documentation updated
- [ ] Follows style guide
- [ ] No breaking changes (or documented)
- [ ] Commit messages clear

**Merge Process:**
1. Review code
2. Request changes if needed
3. Approve when ready
4. Squash and merge (or rebase)
5. Delete branch
6. Close related issues

## Monitoring & Metrics

### Track These Metrics

- **Test Coverage:** Should stay ≥80%
- **Build Time:** Should be <5 minutes
- **Issue Response Time:** Target <48 hours
- **PR Merge Time:** Target <1 week
- **Release Frequency:** Target monthly (minor)

### Health Indicators

**Green (Healthy):**
- All tests passing
- Coverage ≥85%
- No critical vulnerabilities
- <10 open bugs

**Yellow (Attention Needed):**
- 1-2 test failures
- Coverage 75-84%
- 1-2 high vulnerabilities
- 10-20 open bugs

**Red (Action Required):**
- Multiple test failures
- Coverage <75%
- Critical vulnerabilities
- >20 open bugs

## Offboarding

If stepping down as maintainer:

1. **Document** institutional knowledge
2. **Train** new maintainers
3. **Transfer** access (GitHub, Docker Hub)
4. **Update** MAINTAINERS.md
5. **Announce** transition to community

## Emergency Procedures

### Critical Bug in Production

1. **Assess** impact and affected versions
2. **Create** hotfix branch immediately
3. **Fix** and test thoroughly
4. **Release** patch version ASAP
5. **Notify** users via GitHub issue
6. **Post-mortem** after resolution

### Security Breach

1. **Revoke** compromised secrets immediately
2. **Assess** damage and data exposure
3. **Notify** affected users
4. **Fix** vulnerability
5. **Release** security patch
6. **Document** in security advisory

### CI/CD Failure

1. **Check** GitHub Actions status
2. **Review** recent changes
3. **Rollback** if needed
4. **Fix** and redeploy
5. **Monitor** for recurrence

## Useful Commands

```bash
# Run all tests with race detection
go test -race ./...

# Check for data races
go test -race -count=1000 ./internal/generator

# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Profile memory usage
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# Generate test coverage badge
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total | awk '{print $3}'

# Build for all platforms
GOOS=linux GOARCH=amd64 go build -o cardgen-pro-linux-amd64 ./cmd/cardgen-pro
GOOS=darwin GOARCH=amd64 go build -o cardgen-pro-darwin-amd64 ./cmd/cardgen-pro
GOOS=windows GOARCH=amd64 go build -o cardgen-pro-windows-amd64.exe ./cmd/cardgen-pro
```

## Resources

- [Go Release Policy](https://golang.org/doc/devel/release.html)
- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [GitHub Actions Docs](https://docs.github.com/en/actions)

## Contact

For maintainer questions:
- GitHub Discussions: [Project Discussions](https://github.com/felipemacedo1/cardgen-pro/discussions)
- Email: maintainer@cardgen-pro.example.com

---

**Last Updated:** 2025-10-21  
**Maintainer:** Felipe Macedo
