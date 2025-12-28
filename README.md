# thebuildmaestro-website-go

A Go-based static site generator for [thebuildmaestro.com](https://thebuildmaestro.com). This project is a modern rewrite of the original Flask/Python website, maintaining full compatibility with the existing content structure while modernizing the frontend with HTMX and Tailwind CSS.

**Original Flask Project**: [github.com/niski84/thebuildmaestro-website](https://github.com/niski84/thebuildmaestro-website)

## Features

- **File-based Content Management**: Reads INI-style metadata files and Markdown content from directories
- **Static Site Generation**: Pre-renders all pages at build time for fast, static hosting
- **Modern Frontend**: Uses Tailwind CSS and HTMX for a modern, responsive UI
- **Full Compatibility**: Maintains compatibility with the existing Flask project's content structure
- **Atom Feed**: Generates RSS/Atom feed for recent articles
- **Photo Gallery**: Displays photos with thumbnails in a responsive grid
- **Code Highlighting**: Syntax highlighting for code blocks using Chroma

## Project Structure

```
thebuildmaestro-website-go/
├── cmd/
│   └── sitegen/
│       └── main.go              # CLI entry point
├── internal/
│   └── sitegen/
│       ├── sitegen.go           # Core generator logic
│       ├── content.go           # Content loading and metadata parsing
│       ├── templates.go         # Template rendering
│       ├── markdown.go          # Markdown processing
│       ├── feed.go              # Atom feed generation
│       └── sitegen_test.go      # Integration tests
├── templates/
│   ├── base.html                # Base template (HTMX + Tailwind)
│   ├── articles.html            # Articles listing
│   ├── article-display.html     # Individual article
│   ├── code.html                # Code projects listing
│   ├── photos.html              # Photo gallery
│   └── contact.html             # Contact page
├── static/
│   ├── content/                 # Source content (from Flask project)
│   │   ├── articles/
│   │   └── code/
│   ├── files/
│   │   └── photos/
│   ├── css/
│   ├── js/
│   └── images/
├── public/                      # Generated output (build directory)
└── go.mod                       # Go module definition
```

## Content Structure

The generator expects content in the following structure (compatible with the Flask project):

```
static/content/articles/<id>/
  ├── metadata          # INI format with title, author, dates, etc.
  └── README.md        # Markdown content
  └── [other assets]   # Images, files referenced in markdown

static/content/code/<id>/
  ├── metadata
  └── README.md
```

### Metadata File Format

The `metadata` file uses INI format:

```ini
[metadata]
title: Article Title
author: Author Name
last_updated: 2024-01-15
written_on: 2024-01-15
description: Article description
distributions: RHEL_7 CentOS_7 Ubuntu_20.04
url: https://example.com (optional, for external links)
```

## Installation

1. Ensure you have Go 1.22 or later installed
2. Clone or navigate to the project directory
3. Install dependencies:

```bash
go mod download
```

## Usage

### Basic Usage

Generate the static site with default paths:

```bash
go run ./cmd/sitegen
```

### Custom Paths

Specify custom paths for content, templates, static files, and output:

```bash
go run ./cmd/sitegen \
  -content static/content \
  -template templates \
  -static static \
  -out public \
  -domain https://thebuildmaestro.com
```

### Command Line Flags

- `-content`: Path to content directory (default: `static/content`)
- `-template`: Path to templates directory (default: `templates`)
- `-static`: Path to static assets directory (default: `static`)
- `-out`: Path to output directory (default: `public`)
- `-domain`: Canonical domain for URLs (default: `https://thebuildmaestro.com`)

### Build Binary

Build a standalone binary:

```bash
go build -o sitegen ./cmd/sitegen
./sitegen -content static/content -out public
```

## Generated Output

The generator creates the following structure in the output directory:

```
public/
├── index.html              # Redirects to /articles/
├── articles/
│   ├── index.html          # Articles listing
│   └── <id>/
│       ├── index.html      # Article page
│       └── [assets]        # Copied from source
├── code/
│   └── index.html          # Code listing
├── photos/
│   └── index.html          # Photo gallery
├── contact/
│   └── index.html          # Contact page
├── atom.xml                # RSS feed
├── robots.txt
├── favicon.ico
└── static/                 # Copied static assets
```

## Migration from Flask Project

This generator is designed to be compatible with the existing Flask project structure. The original Flask/Python project can be found at [github.com/niski84/thebuildmaestro-website](https://github.com/niski84/thebuildmaestro-website).

### Compatibility

1. **Content Compatibility**: The generator reads the same metadata files and Markdown structure
2. **URL Structure**: Maintains the same URL paths (`/articles/`, `/code/`, etc.)
3. **Asset Handling**: Copies article assets and static files to the same locations

### Migration Steps

1. Clone or download the original Flask project:
   ```bash
   git clone https://github.com/niski84/thebuildmaestro-website.git
   ```

2. Copy your existing content from the Flask project:
   ```bash
   cp -r thebuildmaestro-website/static/content thebuildmaestro-website-go/static/
   cp -r thebuildmaestro-website/static/files thebuildmaestro-website-go/static/
   cp thebuildmaestro-website/static/robots.txt thebuildmaestro-website-go/static/
   cp thebuildmaestro-website/static/favicon.ico thebuildmaestro-website-go/static/
   ```

3. Copy images and other static assets:
   ```bash
   cp -r thebuildmaestro-website/static/images thebuildmaestro-website-go/static/
   ```

4. Generate the site:
   ```bash
   go run ./cmd/sitegen
   ```

5. Deploy the `public/` directory to your hosting service (S3, Netlify, etc.)

## Technology Stack

- **Go 1.22+**: Core language
- **Goldmark**: Markdown processing
- **Chroma**: Syntax highlighting
- **Tailwind CSS**: Styling (via CDN)
- **HTMX**: Interactive features (via CDN)
- **Gorilla Feeds**: Atom feed generation
- **INI Parser**: Metadata file parsing

## Development

### Running Tests

```bash
go test ./internal/sitegen -v
```

### Code Structure

- `internal/sitegen/content.go`: Handles loading and parsing content metadata
- `internal/sitegen/markdown.go`: Converts Markdown to HTML with syntax highlighting
- `internal/sitegen/templates.go`: Manages template loading and rendering
- `internal/sitegen/sitegen.go`: Core generation logic
- `internal/sitegen/feed.go`: Atom feed generation
- `cmd/sitegen/main.go`: CLI entry point

## Logging

All exported functions log their input parameters before executing business logic. The main function logs:
- Start of generation
- Progress through each step
- Completion or errors
- Exits with non-zero code on errors

## License

BSD (same as the original Flask project)

