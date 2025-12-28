package sitegen

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/feeds"
)

// GenerateAtomFeed generates an Atom XML feed from articles metadata
func GenerateAtomFeed(config Config, articles map[string]ContentMetadata, outputDir string) error {
	log.Printf("GenerateAtomFeed: articles count=%d, outputDir=%s", len(articles), outputDir)

	now := time.Now()
	oldestAllowed := now.AddDate(0, -4, 0) // 4 months ago

	feed := &feeds.Feed{
		Title:       "thebuildmaestro - Articles",
		Link:        &feeds.Link{Href: fmt.Sprintf("%s/atom.xml", config.CanonicalDomain)},
		Description: "Recent articles from thebuildmaestro",
		Author:      &feeds.Author{Name: "Nicholas Skitch"},
		Created:     now,
	}

	var feedItems []*feeds.Item

	for articleID, metadata := range articles {
		// Parse dates
		updated, err := time.Parse("2006-01-02", metadata.LastUpdated)
		if err != nil {
			log.Printf("Failed to parse last_updated for article %s: %v, skipping", articleID, err)
			continue
		}

		published, err := time.Parse("2006-01-02", metadata.WrittenOn)
		if err != nil {
			log.Printf("Failed to parse written_on for article %s: %v, skipping", articleID, err)
			continue
		}

		// Filter articles older than 4 months
		if updated.Before(oldestAllowed) {
			continue
		}

		articleURL := fmt.Sprintf("%s/articles/%s/", config.CanonicalDomain, articleID)

		item := &feeds.Item{
			Title:       metadata.Title,
			Link:        &feeds.Link{Href: articleURL},
			Description: metadata.Description,
			Author:      &feeds.Author{Name: metadata.Author},
			Created:     published,
			Updated:     updated,
		}

		feedItems = append(feedItems, item)
	}

	feed.Items = feedItems

	// Generate Atom XML
	atom, err := feed.ToAtom()
	if err != nil {
		return fmt.Errorf("failed to generate Atom feed: %w", err)
	}

	// Write to file
	outputPath := filepath.Join(outputDir, "atom.xml")
	if err := os.WriteFile(outputPath, []byte(atom), 0644); err != nil {
		return fmt.Errorf("failed to write Atom feed: %w", err)
	}

	log.Printf("Generated Atom feed with %d items", len(feedItems))
	return nil
}

