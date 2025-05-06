# Runtime Dockerfile for Terraform Spotify Provider
# Multi-stage build with distroless final image

# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    make \
    bash \
    curl \
    gcc \
    musl-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Update go.mod
RUN go mod tidy

# Build the provider with version information
RUN GIT_COMMIT=$(git rev-parse --short HEAD) \
    && BUILD_DATE=$(date -u '+%Y-%m-%d') \
    && go build -ldflags "-X github.com/ashrafxbilal/terraform-provider-spotify/version.GitCommit=${GIT_COMMIT} -X github.com/ashrafxbilal/terraform-provider-spotify/version.BuildDate=${BUILD_DATE}" \
    -o terraform-provider-spotify

# Terraform stage
FROM hashicorp/terraform:1.5.7 AS terraform

# Final runtime stage
FROM gcr.io/distroless/base-debian11

# Copy the Terraform binary from the terraform stage
COPY --from=terraform /bin/terraform /usr/local/bin/terraform

# Copy the provider binary from the builder stage
COPY --from=builder /app/terraform-provider-spotify /usr/local/bin/terraform-provider-spotify

# Create necessary directories for Terraform
COPY --from=builder /app/docs/examples /examples

# Set working directory
WORKDIR /workspace

# Set environment variables
ENV TF_PLUGIN_CACHE_DIR="/plugins"

# Create a volume for Terraform state
VOLUME ["/workspace", "/plugins"]

# Default command
CMD ["/usr/local/bin/terraform"]