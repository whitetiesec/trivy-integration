# ============================================================
# INTENTIONALLY INSECURE DOCKERFILE — for Trivy demo purposes
# ============================================================
# This Dockerfile contains several IaC misconfigurations that
# Trivy will detect. DO NOT use this in production.

# Flaw 1: Using 'latest' tag — no reproducible builds
FROM golang:latest

# Flaw 2: Running as root (no USER directive)
# Flaw 3: Using ADD instead of COPY (unnecessary remote fetch capability)
ADD . /app
WORKDIR /app

# Flaw 4: Hardcoded secret in ENV
ENV APP_SECRET_KEY=supersecretkey123
ENV DATABASE_URL=postgres://admin:password123@db:5432/trivy_demo

# Flaw 5: Installing packages without version pinning and not cleaning cache
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    netcat-openbsd

# Flaw 6: Using go build without CGO_ENABLED=0 (larger attack surface)
RUN go build -o /trivy-demo .

# Flaw 7: Exposing multiple ports unnecessarily
EXPOSE 8080 22 3306

# Flaw 8: No HEALTHCHECK defined

# Flaw 9: Running the binary as root (no USER instruction)
CMD ["/trivy-demo"]
