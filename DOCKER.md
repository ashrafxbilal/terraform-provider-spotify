# Docker Usage Guide

This document provides instructions for using the Docker containers for the Terraform Spotify Provider.

## Available Docker Images

We provide two Docker images:

1. **Runtime Image**: `ashrafxbilal/terraform-spotify-provider`
   - Contains Terraform CLI and the pre-built provider
   - Optimized for running Terraform configurations with the Spotify provider
   - Minimal size and dependencies for production use

2. **Development Image**: `ashrafxbilal/terraform-spotify-provider-dev`
   - Contains the full Go development environment
   - Includes all build tools and dependencies
   - Designed for contributing to the provider codebase

Both images use multi-stage builds with distroless base images for improved security and reduced size.

## Using the Runtime Image

### Quick Start

```bash
# Pull the image
docker pull ashrafxbilal/terraform-spotify-provider:latest

# Run Terraform with the Spotify provider
docker run -it --rm \
  -v $(pwd):/workspace \
  -e SPOTIFY_CLIENT_ID=your_client_id \
  -e SPOTIFY_CLIENT_SECRET=your_client_secret \
  -e SPOTIFY_REDIRECT_URI=your_redirect_uri \
  -e SPOTIFY_REFRESH_TOKEN=your_refresh_token \
  ashrafxbilal/terraform-spotify-provider:latest init

# Apply your configuration
docker run -it --rm \
  -v $(pwd):/workspace \
  -e SPOTIFY_CLIENT_ID=your_client_id \
  -e SPOTIFY_CLIENT_SECRET=your_client_secret \
  -e SPOTIFY_REDIRECT_URI=your_redirect_uri \
  -e SPOTIFY_REFRESH_TOKEN=your_refresh_token \
  ashrafxbilal/terraform-spotify-provider:latest apply
```

### Using with Example Configurations

```bash
# Run the basic playlist example
docker run -it --rm \
  -v $(pwd):/workspace \
  -e SPOTIFY_CLIENT_ID=your_client_id \
  -e SPOTIFY_CLIENT_SECRET=your_client_secret \
  -e SPOTIFY_REDIRECT_URI=your_redirect_uri \
  -e SPOTIFY_REFRESH_TOKEN=your_refresh_token \
  ashrafxbilal/terraform-spotify-provider:latest -chdir=/examples/basic_playlist apply
```

## Using the Development Image

### Development Workflow

```bash
# Pull the development image
docker pull ashrafxbilal/terraform-spotify-provider-dev:latest

# Start a development container
docker run -it --rm \
  -v $(pwd):/app \
  -w /app \
  ashrafxbilal/terraform-spotify-provider-dev:latest

# Inside the container, you can build and test the provider
go build
make test
```

### Running the Auth Proxy

```bash
# Start a container with port forwarding for the auth proxy
docker run -it --rm \
  -v $(pwd):/app \
  -w /app \
  -p 8080:8080 \
  -e SPOTIFY_CLIENT_ID=your_client_id \
  -e SPOTIFY_CLIENT_SECRET=your_client_secret \
  -e SPOTIFY_REDIRECT_URI=http://localhost:8080/callback \
  ashrafxbilal/terraform-spotify-provider-dev:latest

# Inside the container, run the auth proxy
cd spotify_auth_proxy
go run spotify-auth.go
```

## Building the Images Locally

### Runtime Image

```bash
docker build -t terraform-spotify-provider:local .
```

### Development Image

```bash
docker build -t terraform-spotify-provider-dev:local -f Dockerfile.dev .
```

## Environment Variables

Both images support the following environment variables:

- `SPOTIFY_CLIENT_ID`: Your Spotify API client ID
- `SPOTIFY_CLIENT_SECRET`: Your Spotify API client secret
- `SPOTIFY_REDIRECT_URI`: Your Spotify API redirect URI
- `SPOTIFY_REFRESH_TOKEN`: Your Spotify API refresh token
- `WEATHER_API_KEY`: Your OpenWeatherMap API key (for weather-based playlists)

## Volumes

- `/workspace`: Mount your Terraform configurations here in the runtime image
- `/plugins`: Terraform plugin cache directory in the runtime image
- `/app`: Mount your source code here in the development image

## CI/CD Integration

These Docker images are automatically built and published to Docker Hub when:

1. A new version tag is pushed (e.g., `v0.1.0`)
2. Changes are made to the Dockerfiles
3. The main branch is updated

The images are tagged with:

- Semantic version (e.g., `v0.1.0`)
- Major.minor version (e.g., `0.1`)
- Branch name (for main branch)
- Commit SHA (for all builds)

## Security Considerations

- Both images use distroless base images to minimize attack surface
- No shell or unnecessary utilities in the runtime image
- Secrets should be passed as environment variables, not built into the images