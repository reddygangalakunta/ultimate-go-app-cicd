# Enterprise Go Application & CI/CD Jenkins Pipeline

Production-grade Go microservice with an enterprise multi-stage **Jenkins Declarative Pipeline** (`Jenkinsfile`) utilizing **Docker Container Agents** for all stages, multi-stage hardened `Dockerfile`, and automated GitOps manifest image tag updating via shell script.

---

## 🚀 Key Features

### 1. Enterprise Go Microservice
- **Clean Standard Layout**: Modular structure separating `cmd/`, `internal/config`, `internal/logger`, `internal/model`, `internal/service`, `internal/handler`, and `internal/middleware`.
- **Structured JSON Logging**: Powered by Go standard `log/slog`.
- **Graceful Shutdown**: Listens for SIGINT/SIGTERM OS signals and drains HTTP connections with a configurable timeout.
- **Production Health Checks**: Standard Kubernetes endpoints (`/healthz`, `/livez`, `/readyz`).
- **Domain API**: Orders REST API with request logging, panic recovery middleware, and high unit test coverage.

### 2. Enterprise Hardened Dockerfile
- **Multi-Stage Build**: `golang:1.22-alpine` builder phase + `alpine:3.19` minimal runtime phase.
- **Security-First**: Executed under custom non-root system user `appuser` (UID `10001`).
- **Binary Optimization**: Static compilation (`CGO_ENABLED=0`) with symbol stripping (`-ldflags="-s -w"`) and dynamic version injection.
- **Healthcheck & Annotations**: Built-in container `HEALTHCHECK` directive and standard OCI labels.

### 3. Multi-Stage Containerized Jenkins Pipeline (`Jenkinsfile`)
- **Docker Agent per Stage**: **No local tool pre-installation needed on Jenkins runner host**. Every stage runs in dedicated official Docker containers (`golangci-lint`, `golang`, `docker:26-cli`, `aquasec/trivy`, `alpine/git`).
- **Enterprise Pipeline Stages**:
  1. **Checkout & Environment Init**: Reads semantic version from `VERSION`, generates version tag `v<VERSION>.<BUILD_NUMBER>`.
  2. **Static Analysis & Linting**: Runs `golangci-lint` to enforce code standard policies.
  3. **SAST Security Scan**: Runs `gosec` (Go security scanner) and `govulncheck` (vulnerability database audit).
  4. **Unit Tests & Code Coverage**: Executes `go test -v -race -coverprofile=coverage.out` and checks coverage threshold (`70%`).
  5. **Application Compilation**: Compiles static binary with embedded build tags.
  6. **Docker Image Build**: Builds multi-stage Docker image with version and latest tags.
  7. **Container Vulnerability Scan (Trivy)**: Scans built image for HIGH & CRITICAL CVE vulnerabilities.
  8. **Push to Docker Registry**: Securely logs in using Jenkins credentials store and pushes tagged image.
  9. **Update Git Manifest Image Tag**: Executes `scripts/update-image-tag.sh` to update Kubernetes deployment manifest (`deployments/k8s/deployment.yaml`), commits changes, tags release, and pushes to remote Git repository.

### 4. Git Image Tag Update Script (`scripts/update-image-tag.sh`)
- Robust bash script with `set -euo pipefail`.
- Modifies Kubernetes manifest (`deployments/k8s/deployment.yaml`) using `sed`.
- Automates Git identity configuration, stages changed files, creates commit `[ci skip] chore(deploy): update image tag to <TAG>`, creates release git tag, and pushes back to origin repository.

---

## 📁 Repository Structure

```
.
├── cmd/
│   └── server/
│       └── main.go                 # Main entrypoint with graceful shutdown & signal handling
├── internal/
│   ├── config/
│   │   └── config.go               # Structured configuration loading with env vars & defaults
│   ├── handler/
│   │   ├── health.go               # Healthcheck handlers (/healthz, /readyz, /livez)
│   │   ├── health_test.go          # Handler unit tests
│   │   ├── order.go                # Enterprise domain REST endpoint (Orders API)
│   │   └── order_test.go           # Domain handler tests
│   ├── logger/
│   │   └── logger.go               # Structured JSON logger using log/slog
│   ├── middleware/
│   │   ├── logging.go              # HTTP request logging & duration tracking middleware
│   │   └── recovery.go             # HTTP panic recovery middleware
│   └── service/
│       ├── order_service.go        # Business logic service layer
│       └── order_service_test.go   # Business logic tests
├── deployments/
│   └── k8s/
│       └── deployment.yaml         # Kubernetes deployment manifest (target for image tag update)
├── scripts/
│   └── update-image-tag.sh         # Shell script to update image tag in git manifest & commit
├── .golangci.yml                   # Enterprise linter rules configuration
├── Dockerfile                      # Hardened multi-stage Dockerfile (non-root, minimal image)
├── Jenkinsfile                     # Declarative multi-stage Jenkins Pipeline using Docker agents
├── Makefile                        # Dev automation targets (lint, test, build, docker)
├── VERSION                         # Semantic version file
├── go.mod                          # Go module definition
└── README.md                       # Documentation
```

---

## 🛠 Local Development & Testing

### Using Makefile

```bash
# Run unit tests with race detection
make test

# Generate coverage report (outputs coverage.html)
make coverage

# Run linter
make lint

# Compile static binary
make build

# Run service locally
make run

# Build Docker image
make docker-build

# Test Git image tag updater script (Dry-Run)
make update-tag
```

---

## 🔑 Jenkins Pipeline Credentials Setup

Ensure the following credential IDs are configured in Jenkins:

1. **`docker-registry-credentials`** *(Username with Password)*:
   - Credentials for pushing images to Docker Registry.
2. **`git-credentials`** *(SSH Username with Private Key)*:
   - Credentials for committing and pushing updated Kubernetes deployment manifests back to Git repository.
