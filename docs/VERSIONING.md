# Semantic Versioning Guidelines

This document outlines the semantic versioning guidelines for the Terraform Spotify Provider. Following these guidelines ensures that version numbers communicate meaningful information about the underlying code and its compatibility.

## Semantic Versioning Format

All version numbers follow the format of `MAJOR.MINOR.PATCH` (e.g., `1.2.3`), optionally followed by pre-release identifiers (e.g., `1.2.3-beta.1`).

### Version Components

1. **MAJOR** version increments indicate incompatible API changes
   - Increment when you make breaking changes to the provider
   - Users will need to modify their Terraform configurations when upgrading

2. **MINOR** version increments indicate added functionality in a backward-compatible manner
   - Increment when you add new resources, data sources, or features
   - Existing Terraform configurations should continue to work without modification

3. **PATCH** version increments indicate backward-compatible bug fixes
   - Increment for bug fixes, performance improvements, or minor updates
   - No new functionality or breaking changes

4. **Pre-release** versions (e.g., `-alpha.1`, `-beta.2`, `-rc.1`)
   - Use for versions that are not yet stable or ready for production
   - Clearly communicate the stability level to users

## Version Management Process

### Updating the Version

1. Update the version in `version/version.go`:
   ```go
   // Version is the main version number that is being run at the moment.
   Version = "X.Y.Z"
   
   // VersionPrerelease is a pre-release marker for the version.
   // Set to empty string for final releases.
   VersionPrerelease = "" // or "dev", "alpha", "beta", "rc1", etc.
   ```

2. Update the `CHANGELOG.md` file with detailed information about the changes in the new version

3. Create a new release using the Makefile target:
   ```sh
   make release
   ```

### Release Workflow

When a new version tag is pushed (e.g., `v1.2.3`), the following automated processes occur:

1. GitHub Actions workflow builds the provider for multiple platforms
2. GoReleaser creates release artifacts and generates checksums
3. Docker images are built and tagged with the version

## Version Compatibility

### Terraform Core Compatibility

Ensure that the provider remains compatible with the specified Terraform versions in the documentation. Major version bumps may be required when dropping support for older Terraform versions.

### Spotify API Compatibility

Track changes to the Spotify API and update the provider accordingly:
- Breaking API changes may require a major version bump
- New API features can be added in minor version releases
- API bug fixes or small adjustments can be included in patch releases

## Communication

Clearly communicate version changes to users:

1. Maintain a detailed `CHANGELOG.md` file
2. Document upgrade paths for major version changes
3. Provide deprecation notices before removing features

## Example Version Increments

- **MAJOR** (`1.0.0` → `2.0.0`): Changing resource attribute types, removing resources or attributes
- **MINOR** (`1.0.0` → `1.1.0`): Adding new resources, data sources, or optional attributes
- **PATCH** (`1.0.0` → `1.0.1`): Fixing bugs, improving error messages, performance optimizations

## References

- [Semantic Versioning 2.0.0](https://semver.org/)
- [HashiCorp Versioning Specification](https://www.terraform.io/plugin/sdkv2/best-practices/versioning)