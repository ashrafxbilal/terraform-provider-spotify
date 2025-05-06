# Terraform Spotify Provider Documentation Site

This directory contains the source files for the Terraform Spotify Provider documentation site, which is hosted on GitHub Pages.

## Overview

The documentation site provides comprehensive information about the Terraform Spotify Provider, including:

- Getting started guides
- Resource and data source documentation
- Examples and use cases
- API reference

## Local Development

To run the documentation site locally:

1. Install Jekyll and its dependencies:

```bash
gem install jekyll bundler
bundle install
```

2. Start the local server:

```bash
bundle exec jekyll serve
```

3. Open your browser and navigate to `http://localhost:4000/terraform-provider-spotify/`

## File Structure

- `index.html`: Main landing page
- `styles.css`: CSS styles for the site
- `script.js`: JavaScript for interactive elements
- `_config.yml`: Jekyll configuration
- `resources/`: Documentation for provider resources
- `data-sources/`: Documentation for provider data sources
- `examples/`: Example configurations and use cases

## Contributing

To contribute to the documentation:

1. Fork the repository
2. Create a new branch for your changes
3. Make your changes to the documentation
4. Submit a pull request

## Deployment

The documentation site is automatically deployed to GitHub Pages when changes are pushed to the main branch.