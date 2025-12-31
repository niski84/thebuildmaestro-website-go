package sitegen

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// TemplateData holds data to be passed to templates
type TemplateData struct {
	// Common fields
	Title       string
	ActivePage  string
	CanonicalDomain string
	Year        int
	
	// For article/code listings
	Metadata    map[string]ContentMetadata
	
	// For individual article
	Article     ContentMetadata
	Content     template.HTML
	ArticleURL  string
	ArticleID   string
	
	// For photos
	Photos      []PhotoInfo
	PhotoPath   string
	ThumbnailPath string
	
	// For Atom feed link
	AtomFeedURL string
}

// PhotoInfo represents a photo in the gallery
type PhotoInfo struct {
	Filename string
	Index    int
}

// TemplateManager handles template loading and rendering
type TemplateManager struct {
	templates map[string]*template.Template
	funcMap   template.FuncMap
}

// NewTemplateManager creates a new template manager and loads all templates
func NewTemplateManager(templateDir string) (*TemplateManager, error) {
	log.Printf("NewTemplateManager: templateDir=%s", templateDir)

	tm := &TemplateManager{
		templates: make(map[string]*template.Template),
		funcMap: template.FuncMap{
			"urlFor": func(route string, args ...string) string {
				return urlFor(route, args...)
			},
			"join": strings.Join,
			"split": func(s, sep string) []string {
				return strings.Split(s, sep)
			},
			"hasPrefix": strings.HasPrefix,
			"mod": func(a, b int) int {
				return a % b
			},
			"add": func(a, b int) int {
				return a + b
			},
		},
	}

	// Load base template content once
	basePath := filepath.Join(templateDir, "base.html")
	baseContent, err := os.ReadFile(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read base template: %w", err)
	}

	// Templates that use base template (need inheritance)
	baseTemplates := []string{
		"articles.html",
		"article-display.html",
		"code.html",
		"leanlabs.html",
		"photos.html",
		"contact.html",
	}

	// Standalone templates (full HTML pages)
	standaloneTemplates := []string{
		"home.html",
		"about.html",
		"leanlabs.html",
		"service-saas.html",
		"moodys.html",
		"work-with-me.html",
	}

	// Parse base + each child template separately to avoid block conflicts
	for _, filename := range baseTemplates {
		templatePath := filepath.Join(templateDir, filename)
		content, err := os.ReadFile(templatePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read template %s: %w", filename, err)
		}

		// Create a new template set for this page
		pageTemplates := template.New("").Funcs(tm.funcMap)
		
		// Parse base template
		pageTemplates, err = pageTemplates.Parse(string(baseContent))
		if err != nil {
			return nil, fmt.Errorf("failed to parse base template: %w", err)
		}

		// Parse child template
		pageTemplates, err = pageTemplates.Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", filename, err)
		}

		// Store reference to the base template (which will use this child's blocks)
		tm.templates[filename] = pageTemplates.Lookup("base")
	}

	// Parse standalone templates (they have their own full HTML structure)
	for _, filename := range standaloneTemplates {
		templatePath := filepath.Join(templateDir, filename)
		content, err := os.ReadFile(templatePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read template %s: %w", filename, err)
		}

		// Create a new template set for standalone template
		pageTemplates := template.New("").Funcs(tm.funcMap)
		
		// Parse standalone template (it defines its own template name)
		pageTemplates, err = pageTemplates.Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", filename, err)
		}

		// Store reference to the template (lookup by the define name)
		tm.templates[filename] = pageTemplates.Lookup(filename)
	}

	return tm, nil
}

// RenderTemplate renders a template with the given data
func (tm *TemplateManager) RenderTemplate(templateName string, data TemplateData) (string, error) {
	log.Printf("RenderTemplate: templateName=%s", templateName)

	tmpl, ok := tm.templates[templateName]
	if !ok {
		return "", fmt.Errorf("template not found: %s", templateName)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// urlFor generates URLs for routes (similar to Flask's url_for)
func urlFor(route string, args ...string) string {
	switch route {
	case "display_articles_listing":
		return "/articles/"
	case "display_code_projects_listing":
		return "/code/"
	case "display_photo_gallery":
		return "/photos/"
	case "render_contact_page":
		return "/contact/"
	case "render_homepage_redirect":
		return "/"
	case "render_single_article":
		if len(args) > 0 {
			return fmt.Sprintf("/articles/%s/", args[0])
		}
		return "/articles/"
	case "generate_atom_rss_feed":
		return "/atom.xml"
	default:
		return "/"
	}
}

