package sitegen

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

// ContentMetadata represents the metadata for a content item (article or code project)
type ContentMetadata struct {
	Title        string
	Author       string
	LastUpdated  string
	WrittenOn    string
	Distributions string
	Description  string
	URL          string
	ID           string
}

// ContentRegistry maps content type to content ID to metadata
type ContentRegistry map[string]map[string]ContentMetadata

// LoadContentRegistry walks the content directories and builds a registry of all content
func LoadContentRegistry(contentPath string) (ContentRegistry, error) {
	log.Printf("LoadContentRegistry: contentPath=%s", contentPath)

	registry := make(ContentRegistry)
	registry["articles"] = make(map[string]ContentMetadata)
	registry["code"] = make(map[string]ContentMetadata)

	contentTypes := []string{"articles", "code"}

	for _, contentType := range contentTypes {
		typePath := filepath.Join(contentPath, contentType)
		
		// Check if directory exists
		if _, err := os.Stat(typePath); os.IsNotExist(err) {
			log.Printf("Content directory does not exist: %s, skipping", typePath)
			continue
		}

		entries, err := os.ReadDir(typePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read content directory %s: %w", typePath, err)
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			contentID := entry.Name()
			contentDir := filepath.Join(typePath, contentID)
			metadataPath := filepath.Join(contentDir, "metadata")

			// Check if metadata file exists
			if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
				log.Printf("Metadata file not found: %s, skipping", metadataPath)
				continue
			}

			metadata, err := parseMetadataFile(metadataPath, contentID)
			if err != nil {
				log.Printf("Failed to parse metadata file %s: %v, skipping", metadataPath, err)
				continue
			}

			registry[contentType][contentID] = metadata
		}
	}

	return registry, nil
}

// parseMetadataFile reads and parses an INI-style metadata file
func parseMetadataFile(metadataPath string, contentID string) (ContentMetadata, error) {
	log.Printf("parseMetadataFile: metadataPath=%s, contentID=%s", metadataPath, contentID)

	cfg, err := ini.Load(metadataPath)
	if err != nil {
		return ContentMetadata{}, fmt.Errorf("failed to load INI file: %w", err)
	}

	metadata := ContentMetadata{
		ID: contentID,
	}

	section := cfg.Section("metadata")
	if section == nil {
		return ContentMetadata{}, fmt.Errorf("metadata section not found in %s", metadataPath)
	}

	// Extract all fields, defaulting to empty string if not present
	metadata.Title = section.Key("title").MustString("")
	metadata.Author = section.Key("author").MustString("")
	metadata.LastUpdated = section.Key("last_updated").MustString("")
	metadata.WrittenOn = section.Key("written_on").MustString("")
	metadata.Distributions = section.Key("distributions").MustString("")
	metadata.Description = section.Key("description").MustString("")
	metadata.URL = section.Key("url").MustString("")

	return metadata, nil
}

// ReadMarkdownFile reads the README.md file from a content directory
func ReadMarkdownFile(contentPath, contentType, contentID string) ([]byte, error) {
	log.Printf("ReadMarkdownFile: contentPath=%s, contentType=%s, contentID=%s", contentPath, contentType, contentID)

	readmePath := filepath.Join(contentPath, contentType, contentID, "README.md")
	
	// Security check: ensure the path is within contentPath
	absContentPath, err := filepath.Abs(contentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute content path: %w", err)
	}
	
	absReadmePath, err := filepath.Abs(readmePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute readme path: %w", err)
	}

	if !strings.HasPrefix(absReadmePath, absContentPath) {
		return nil, fmt.Errorf("readme path outside content directory: %s", readmePath)
	}

	// Sanitize contentID to prevent directory traversal
	sanitizedID := strings.ReplaceAll(contentID, "..", "")
	if sanitizedID != contentID {
		return nil, fmt.Errorf("invalid content ID: %s", contentID)
	}

	data, err := os.ReadFile(readmePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read markdown file: %w", err)
	}

	return data, nil
}

// GetArticleAssets returns a list of asset files in an article directory (excluding README.md and metadata)
func GetArticleAssets(contentPath, articleID string) ([]string, error) {
	log.Printf("GetArticleAssets: contentPath=%s, articleID=%s", contentPath, articleID)

	articleDir := filepath.Join(contentPath, "articles", articleID)
	
	// Security check
	absContentPath, err := filepath.Abs(contentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute content path: %w", err)
	}
	
	absArticleDir, err := filepath.Abs(articleDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute article dir: %w", err)
	}

	if !strings.HasPrefix(absArticleDir, absContentPath) {
		return nil, fmt.Errorf("article path outside content directory: %s", articleDir)
	}

	entries, err := os.ReadDir(articleDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read article directory: %w", err)
	}

	var assets []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		if filename != "README.md" && filename != "metadata" {
			assets = append(assets, filename)
		}
	}

	return assets, nil
}

// FilterPublishedArticles filters out articles with future publish dates
// Articles with WrittenOn date in the future will be excluded
func FilterPublishedArticles(articles map[string]ContentMetadata) map[string]ContentMetadata {
	log.Printf("FilterPublishedArticles: input articles count=%d", len(articles))
	
	now := time.Now()
	filtered := make(map[string]ContentMetadata)
	
	for articleID, metadata := range articles {
		// If WrittenOn is empty, include the article (backward compatibility)
		if metadata.WrittenOn == "" {
			filtered[articleID] = metadata
			continue
		}
		
		// Parse the WrittenOn date
		writtenOn, err := time.Parse("2006-01-02", metadata.WrittenOn)
		if err != nil {
			log.Printf("Failed to parse written_on date '%s' for article %s: %v, including article", metadata.WrittenOn, articleID, err)
			// If date parsing fails, include the article to avoid breaking existing content
			filtered[articleID] = metadata
			continue
		}
		
		// Compare dates (only date, not time)
		writtenOnDate := time.Date(writtenOn.Year(), writtenOn.Month(), writtenOn.Day(), 0, 0, 0, 0, time.UTC)
		nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		
		// Only include articles where WrittenOn is today or in the past
		if !writtenOnDate.After(nowDate) {
			filtered[articleID] = metadata
		} else {
			log.Printf("Excluding article %s (scheduled for future: %s)", articleID, metadata.WrittenOn)
		}
	}
	
	log.Printf("FilterPublishedArticles: output articles count=%d", len(filtered))
	return filtered
}

