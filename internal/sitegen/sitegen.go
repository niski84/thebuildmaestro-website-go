package sitegen

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Config holds the configuration for site generation
type Config struct {
	ContentPath  string
	TemplateDir  string
	StaticDir    string
	OutputDir    string
	CanonicalDomain string
}

// GenerateSite generates the complete static site
func GenerateSite(config Config) error {
	log.Printf("GenerateSite: config=%+v", config)

	// Clean and create output directory
	if err := prepareOutputDir(config.OutputDir); err != nil {
		return fmt.Errorf("failed to prepare output directory: %w", err)
	}

	// Load content registry
	registry, err := LoadContentRegistry(config.ContentPath)
	if err != nil {
		return fmt.Errorf("failed to load content registry: %w", err)
	}

	// Initialize template manager
	tm, err := NewTemplateManager(config.TemplateDir)
	if err != nil {
		return fmt.Errorf("failed to initialize template manager: %w", err)
	}

	// Generate homepage
	if err := generateHomepage(config, tm); err != nil {
		return fmt.Errorf("failed to generate homepage: %w", err)
	}

	// Generate about page
	if err := generateAboutPage(config, tm); err != nil {
		return fmt.Errorf("failed to generate about page: %w", err)
	}

	// Generate lean labs page
	if err := generateLeanLabsPage(config, tm); err != nil {
		return fmt.Errorf("failed to generate lean labs page: %w", err)
	}

	// Generate service detail pages
	if err := generateServiceSaasPage(config, tm); err != nil {
		return fmt.Errorf("failed to generate SaaS service page: %w", err)
	}

	// Generate Moody's Analytics page
	if err := generateMoodysPage(config, tm); err != nil {
		return fmt.Errorf("failed to generate Moody's Analytics page: %w", err)
	}

	// Generate Work with me page
	if err := generateWorkWithMePage(config, tm); err != nil {
		return fmt.Errorf("failed to generate Work with me page: %w", err)
	}

	// Filter articles to exclude future-dated posts
	publishedArticles := FilterPublishedArticles(registry["articles"])

	// Generate articles listing
	if err := generateArticlesListing(config, tm, publishedArticles); err != nil {
		return fmt.Errorf("failed to generate articles listing: %w", err)
	}

	// Generate individual articles
	if err := generateArticles(config, tm, publishedArticles); err != nil {
		return fmt.Errorf("failed to generate articles: %w", err)
	}

	// Generate code listing
	if err := generateCodeListing(config, tm, registry["code"]); err != nil {
		return fmt.Errorf("failed to generate code listing: %w", err)
	}

	// Generate photos page
	if err := generatePhotosPage(config, tm); err != nil {
		return fmt.Errorf("failed to generate photos page: %w", err)
	}

	// Generate contact page
	if err := generateContactPage(config, tm); err != nil {
		return fmt.Errorf("failed to generate contact page: %w", err)
	}

	// Copy static assets
	if err := copyStaticAssets(config.StaticDir, config.OutputDir); err != nil {
		return fmt.Errorf("failed to copy static assets: %w", err)
	}

	// Copy article assets (use published articles to avoid copying assets for future posts)
	if err := copyArticleAssets(config.ContentPath, config.OutputDir, publishedArticles); err != nil {
		return fmt.Errorf("failed to copy article assets: %w", err)
	}

	// Copy robots.txt and favicon.ico
	if err := copyStaticFiles(config.StaticDir, config.OutputDir); err != nil {
		return fmt.Errorf("failed to copy static files: %w", err)
	}

	// Generate Atom feed (only include published articles)
	if err := GenerateAtomFeed(config, publishedArticles, config.OutputDir); err != nil {
		return fmt.Errorf("failed to generate Atom feed: %w", err)
	}

	// Generate code highlighting CSS
	if err := generateCodeHighlightCSS(config.OutputDir); err != nil {
		return fmt.Errorf("failed to generate code highlighting CSS: %w", err)
	}

	log.Printf("Site generation completed successfully")
	return nil
}

// generateCodeHighlightCSS generates and writes the code highlighting CSS file
func generateCodeHighlightCSS(outputDir string) error {
	log.Printf("generateCodeHighlightCSS: outputDir=%s", outputDir)

	css, err := GetCodeHighlightCSS()
	if err != nil {
		return fmt.Errorf("failed to get code highlight CSS: %w", err)
	}

	cssDir := filepath.Join(outputDir, "static", "css")
	if err := os.MkdirAll(cssDir, 0755); err != nil {
		return fmt.Errorf("failed to create css directory: %w", err)
	}

	cssPath := filepath.Join(cssDir, "codehilite.css")
	if err := os.WriteFile(cssPath, []byte(css), 0644); err != nil {
		return fmt.Errorf("failed to write CSS file: %w", err)
	}

	return nil
}

// prepareOutputDir removes and recreates the output directory
func prepareOutputDir(outputDir string) error {
	log.Printf("prepareOutputDir: outputDir=%s", outputDir)

	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("failed to remove output directory: %w", err)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return nil
}

// generateHomepage generates the new homepage
func generateHomepage(config Config, tm *TemplateManager) error {
	log.Printf("generateHomepage")

	outputPath := filepath.Join(config.OutputDir, "index.html")
	
	data := TemplateData{
		Title:         "Home",
		CanonicalDomain: config.CanonicalDomain,
		Year:          time.Now().Year(),
	}

	html, err := tm.RenderTemplate("home.html", data)
	if err != nil {
		return fmt.Errorf("failed to render home template: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write homepage: %w", err)
	}

	return nil
}

// generateAboutPage generates the about page
func generateAboutPage(config Config, tm *TemplateManager) error {
	log.Printf("generateAboutPage")

	outputPath := filepath.Join(config.OutputDir, "about", "index.html")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create about directory: %w", err)
	}

	data := TemplateData{
		Title:         "About Me",
		CanonicalDomain: config.CanonicalDomain,
		Year:          time.Now().Year(),
	}

	html, err := tm.RenderTemplate("about.html", data)
	if err != nil {
		return fmt.Errorf("failed to render about template: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write about page: %w", err)
	}

	return nil
}

// generateLeanLabsPage generates the lean labs page
func generateLeanLabsPage(config Config, tm *TemplateManager) error {
	log.Printf("generateLeanLabsPage")

	outputPath := filepath.Join(config.OutputDir, "leanlabs", "index.html")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create leanlabs directory: %w", err)
	}

	data := TemplateData{
		Title:         "Lean Labs | Nicholas Skitch",
		CanonicalDomain: config.CanonicalDomain,
		Year:          time.Now().Year(),
	}

	html, err := tm.RenderTemplate("leanlabs.html", data)
	if err != nil {
		return fmt.Errorf("failed to render leanlabs template: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write leanlabs page: %w", err)
	}

	return nil
}

// generateServiceSaasPage generates the SaaS service detail page
func generateServiceSaasPage(config Config, tm *TemplateManager) error {
	log.Printf("generateServiceSaasPage")

	outputPath := filepath.Join(config.OutputDir, "services", "saas", "index.html")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create services/saas directory: %w", err)
	}

	data := TemplateData{
		Title:         "SaaS & Platform Engineering | Nicholas Skitch",
		CanonicalDomain: config.CanonicalDomain,
		Year:          time.Now().Year(),
	}

	html, err := tm.RenderTemplate("service-saas.html", data)
	if err != nil {
		return fmt.Errorf("failed to render service-saas template: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write service-saas page: %w", err)
	}

	return nil
}

// generateMoodysPage generates the Moody's Analytics page
func generateMoodysPage(config Config, tm *TemplateManager) error {
	log.Printf("generateMoodysPage")

	outputPath := filepath.Join(config.OutputDir, "moodys", "index.html")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create moodys directory: %w", err)
	}

	data := TemplateData{
		Title:           "Moody's Analytics | Nicholas Skitch",
		CanonicalDomain: config.CanonicalDomain,
		Year:            time.Now().Year(),
	}

	html, err := tm.RenderTemplate("moodys.html", data)
	if err != nil {
		return fmt.Errorf("failed to render moodys template: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write moodys page: %w", err)
	}

	return nil
}

// generateWorkWithMePage generates the Work with me page
func generateWorkWithMePage(config Config, tm *TemplateManager) error {
	log.Printf("generateWorkWithMePage")

	outputPath := filepath.Join(config.OutputDir, "work-with-me", "index.html")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create work-with-me directory: %w", err)
	}

	data := TemplateData{
		Title:           "Work with me | Nicholas Skitch",
		CanonicalDomain: config.CanonicalDomain,
		Year:            time.Now().Year(),
	}

	html, err := tm.RenderTemplate("work-with-me.html", data)
	if err != nil {
		return fmt.Errorf("failed to render work-with-me template: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write work-with-me page: %w", err)
	}

	return nil
}

// generateArticlesListing generates the articles listing page
func generateArticlesListing(config Config, tm *TemplateManager, articles map[string]ContentMetadata) error {
	log.Printf("generateArticlesListing: articles count=%d", len(articles))

	outputPath := filepath.Join(config.OutputDir, "articles", "index.html")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create articles directory: %w", err)
	}

	data := TemplateData{
		Title:         "Articles",
		ActivePage:    "Articles",
		Metadata:      articles,
		CanonicalDomain: config.CanonicalDomain,
		AtomFeedURL:   "/atom.xml",
	}

	html, err := tm.RenderTemplate("articles.html", data)
	if err != nil {
		return fmt.Errorf("failed to render articles template: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write articles listing: %w", err)
	}

	return nil
}

// generateArticles generates individual article pages
func generateArticles(config Config, tm *TemplateManager, articles map[string]ContentMetadata) error {
	log.Printf("generateArticles: articles count=%d", len(articles))

	for articleID, metadata := range articles {
		// Read markdown file
		markdown, err := ReadMarkdownFile(config.ContentPath, "articles", articleID)
		if err != nil {
			log.Printf("Failed to read markdown for article %s: %v, skipping", articleID, err)
			continue
		}

		// Convert to HTML
		htmlContent, err := ConvertMarkdownToHTML(markdown)
		if err != nil {
			log.Printf("Failed to convert markdown for article %s: %v, skipping", articleID, err)
			continue
		}

		// Generate article URL
		articleURL := fmt.Sprintf("%s/articles/%s/", config.CanonicalDomain, articleID)

		// Create output directory
		outputPath := filepath.Join(config.OutputDir, "articles", articleID, "index.html")
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("failed to create article directory: %w", err)
		}

		data := TemplateData{
			Title:         metadata.Title,
			ActivePage:    "Articles",
			Article:       metadata,
			Content:       template.HTML(htmlContent),
			ArticleURL:    articleURL,
			ArticleID:     articleID,
			CanonicalDomain: config.CanonicalDomain,
		}

		html, err := tm.RenderTemplate("article-display.html", data)
		if err != nil {
			log.Printf("Failed to render article template for %s: %v, skipping", articleID, err)
			continue
		}

		if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
			return fmt.Errorf("failed to write article: %w", err)
		}
	}

	return nil
}

// generateCodeListing generates the code projects listing page
func generateCodeListing(config Config, tm *TemplateManager, codeProjects map[string]ContentMetadata) error {
	log.Printf("generateCodeListing: code projects count=%d", len(codeProjects))

	outputPath := filepath.Join(config.OutputDir, "code", "index.html")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create code directory: %w", err)
	}

	data := TemplateData{
		Title:         "Code",
		ActivePage:    "Code",
		Metadata:      codeProjects,
		CanonicalDomain: config.CanonicalDomain,
	}

	html, err := tm.RenderTemplate("code.html", data)
	if err != nil {
		return fmt.Errorf("failed to render code template: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write code listing: %w", err)
	}

	return nil
}

// generatePhotosPage generates the photos gallery page
func generatePhotosPage(config Config, tm *TemplateManager) error {
	log.Printf("generatePhotosPage")

	photoPath := filepath.Join(config.StaticDir, "files", "photos")
	entries, err := os.ReadDir(photoPath)
	if err != nil {
		return fmt.Errorf("failed to read photos directory: %w", err)
	}

	var photos []PhotoInfo
	index := 0
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "thumbnails" {
			continue
		}
		if strings.HasSuffix(strings.ToLower(entry.Name()), ".jpg") ||
			strings.HasSuffix(strings.ToLower(entry.Name()), ".jpeg") ||
			strings.HasSuffix(strings.ToLower(entry.Name()), ".png") {
			photos = append(photos, PhotoInfo{
				Filename: entry.Name(),
				Index:    index,
			})
			index++
		}
	}

	outputPath := filepath.Join(config.OutputDir, "photos", "index.html")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create photos directory: %w", err)
	}

	data := TemplateData{
		Title:         "Photos",
		ActivePage:    "Photos",
		Photos:        photos,
		PhotoPath:     "/static/files/photos",
		ThumbnailPath: "/static/files/photos/thumbnails",
		CanonicalDomain: config.CanonicalDomain,
	}

	html, err := tm.RenderTemplate("photos.html", data)
	if err != nil {
		return fmt.Errorf("failed to render photos template: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write photos page: %w", err)
	}

	return nil
}

// generateContactPage generates the contact page
func generateContactPage(config Config, tm *TemplateManager) error {
	log.Printf("generateContactPage")

	outputPath := filepath.Join(config.OutputDir, "contact", "index.html")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create contact directory: %w", err)
	}

	data := TemplateData{
		Title:         "Contact",
		ActivePage:    "Contact",
		CanonicalDomain: config.CanonicalDomain,
	}

	html, err := tm.RenderTemplate("contact.html", data)
	if err != nil {
		return fmt.Errorf("failed to render contact template: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write contact page: %w", err)
	}

	return nil
}

// copyStaticAssets copies static directory to output
func copyStaticAssets(staticDir, outputDir string) error {
	log.Printf("copyStaticAssets: staticDir=%s, outputDir=%s", staticDir, outputDir)

	// Copy everything except content/ directory (we handle that separately)
	return filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(staticDir, path)
		if err != nil {
			return err
		}

		// Skip content directory (we copy article assets separately)
		if strings.HasPrefix(relPath, "content") {
			return nil
		}

		destPath := filepath.Join(outputDir, "static", relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// Read and write file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(destPath, data, info.Mode())
	})
}

// copyArticleAssets copies assets from article directories
func copyArticleAssets(contentPath, outputDir string, articles map[string]ContentMetadata) error {
	log.Printf("copyArticleAssets: articles count=%d", len(articles))

	for articleID := range articles {
		assets, err := GetArticleAssets(contentPath, articleID)
		if err != nil {
			log.Printf("Failed to get assets for article %s: %v, skipping", articleID, err)
			continue
		}

		for _, asset := range assets {
			srcPath := filepath.Join(contentPath, "articles", articleID, asset)
			destPath := filepath.Join(outputDir, "articles", articleID, asset)

			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return fmt.Errorf("failed to create asset directory: %w", err)
			}

			data, err := os.ReadFile(srcPath)
			if err != nil {
				log.Printf("Failed to read asset %s: %v, skipping", srcPath, err)
				continue
			}

			if err := os.WriteFile(destPath, data, 0644); err != nil {
				return fmt.Errorf("failed to write asset: %w", err)
			}
		}
	}

	return nil
}

// copyStaticFiles copies robots.txt and favicon.ico to root
func copyStaticFiles(staticDir, outputDir string) error {
	log.Printf("copyStaticFiles: staticDir=%s, outputDir=%s", staticDir, outputDir)

	files := []string{"robots.txt", "favicon.ico"}
	for _, filename := range files {
		srcPath := filepath.Join(staticDir, filename)
		destPath := filepath.Join(outputDir, filename)

		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			continue // File doesn't exist, skip
		}

		data, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", filename, err)
		}

		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
	}

	return nil
}

