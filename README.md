# Trivy Security Scanning Integration Demo

A deliberately vulnerable Go application designed to demonstrate [Trivy](https://trivy.dev)'s security scanning capabilities across four key areas: **dependency vulnerabilities**, **IaC misconfigurations**, **secret detection**, and **license compliance**.

> **Warning**: This repository contains intentional security flaws for educational and demonstration purposes. **Never deploy this to production.**

## Purpose

This repository serves as a live demo environment to showcase how Trivy integrates into CI/CD pipelines via GitHub Actions. When a pull request is opened, the workflows in `.github/workflows/` automatically scan the codebase and report findings directly in the PR as a structured comment.

## What Trivy Detects in This Repo

### 1. Dependency Vulnerabilities (`go.mod`)

The application pins outdated Go modules with known CVEs:

| Package | Issue |
|---------|-------|
| `golang.org/x/crypto v0.0.0-20220314...` | Multiple CVEs in older crypto versions |
| `golang.org/x/net v0.0.0-20220225...` | HTTP/2 and URL parsing vulnerabilities |
| `golang.org/x/text v0.3.7` | Known text processing vulnerabilities |
| `github.com/dgrijalva/jwt-go v3.2.0` | Deprecated library with known JWT validation issues |
| `github.com/gin-gonic/gin v1.9.0` | Older version with known issues |

### 2. IaC Misconfigurations

**Dockerfile** (`Dockerfile`):
- Uses `golang:latest` (no pinned version)
- Runs as root (no `USER` directive)
- Uses `ADD` instead of `COPY`
- Hardcoded secrets in `ENV`
- No `HEALTHCHECK`
- Exposes unnecessary ports (22, 3306)
- No apt cache cleanup

**Terraform** (`infra/main.tf`):
- S3 bucket with public read access and no encryption
- Security group open to `0.0.0.0/0` on all ports
- RDS instance publicly accessible, unencrypted, no backups
- Hardcoded database password
- IAM policy with `Action: *` / `Resource: *`
- CloudWatch log group without encryption

### 3. Secret Leaks

**`config/credentials.json`** — contains fake but realistically formatted secrets:
- AWS Access Key ID and Secret Access Key
- GitHub Personal Access Token
- Stripe live API key
- SendGrid API key
- Slack webhook URL
- Database credentials

**`Dockerfile`** — hardcoded `ENV` values with secrets.

### 4. License Compliance

The `go.mod` includes dependencies under various licenses. Trivy's license scanner checks for copyleft licenses (GPL, AGPL, SSPL, EUPL) that could impose obligations on your codebase.

## GitHub Actions Workflows

### PR Scan (`.github/workflows/trivy-pr-scan.yml`)

Triggered on every pull request. Runs four parallel scans and posts a **sticky comment** on the PR with:

- Summary table (pass/fail per category)
- Expandable details for each finding
- Auto-fix PR creation for dependencies with known fixed versions

| Scan | What it checks |
|------|----------------|
| `vuln` | Known CVEs in Go dependencies |
| `secret` | Hardcoded secrets and API keys |
| `config` | Dockerfile and Terraform misconfigurations |
| `license` | Copyleft / restrictive licenses |

The workflow **fails the check** if any issues are found, acting as a quality gate.

### Scheduled Scan (`.github/workflows/trivy-scheduled.yml`)

Runs Monday–Friday at 07:00 UTC. Performs a full scan of the `main` branch and creates a PR with a markdown report if findings are detected. This catches newly disclosed CVEs in existing dependencies.

## Running Locally

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/your-org/trivy-integration.git
cd trivy-integration

# Start the application and database
docker compose up --build

# The API is available at http://localhost:8080
```

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Health check |
| `GET` | `/api/documents` | List all documents |
| `GET` | `/api/documents/:id` | Get a document by ID |
| `POST` | `/api/documents` | Create a document |
| `POST` | `/api/login` | Get a JWT token |

### Example Requests

```bash
# Health check
curl http://localhost:8080/health

# List documents
curl http://localhost:8080/api/documents

# Login (returns JWT)
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "password": "demo"}'
```

## Running Trivy Locally

You can run the same scans locally that the CI pipeline runs:

```bash
# Install Trivy (macOS)
brew install trivy

# Scan for vulnerabilities in dependencies
trivy fs --scanners vuln --severity CRITICAL,HIGH .

# Scan for secrets
trivy fs --scanners secret .

# Scan IaC files (Dockerfile + Terraform)
trivy config .

# Scan for license issues
trivy fs --scanners license .

# Run all scanners at once
trivy fs --scanners vuln,secret,misconfig,license .
```

## Demo Workflow

1. **Create a feature branch** from `main`
2. **Open a pull request** — the Trivy PR scan triggers automatically
3. **Review the PR comment** — it shows a structured report with all findings
4. **Observe the failed check** — the workflow fails because issues were found
5. **Check for auto-fix PRs** — for dependencies with known fixes, the workflow creates separate fix PRs
6. **Discuss remediation** — use the findings to discuss security best practices

## Project Structure

```
.
├── .github/workflows/
│   ├── trivy-pr-scan.yml        # PR security scan + auto-fix
│   └── trivy-scheduled.yml      # Scheduled daily scan
├── config/
│   └── credentials.json         # Intentionally leaked secrets
├── infra/
│   └── main.tf                  # Insecure Terraform config
├── Dockerfile                   # Insecure Dockerfile
├── docker-compose.yml           # Local deployment
├── go.mod / go.sum              # Vulnerable dependencies
├── init.sql                     # Database schema
├── main.go                      # Application code
└── README.md
```

## License

This project is for internal demonstration purposes only.
