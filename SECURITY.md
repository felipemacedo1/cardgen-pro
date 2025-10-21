# Security Guidelines

## ⚠️ Critical Warnings

### Primary Warning

**cardgen-pro generates SYNTHETIC TEST DATA for sandbox/development environments ONLY.**

- ❌ **NEVER** use generated PANs on production payment networks
- ❌ **NEVER** use this tool with real cardholder data
- ❌ **NEVER** bypass PCI-DSS or security requirements with this tool
- ❌ **NEVER** commit secrets to version control
- ❌ **NEVER** log full PANs or CVCs in production systems

### Scope of Use

**Approved Use Cases:**
- ✅ Development environment testing
- ✅ QA/staging sandbox testing
- ✅ Integration testing with mock payment gateways
- ✅ Load testing payment workflows
- ✅ Training and demonstrations

**Prohibited Use Cases:**
- ❌ Production payment processing
- ❌ Real financial transactions
- ❌ Customer-facing systems
- ❌ Compliance testing with real data
- ❌ PCI-DSS scope reduction attempts

## 🔐 Secret Management

### Why Secrets Matter

The `CARDGEN_SECRET` is used for deterministic CVC generation via HMAC-SHA256. While this is synthetic data, proper secret management is a **best practice** that should be maintained in all environments.

### Secret Requirements

- **Length:** Minimum 32 characters (256 bits recommended)
- **Complexity:** Cryptographically random (use `openssl rand -base64 32`)
- **Uniqueness:** Different secret per environment
- **Rotation:** Rotate every 90 days or on suspected compromise

### Providing Secrets

#### Environment Variables (Development)

```bash
# Set in shell session
export CARDGEN_SECRET="your-dev-secret-key-here"

# Verify (DO NOT echo in production scripts)
echo $CARDGEN_SECRET
```

#### AWS Secrets Manager (Recommended for AWS)

```bash
# Store secret
aws secretsmanager create-secret \
  --name cardgen/dev/secret \
  --secret-string "your-dev-secret-key"

# Retrieve in application
export CARDGEN_SECRET=$(aws secretsmanager get-secret-value \
  --secret-id cardgen/dev/secret \
  --query SecretString \
  --output text)
```

#### HashiCorp Vault (Recommended for Multi-Cloud)

```bash
# Store secret
vault kv put secret/cardgen/dev secret="your-dev-secret-key"

# Retrieve in application
export CARDGEN_SECRET=$(vault kv get -field=secret secret/cardgen/dev)
```

#### Kubernetes Secrets

```bash
# Create secret
kubectl create secret generic cardgen-secret \
  --from-literal=secret=your-dev-secret-key

# Use in pod
apiVersion: v1
kind: Pod
metadata:
  name: cardgen-worker
spec:
  containers:
  - name: cardgen
    image: cardgen-pro:latest
    env:
    - name: CARDGEN_SECRET
      valueFrom:
        secretKeyRef:
          name: cardgen-secret
          key: secret
```

#### Docker Secrets (Docker Swarm)

```bash
# Create secret
echo "your-dev-secret-key" | docker secret create cardgen_secret -

# Use in service
docker service create \
  --name cardgen-api \
  --secret cardgen_secret \
  --env CARDGEN_SECRET_FILE=/run/secrets/cardgen_secret \
  cardgen-pro:latest serve
```

### What NOT to Do

```bash
# ❌ NEVER hardcode in source code
const SECRET = "hardcoded-secret" // WRONG!

# ❌ NEVER commit to git
git add .env
git commit -m "Added secrets" // WRONG!

# ❌ NEVER pass as command-line argument (visible in ps)
cardgen-pro generate --secret "my-secret" // WRONG! Use env var

# ❌ NEVER log secrets
log.Printf("Using secret: %s", secret) // WRONG!
```

## 🎭 PAN Masking

### Always Mask PANs

**Never log or display full PANs.** Use masking format: `first6****last4`

```go
import "github.com/felipemacedo/cardgen-pro/internal/generator"

pan := "4000000000000002"
masked := generator.MaskPAN(pan)
// Output: "400000******0002"

// Safe to log
log.Printf("Processing card: %s", masked)
```

### Masking Rules

| Data Type | Masking Format | Example |
|-----------|---------------|---------|
| PAN | `first6****last4` | `400000******0002` |
| CVC | Never log | `***` |
| Expiry | Safe to log | `12/2027` |
| Cardholder Name | First initial + asterisks | `J*** D**` |

## 🗑️ Data Cleanup

### Remove Generated Files

```bash
# After testing, securely delete fixtures
shred -vfz -n 10 cards.json
shred -vfz -n 10 orders_with_cvc.json

# Or use secure delete (macOS)
rm -P cards.json

# Verify deletion
ls -la cards.json # should not exist
```

### Container Cleanup

```bash
# Remove container volumes after testing
docker rm -v cardgen-test-container

# Clean up temporary files in containers
docker run --rm cardgen-pro sh -c "rm -rf /tmp/*"
```

## 🔄 Secret Rotation

### When to Rotate

- **Scheduled:** Every 90 days
- **On compromise:** Immediately if secret is exposed
- **On team changes:** When developers leave
- **On environment promotion:** Different secrets per environment

### Rotation Process

1. Generate new secret: `openssl rand -base64 32`
2. Update secret in secret manager
3. Restart applications with new secret
4. Verify CVC generation still works (deterministic check may fail for old data)
5. Revoke old secret access

### Impact of Rotation

⚠️ **Important:** Rotating the secret will change generated CVCs for the same PAN/expiry.

If you need to **preserve** CVCs across rotation (e.g., for regression tests):
- Store test fixtures with CVCs in version control (masked PANs)
- Use fixture files instead of regenerating

## 🛡️ Access Controls

### Principle of Least Privilege

**Who should have access:**
- ✅ QA engineers (generate test data)
- ✅ Backend developers (integration testing)
- ✅ DevOps/SRE (sandbox infrastructure)

**Who should NOT have access:**
- ❌ Production systems
- ❌ Customer support
- ❌ Third-party vendors (unless specifically approved)

### API Server Security

If running the HTTP API (`cardgen-pro serve`):

```bash
# Generate strong token
TOKEN=$(openssl rand -hex 32)

# Start server with token
cardgen-pro serve --port 8080 --token $TOKEN

# Client authentication
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/v1/cards
```

**Additional Recommendations:**
- Run behind reverse proxy (nginx, Caddy)
- Enable TLS/HTTPS (even in sandbox)
- Implement IP whitelisting
- Use network policies (Kubernetes NetworkPolicy, AWS Security Groups)
- Monitor and alert on suspicious access patterns

## 📊 Logging & Monitoring

### Safe Logging Practices

```go
// ✅ CORRECT: Masked PAN
log.Printf("Generated card: %s", card.MaskedPAN)

// ✅ CORRECT: No sensitive data
log.Printf("Generated %d cards for brand %s", count, brand)

// ❌ WRONG: Full PAN
log.Printf("Generated card: %s", card.PAN)

// ❌ WRONG: CVC
log.Printf("Card CVC: %s", card.CVC)
```

### Monitoring

Monitor for:
- Excessive generation requests (potential abuse)
- Failed CVC generation (secret misconfiguration)
- API 401/403 errors (unauthorized access attempts)
- Large batch generation (>1000 cards) - may indicate misuse

### Audit Trails

For compliance, maintain audit logs:
- Who generated cards (user/service account)
- When (timestamp)
- How many (count)
- Which environment (dev/staging/sandbox)

## 🚨 Incident Response

### If Secret is Compromised

1. **Immediately** rotate the secret
2. Review access logs for unauthorized use
3. Notify security team
4. If data was exposed externally, assess impact (remember: this is synthetic test data)

### If Tool is Misused

1. Investigate scope of misuse
2. Disable access for responsible party
3. Review and strengthen access controls
4. Re-train team on proper use

### If Real PANs are Processed

**This is a CRITICAL incident:**

1. **STOP** using the tool immediately
2. Isolate affected systems
3. Notify PCI-DSS compliance team
4. Follow your organization's data breach response plan
5. Determine if PCI-DSS scope has been impacted
6. Engage legal/compliance as needed

**Remember:** This tool is NOT designed for real cardholder data.

## ✅ Compliance Checklist

Use this checklist for security reviews:

- [ ] Secret stored in secret manager (not hardcoded)
- [ ] Secret rotation policy in place (90 days)
- [ ] PANs always masked in logs
- [ ] CVCs never logged
- [ ] Generated files cleaned up after tests
- [ ] Access controls documented
- [ ] Tool used only in sandbox/dev/QA environments
- [ ] API server uses strong authentication
- [ ] Monitoring and alerting configured
- [ ] Team trained on proper use
- [ ] Incident response plan includes this tool

## 📚 Additional Resources

- [PCI-DSS Requirements](https://www.pcisecuritystandards.org/)
- [OWASP Secrets Management](https://cheatsheetseries.owasp.org/cheatsheets/Secrets_Management_Cheat_Sheet.html)
- [NIST Cryptographic Standards](https://csrc.nist.gov/projects/cryptographic-standards-and-guidelines)
- [AWS Secrets Manager Best Practices](https://docs.aws.amazon.com/secretsmanager/latest/userguide/best-practices.html)

## 🙋 Questions?

For security-related questions or to report vulnerabilities:

- **GitHub Issues:** [Report security issue](https://github.com/felipemacedo1/cardgen-pro/issues/new?labels=security)
- **Email:** security@cardgen-pro.example.com (if you set up a security contact)

---

**Security is everyone's responsibility. Use this tool wisely.**
