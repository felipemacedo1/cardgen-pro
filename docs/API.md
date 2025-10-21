# API Documentation

REST API for **cardgen-pro** fixture serving (sandbox/development only).

## ⚠️ Security Warning

**This API is for TEST/SANDBOX environments only.**
- Never expose publicly
- Use strong authentication tokens
- Implement network restrictions
- Enable TLS/HTTPS in production-like environments

## Starting the Server

```bash
# Basic usage
export CARDGEN_SECRET="your-dev-secret"
cardgen-pro serve --port 8080 --token "your-auth-token"

# With custom port
cardgen-pro serve --port 3000 --token "dev-token-xyz"

# Docker
docker run -p 8080:8080 \
  -e CARDGEN_SECRET="your-dev-secret" \
  cardgen-pro:latest \
  serve --port 8080 --token "your-auth-token"
```

## Authentication

All protected endpoints require Bearer token authentication:

```bash
curl -H "Authorization: Bearer your-auth-token" \
  http://localhost:8080/v1/cards
```

**Unauthorized requests return 401:**

```json
{
  "error": "Unauthorized: invalid token"
}
```

## Rate Limiting

- **Limit:** 100 requests per minute per IP
- **Window:** Rolling 60 seconds
- **Response:** 429 Too Many Requests

```json
{
  "error": "Rate limit exceeded"
}
```

## Endpoints

### Health Check

**Public endpoint** - no authentication required

```http
GET /health
```

**Response: 200 OK**

```json
{
  "status": "ok",
  "time": "2025-10-21T12:00:00Z"
}
```

**Example:**

```bash
curl http://localhost:8080/health
```

---

### Generate Cards

**Protected endpoint** - requires authentication

```http
GET /v1/cards?brand={brand}&count={count}&bin={bin}&secret={secret}
```

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `brand` | string | No | `visa` | Card brand: `visa`, `mastercard`, `amex` |
| `count` | integer | No | `10` | Number of cards (max 100) |
| `bin` | string | No | - | Custom BIN (6 digits) |
| `secret` | string | No | - | CVC generation secret |

**Response: 200 OK**

```json
{
  "cards": [
    {
      "pan": "4000000000000002",
      "masked_pan": "400000******0002",
      "brand": "Visa",
      "expiry_month": 12,
      "expiry_year": 2027,
      "cvc": "382",
      "track2": "4000000000000002=27122011234",
      "iso_fields": {
        "2": "4000000000000002",
        "3": "000000",
        "4": "000000010000",
        "11": "123456",
        "14": "2712",
        "49": "986"
      },
      "generated_at": "2025-10-21T12:00:00Z"
    }
  ],
  "count": 1
}
```

**Examples:**

```bash
# Generate 5 Visa cards
curl -H "Authorization: Bearer your-token" \
  "http://localhost:8080/v1/cards?brand=visa&count=5"

# Generate Mastercard with custom BIN
curl -H "Authorization: Bearer your-token" \
  "http://localhost:8080/v1/cards?brand=mastercard&bin=510000&count=3"

# Generate with CVC secret
curl -H "Authorization: Bearer your-token" \
  "http://localhost:8080/v1/cards?brand=amex&count=2&secret=my-secret"
```

**Error Responses:**

```http
401 Unauthorized
{
  "error": "Unauthorized: missing Authorization header"
}

429 Too Many Requests
{
  "error": "Rate limit exceeded"
}

500 Internal Server Error
{
  "error": "Failed to generate card: ..."
}
```

---

### List Test Scenarios

**Protected endpoint** - requires authentication

```http
GET /v1/scenarios
```

**Response: 200 OK**

```json
[
  {
    "id": "success_auth",
    "name": "Successful Authorization",
    "description": "Standard approved transaction",
    "response_code": "00",
    "response_text": "Approved",
    "amount": 10000,
    "currency": "986",
    "card_brand": "visa",
    "expected_outcome": "Transaction approved, auth code generated"
  },
  {
    "id": "insufficient_funds",
    "name": "Insufficient Funds",
    "description": "Card has insufficient balance",
    "response_code": "51",
    "response_text": "Insufficient funds",
    "amount": 100000,
    "currency": "986",
    "card_brand": "visa",
    "expected_outcome": "Decline due to insufficient balance"
  }
]
```

**Example:**

```bash
curl -H "Authorization: Bearer your-token" \
  http://localhost:8080/v1/scenarios | jq .
```

## Client Examples

### cURL

```bash
#!/bin/bash

TOKEN="your-auth-token"
BASE_URL="http://localhost:8080"

# Generate cards
curl -H "Authorization: Bearer $TOKEN" \
  "$BASE_URL/v1/cards?brand=visa&count=5" \
  | jq -r '.cards[].masked_pan'

# List scenarios
curl -H "Authorization: Bearer $TOKEN" \
  "$BASE_URL/v1/scenarios" \
  | jq -r '.[].id'
```

### Python

```python
import requests

BASE_URL = "http://localhost:8080"
TOKEN = "your-auth-token"
HEADERS = {"Authorization": f"Bearer {TOKEN}"}

# Generate cards
response = requests.get(
    f"{BASE_URL}/v1/cards",
    headers=HEADERS,
    params={"brand": "visa", "count": 5}
)
cards = response.json()["cards"]

for card in cards:
    print(f"PAN: {card['masked_pan']}, CVC: {card['cvc']}")

# List scenarios
response = requests.get(f"{BASE_URL}/v1/scenarios", headers=HEADERS)
scenarios = response.json()

for scenario in scenarios:
    print(f"{scenario['id']}: {scenario['name']}")
```

### JavaScript (Node.js)

```javascript
const axios = require('axios');

const BASE_URL = 'http://localhost:8080';
const TOKEN = 'your-auth-token';
const headers = { Authorization: `Bearer ${TOKEN}` };

// Generate cards
async function generateCards() {
  const response = await axios.get(`${BASE_URL}/v1/cards`, {
    headers,
    params: { brand: 'visa', count: 5 }
  });
  
  response.data.cards.forEach(card => {
    console.log(`PAN: ${card.masked_pan}, CVC: ${card.cvc}`);
  });
}

// List scenarios
async function listScenarios() {
  const response = await axios.get(`${BASE_URL}/v1/scenarios`, { headers });
  
  response.data.forEach(scenario => {
    console.log(`${scenario.id}: ${scenario.name}`);
  });
}

generateCards();
listScenarios();
```

### Go

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

const (
    baseURL = "http://localhost:8080"
    token   = "your-auth-token"
)

type CardsResponse struct {
    Cards []Card `json:"cards"`
    Count int    `json:"count"`
}

type Card struct {
    PAN       string `json:"pan"`
    MaskedPAN string `json:"masked_pan"`
    Brand     string `json:"brand"`
    CVC       string `json:"cvc"`
}

func main() {
    // Generate cards
    req, _ := http.NewRequest("GET", baseURL+"/v1/cards?brand=visa&count=5", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    
    client := &http.Client{}
    resp, _ := client.Do(req)
    defer resp.Body.Close()
    
    var cardsResp CardsResponse
    json.NewDecoder(resp.Body).Decode(&cardsResp)
    
    for _, card := range cardsResp.Cards {
        fmt.Printf("PAN: %s, CVC: %s\n", card.MaskedPAN, card.CVC)
    }
}
```

## Deployment Examples

### Docker Compose

```yaml
version: '3.8'

services:
  cardgen-api:
    image: cardgen-pro:latest
    command: serve --port 8080 --token ${API_TOKEN}
    environment:
      - CARDGEN_SECRET=${CARDGEN_SECRET}
    ports:
      - "8080:8080"
    restart: unless-stopped
    networks:
      - cardgen-net

networks:
  cardgen-net:
    driver: bridge
```

### Kubernetes

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cardgen-secrets
type: Opaque
stringData:
  secret: "your-dev-secret"
  token: "your-auth-token"

---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: cardgen-api
spec:
  replicas: 2
  selector:
    matchLabels:
      app: cardgen-api
  template:
    metadata:
      labels:
        app: cardgen-api
    spec:
      containers:
      - name: cardgen
        image: cardgen-pro:latest
        args: ["serve", "--port", "8080", "--token", "$(API_TOKEN)"]
        env:
        - name: CARDGEN_SECRET
          valueFrom:
            secretKeyRef:
              name: cardgen-secrets
              key: secret
        - name: API_TOKEN
          valueFrom:
            secretKeyRef:
              name: cardgen-secrets
              key: token
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10

---

apiVersion: v1
kind: Service
metadata:
  name: cardgen-api
spec:
  selector:
    app: cardgen-api
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

## Security Best Practices

### 1. Use Strong Tokens

```bash
# Generate secure token
openssl rand -hex 32

# Use in server
cardgen-pro serve --token $(openssl rand -hex 32)
```

### 2. Enable TLS

```bash
# Behind reverse proxy (nginx)
server {
    listen 443 ssl;
    server_name api.cardgen.local;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Authorization $http_authorization;
    }
}
```

### 3. Network Restrictions

```bash
# iptables (allow only from specific IPs)
iptables -A INPUT -p tcp --dport 8080 -s 10.0.0.0/8 -j ACCEPT
iptables -A INPUT -p tcp --dport 8080 -j DROP

# AWS Security Group
- Type: Custom TCP
- Port: 8080
- Source: 10.0.0.0/8
```

### 4. Monitoring

```bash
# Log all requests
cardgen-pro serve --port 8080 --token $TOKEN 2>&1 | tee /var/log/cardgen.log

# Monitor rate limits
grep "Rate limit exceeded" /var/log/cardgen.log | wc -l

# Monitor auth failures
grep "Unauthorized" /var/log/cardgen.log | wc -l
```

## Troubleshooting

### Server Won't Start

```bash
# Check if port is available
lsof -i :8080

# Check environment variables
env | grep CARDGEN

# Check logs
cardgen-pro serve --port 8080 --token test 2>&1
```

### 401 Unauthorized

```bash
# Verify token format
curl -v -H "Authorization: Bearer your-token" http://localhost:8080/v1/cards

# Should be "Bearer <token>", not just "<token>"
```

### 429 Rate Limited

```bash
# Wait 60 seconds for window to reset
sleep 60

# Or restart server to reset limits
```

## API Versioning

Current version: **v1**

Future versions will be available at `/v2/`, `/v3/`, etc., maintaining backward compatibility.

## Support

- **Issues:** [GitHub Issues](https://github.com/felipemacedo1/cardgen-pro/issues)
- **Discussions:** [GitHub Discussions](https://github.com/felipemacedo1/cardgen-pro/discussions)

---

**API Version:** 1.0.0  
**Last Updated:** 2025-10-21
