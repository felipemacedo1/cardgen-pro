# GitHub Secrets Configuration

This document describes the required secrets and variables for CI/CD workflows.

## Required Secrets

### 1. CODECOV_TOKEN (Optional)
**Used by:** CI workflow  
**Purpose:** Upload test coverage to Codecov  
**How to get:**
1. Go to https://codecov.io
2. Sign in with GitHub
3. Add your repository
4. Copy the token from Settings → General

**How to set:**
```
Settings → Secrets and variables → Actions → New repository secret
Name: CODECOV_TOKEN
Value: <your-token>
```

### 2. DOCKERHUB_USERNAME (Optional)
**Used by:** Release workflow  
**Purpose:** Authenticate with Docker Hub to push images  
**Value:** Your Docker Hub username

### 3. DOCKERHUB_TOKEN (Optional)
**Used by:** Release workflow  
**Purpose:** Authenticate with Docker Hub to push images  
**How to get:**
1. Go to https://hub.docker.com
2. Account Settings → Security → New Access Token
3. Copy the token

**How to set:**
```
Settings → Secrets and variables → Actions → New repository secret
Name: DOCKERHUB_USERNAME
Value: <your-dockerhub-username>

Name: DOCKERHUB_TOKEN
Value: <your-access-token>
```

## Repository Variables

### DOCKER_ENABLED
**Used by:** CI and Release workflows  
**Purpose:** Enable/disable Docker build steps  
**Default:** `false`  
**Values:** `true` or `false`

**How to set:**
```
Settings → Secrets and variables → Actions → Variables → New repository variable
Name: DOCKER_ENABLED
Value: false
```

Set to `true` only when:
- ✅ DOCKERHUB_USERNAME is configured
- ✅ DOCKERHUB_TOKEN is configured
- ✅ You want to publish Docker images

## Workflow Behavior

### Without Secrets
If secrets are not configured:
- ✅ **CI workflow:** Runs all tests, skips Codecov upload and Docker steps
- ✅ **Release workflow:** Creates GitHub release with binaries, skips Docker push

### With Secrets
If all secrets are configured and DOCKER_ENABLED=true:
- ✅ **CI workflow:** Runs all tests, uploads coverage, builds Docker images
- ✅ **Release workflow:** Creates GitHub release, pushes Docker images to Docker Hub

## Testing Configuration

After configuring secrets, test with:

```bash
# Trigger CI workflow
git push origin main

# Trigger release workflow
git tag v1.0.1
git push origin v1.0.1
```

## Security Best Practices

1. ✅ **Never commit secrets** to source code
2. ✅ **Rotate tokens** regularly (every 90 days)
3. ✅ **Use minimal permissions** for tokens
4. ✅ **Monitor workflow logs** for exposed secrets
5. ✅ **Enable secret scanning** in repository settings

## Troubleshooting

### "Username and password required"
**Cause:** DOCKERHUB_USERNAME or DOCKERHUB_TOKEN not set  
**Fix:** Either configure the secrets or set DOCKER_ENABLED=false

### "Codecov upload failed"
**Cause:** CODECOV_TOKEN not set  
**Fix:** This is non-critical. Workflow continues but coverage is not uploaded.

### Docker job skipped
**Cause:** DOCKER_ENABLED variable is false or not set  
**Fix:** Set DOCKER_ENABLED=true if you want Docker builds

## References

- [GitHub Secrets Documentation](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- [Codecov Documentation](https://docs.codecov.com/docs)
- [Docker Hub Access Tokens](https://docs.docker.com/docker-hub/access-tokens/)
