package sitegen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateSite(t *testing.T) {
	// Create temporary directory structure
	tmpDir, err := os.MkdirTemp("", "sitegen-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directory structure
	contentDir := filepath.Join(tmpDir, "content")
	articlesDir := filepath.Join(contentDir, "articles", "test-article")
	codeDir := filepath.Join(contentDir, "code", "test-code")
	templateDir := filepath.Join(tmpDir, "templates")
	staticDir := filepath.Join(tmpDir, "static")
	outputDir := filepath.Join(tmpDir, "public")

	// Create directories
	if err := os.MkdirAll(articlesDir, 0755); err != nil {
		t.Fatalf("Failed to create articles directory: %v", err)
	}
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		t.Fatalf("Failed to create code directory: %v", err)
	}
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}
	if err := os.MkdirAll(staticDir, 0755); err != nil {
		t.Fatalf("Failed to create static directory: %v", err)
	}

	// Create test metadata file for article
	articleMetadata := `[metadata]
title: Test Article
author: Test Author
last_updated: 2024-01-15
written_on: 2024-01-15
description: This is a test article
`
	if err := os.WriteFile(filepath.Join(articlesDir, "metadata"), []byte(articleMetadata), 0644); err != nil {
		t.Fatalf("Failed to write article metadata: %v", err)
	}

	// Create test README.md for article
	articleReadme := "# Test Article\n\nThis is a test article content.\n\n## Section 1\n\nSome content here.\n\n```python\ndef hello():\n    print(\"Hello, World!\")\n```\n"
	if err := os.WriteFile(filepath.Join(articlesDir, "README.md"), []byte(articleReadme), 0644); err != nil {
		t.Fatalf("Failed to write article README: %v", err)
	}

	// Create test metadata file for code
	codeMetadata := `[metadata]
title: Test Code Project
author: Test Author
last_updated: 2024-01-10
written_on: 2024-01-10
description: This is a test code project
url: https://example.com/project
`
	if err := os.WriteFile(filepath.Join(codeDir, "metadata"), []byte(codeMetadata), 0644); err != nil {
		t.Fatalf("Failed to write code metadata: %v", err)
	}

	// Create base template
	baseTemplate := `<!DOCTYPE html>
<html>
<head>
	<title>{{if .Title}}{{.Title}} - {{end}}Test Site</title>
</head>
<body>
	{{block "content" .}}{{end}}
</body>
</html>
`
	if err := os.WriteFile(filepath.Join(templateDir, "base.html"), []byte(baseTemplate), 0644); err != nil {
		t.Fatalf("Failed to write base template: %v", err)
	}

	// Create articles template
	articlesTemplate := `{{define "articles.html"}}
{{block "content" .}}
<h1>Articles</h1>
{{range $id, $details := .Metadata}}
	<h2>{{$details.Title}}</h2>
	<p>{{$details.Description}}</p>
{{end}}
{{end}}
{{end}}
`
	if err := os.WriteFile(filepath.Join(templateDir, "articles.html"), []byte(articlesTemplate), 0644); err != nil {
		t.Fatalf("Failed to write articles template: %v", err)
	}

	// Create article-display template
	articleDisplayTemplate := `{{define "article-display.html"}}
{{block "content" .}}
<h1>{{.Article.Title}}</h1>
<div>{{.Content}}</div>
{{end}}
{{end}}
`
	if err := os.WriteFile(filepath.Join(templateDir, "article-display.html"), []byte(articleDisplayTemplate), 0644); err != nil {
		t.Fatalf("Failed to write article-display template: %v", err)
	}

	// Create code template
	codeTemplate := `{{define "code.html"}}
{{block "content" .}}
<h1>Code</h1>
{{range $id, $details := .Metadata}}
	<h2>{{$details.Title}}</h2>
	<p>{{$details.Description}}</p>
{{end}}
{{end}}
{{end}}
`
	if err := os.WriteFile(filepath.Join(templateDir, "code.html"), []byte(codeTemplate), 0644); err != nil {
		t.Fatalf("Failed to write code template: %v", err)
	}

	// Create photos template
	photosTemplate := `{{define "photos.html"}}
{{block "content" .}}
<h1>Photos</h1>
{{end}}
{{end}}
`
	if err := os.WriteFile(filepath.Join(templateDir, "photos.html"), []byte(photosTemplate), 0644); err != nil {
		t.Fatalf("Failed to write photos template: %v", err)
	}

	// Create contact template
	contactTemplate := `{{define "contact.html"}}
{{block "content" .}}
<h1>Contact</h1>
{{end}}
{{end}}
`
	if err := os.WriteFile(filepath.Join(templateDir, "contact.html"), []byte(contactTemplate), 0644); err != nil {
		t.Fatalf("Failed to write contact template: %v", err)
	}

	// Create static files directory
	staticFilesDir := filepath.Join(staticDir, "files", "photos")
	if err := os.MkdirAll(staticFilesDir, 0755); err != nil {
		t.Fatalf("Failed to create photos directory: %v", err)
	}

	// Create robots.txt
	if err := os.WriteFile(filepath.Join(staticDir, "robots.txt"), []byte("User-agent: *\nDisallow:"), 0644); err != nil {
		t.Fatalf("Failed to write robots.txt: %v", err)
	}

	// Run generator
	config := Config{
		ContentPath:    contentDir,
		TemplateDir:    templateDir,
		StaticDir:      staticDir,
		OutputDir:      outputDir,
		CanonicalDomain: "https://example.com",
	}

	if err := GenerateSite(config); err != nil {
		t.Fatalf("GenerateSite failed: %v", err)
	}

	// Verify output files exist
	expectedFiles := []string{
		"index.html",
		"articles/index.html",
		"articles/test-article/index.html",
		"code/index.html",
		"photos/index.html",
		"contact/index.html",
		"atom.xml",
		"robots.txt",
	}

	for _, file := range expectedFiles {
		filePath := filepath.Join(outputDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file does not exist: %s", filePath)
		}
	}
}

