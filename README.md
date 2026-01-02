# whoami

whoami is a container-native microservice designed to provide high-concurrency system metadata and network introspection. Purpose-built for orchestration environments, it is delivered as a zero-dependency distroless image.

## 📦 Container Product Specifications

whoami is distributed exclusively as a container image. It is built on the scratch base image, ensuring it contains no OS vulnerabilities, no shell, and no unnecessary binaries.

Base Image: scratch \
Binary Type: Statically linked Go (CGO disabled) \
Signal Handling: Full support for SIGTERM / SIGINT for graceful orchestration shutdowns. \
Security: Runs as a non-privileged user (UID 1000) by default.

## 🚀 Quick Start

Pull the latest version from the GitHub Container Registry:
```bash
docker pull ghcr.io/structx/whoami:latest
```

Run the container:
```bash 
docker run -d \
  --name whoami-prod \
  -p 8080:8080 \
  ghcr.io/structx/whoami:latest
```

## 🏗 Build & Architecture

The product is compiled using a multi-stage process to eliminate build-tool leakage and minimize the final image size.
Build Flags

The binary is optimized for production using the following parameters:

`-ldflags="-s -w"` Strips symbol and DWARF tables.\
`-trimpath` Removes file system paths from the compiled binary for better privacy/security.\
`CGO_ENABLED=0` Ensures a pure-Go static binary.

## 📖 API Documentation

The following endpoints are exposed on the container's configured port (default 8080):\
Endpoint	Format	Description\
`/`	        JSON	Returns a concise string of hostname (OS/Arch).\
`/health`	JSON	Liveness/Readiness probe (Returns 200 OK).


# ⚙️ Environment Configuration
Variable	        Default	Purpose\
PORT	8080	    The internal container port for the API server.\
HOST    127.0.0.1   internal container host for the API server.

