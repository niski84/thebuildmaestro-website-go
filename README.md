# thebuildmaestro-website-go

This is the Go version of my personal website, [thebuildmaestro.com](https://thebuildmaestro.com). I originally built it with Python and Flask, but I've been doing more work in Go lately and decided I wanted to maintain it in Go going forward.

The original Flask version is still available at [github.com/niski84/thebuildmaestro-website](https://github.com/niski84/thebuildmaestro-website) if you want to see what it looked like before.

## What It Does

It's a static site generator. I write my blog posts and project descriptions in Markdown files, and this tool converts them to HTML. The content structure stayed the same from the Flask version—I just rewrote the generator in Go.

## Features

- Reads markdown files and converts them to HTML
- Uses the same content structure as the Flask version (INI metadata files, README.md files)
- Generates a static site you can host anywhere
- Uses Tailwind CSS and HTMX for the frontend (replaced Bootstrap and jQuery)
- Syntax highlighting for code blocks
- Generates an Atom feed for RSS readers
- Photo gallery with thumbnails

## Project Structure

```
thebuildmaestro-website-go/
├── cmd/sitegen/main.go          # The CLI tool
├── internal/sitegen/            # The actual generator code
│   ├── sitegen.go               # Main generation logic
│   ├── content.go                # Reads and parses content
│   ├── templates.go             # Handles HTML templates
│   ├── markdown.go              # Converts markdown to HTML
│   └── feed.go                  # Generates the Atom feed
├── templates/                    # HTML templates
├── static/                       # Source content and assets
│   ├── content/                  # My blog posts and code projects
│   ├── images/                  # Images
│   └── css/js/                 # CSS and JavaScript
└── public/                       # Generated output (this gets deployed)
```

## Content Structure

I organize my content like this:

```
static/content/articles/<article-id>/
  ├── metadata          # INI file with title, author, dates, etc.
  └── README.md        # The actual markdown content
  └── [images/files]   # Any images or files for that article

static/content/code/<project-id>/
  ├── metadata
  └── README.md
```

The metadata file is just a simple INI file:

```ini
[metadata]
title: My Article Title
author: Nicholas Skitch
last_updated: 2024-01-15
written_on: 2024-01-15
description: What this article is about
```

## Installation

You need Go 1.22 or later. Then just:

```bash
go mod download
```

## Usage

Run the generator:

```bash
go run ./cmd/sitegen
```

This will read everything from `static/content`, convert it to HTML, and output it to the `public/` directory.

You can also customize the paths:

```bash
go run ./cmd/sitegen \
  -content static/content \
  -template templates \
  -static static \
  -out public \
  -domain https://thebuildmaestro.com
```

Or build a binary:

```bash
go build -o sitegen ./cmd/sitegen
./sitegen
```

## What Gets Generated

The generator creates a `public/` directory with:

- `index.html` - redirects to `/articles/`
- `articles/` - all my blog posts
- `code/` - code project listings
- `photos/` - photo gallery
- `contact/` - contact page
- `atom.xml` - RSS feed
- `static/` - all the CSS, images, JS files

Then I just deploy the `public/` directory to wherever I'm hosting (S3, Netlify, etc.).

## Migration from Flask

If you're migrating from the Flask version, the good news is the content structure is exactly the same. I kept all the same file formats and directory layouts so the migration was straightforward.

The original Flask project is at [github.com/niski84/thebuildmaestro-website](https://github.com/niski84/thebuildmaestro-website).

To migrate:

1. Clone the original Flask project to get the content:
   ```bash
   git clone https://github.com/niski84/thebuildmaestro-website.git
   ```

2. Copy your content over:
   ```bash
   cp -r thebuildmaestro-website/static/content thebuildmaestro-website-go/static/
   cp -r thebuildmaestro-website/static/files thebuildmaestro-website-go/static/
   cp -r thebuildmaestro-website/static/images thebuildmaestro-website-go/static/
   cp thebuildmaestro-website/static/robots.txt thebuildmaestro-website-go/static/
   cp thebuildmaestro-website/static/favicon.ico thebuildmaestro-website-go/static/
   ```

3. Generate the site:
   ```bash
   go run ./cmd/sitegen
   ```

4. Deploy the `public/` directory.

That's it. The content structure is compatible, so everything should just work.

## Tech Stack

- **Go** - The language I rewrote it in
- **Goldmark** - Markdown parser (Go equivalent of Python-Markdown)
- **Chroma** - Syntax highlighting (replaces Pygments)
- **Tailwind CSS** - For styling (via CDN)
- **HTMX** - For any interactive stuff (via CDN)
- **Gorilla Feeds** - For generating the Atom feed

## Why Go?

Honestly? Because I wanted to. I've been writing more Go code lately, and it felt like the right time to consolidate. Go compiles to a single binary, which makes deployment trivial. No virtual environments, no pip installs, no runtime dependencies. Just build it and ship it.

Plus, Go's standard library is really solid. File operations, templating, HTTP serving—it's all there. I didn't need to pull in a bunch of external dependencies. The whole generator is pretty lean.

## Running Tests

```bash
go test ./internal/sitegen -v
```

## License

BSD (same as the original Flask project)
