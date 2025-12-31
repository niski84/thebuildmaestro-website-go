package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"thebuildmaestro-website-go/internal/sitegen"
)

func main() {
	log.Printf("Starting site generator")

	// Parse command line flags
	contentPath := flag.String("content", "static/content", "Path to content directory")
	templateDir := flag.String("template", "templates", "Path to templates directory")
	staticDir := flag.String("static", "static", "Path to static assets directory")
	outputDir := flag.String("out", "public", "Path to output directory")
	
	// Check for environment variable first, then use default
	defaultDomain := os.Getenv("CANONICAL_DOMAIN")
	if defaultDomain == "" {
		defaultDomain = "https://thebuildmaestro.com"
	}
	canonicalDomain := flag.String("domain", defaultDomain, "Canonical domain for URLs")

	flag.Parse()

	log.Printf("Configuration: content=%s, template=%s, static=%s, out=%s, domain=%s",
		*contentPath, *templateDir, *staticDir, *outputDir, *canonicalDomain)

	// Validate directories exist
	if err := validateDirectories(*contentPath, *templateDir, *staticDir); err != nil {
		log.Printf("Validation error: %v", err)
		os.Exit(1)
	}

	// Create config
	config := sitegen.Config{
		ContentPath:    *contentPath,
		TemplateDir:    *templateDir,
		StaticDir:      *staticDir,
		OutputDir:      *outputDir,
		CanonicalDomain: *canonicalDomain,
	}

	// Generate site
	if err := sitegen.GenerateSite(config); err != nil {
		log.Printf("Site generation failed: %v", err)
		os.Exit(1)
	}

	log.Printf("Site generation completed successfully")
}

// validateDirectories checks that required directories exist
func validateDirectories(contentPath, templateDir, staticDir string) error {
	dirs := map[string]string{
		"content":  contentPath,
		"template": templateDir,
		"static":   staticDir,
	}

	for name, path := range dirs {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("%s directory does not exist: %s", name, path)
		}
	}

	return nil
}

