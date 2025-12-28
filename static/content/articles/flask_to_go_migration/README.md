I've been running thebuildmaestro.com for a while now, and the original version was built with Python and Flask. It worked fine. The concept was simple: read markdown files from directories, convert them to HTML, and serve them up. I used Frozen-Flask to pre-render everything into static HTML files, which meant I could host it anywhere without needing a Python runtime.

But here's the thing. I've been doing more and more work in Go lately, and I started thinking: why am I still maintaining a Python project when I could just rewrite this in Go? The core idea doesn't change. It's still markdown to HTML. It's still file-based content. It's still a static site generator.

So I did it.

## What Changed (And What Didn't)

The philosophy stayed the same. Content lives in markdown files. Each article or project gets its own directory with a `README.md` and a `metadata` file. The generator walks those directories, parses the markdown, converts it to HTML, and spits out a static site.

What changed was the tooling. Instead of Python and Frozen-Flask, I'm now using:

- **Go's standard library** for file operations and templating
- **Goldmark** for markdown parsing (the Go equivalent of Python-Markdown)
- **Chroma** for syntax highlighting (replacing Pygments)
- **Go's html/template** package for rendering
- **Tailwind CSS** and **HTMX** for the frontend (replacing Bootstrap and jQuery)

The result is a single binary that does everything. No virtual environments. No pip installs. No runtime dependencies. Just compile it and run it.

## The Migration

I kept all the existing content structure. The markdown files, the metadata format, the directory layout—all of it stayed the same. That made the migration straightforward. I just had to port the logic from Python to Go.

The templates got a refresh too. I converted from Jinja2 to Go templates, swapped Bootstrap for Tailwind, and replaced jQuery with HTMX. The site looks more modern now, but it still feels like the same site. Same content, same structure, just faster and simpler to deploy.

## Why Go?

Honestly? Because I wanted to. I've been writing more Go code lately, and it felt like the right time to consolidate. Go compiles to a single binary, which makes deployment trivial. No need to worry about Python versions or dependencies. Just build it and ship it.

Plus, Go's standard library is really solid. File operations, templating, HTTP serving—it's all there. I didn't need to pull in a bunch of external dependencies. The whole generator is pretty lean.

## The Repos

If you want to see the code, both versions are on GitHub:

- **Original Flask version**: [github.com/niski84/thebuildmaestro-website](https://github.com/niski84/thebuildmaestro-website)
- **New Go version**: [github.com/niski84/thebuildmaestro-website-go](https://github.com/niski84/thebuildmaestro-website-go)

The old one still works. The new one does the same thing, just with different tools. Sometimes a rewrite is just about using the tools you want to use.

