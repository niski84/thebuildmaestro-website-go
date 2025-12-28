package sitegen

import (
	"bytes"
	"fmt"
	"log"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

var md goldmark.Markdown

func init() {
	// Initialize goldmark with extensions
	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, // GitHub Flavored Markdown (tables, strikethrough, etc.)
			extension.DefinitionList,
			extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"), // Use GitHub style for code highlighting
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(false),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithHardWraps(),
			goldmarkhtml.WithXHTML(),
		),
	)
}

// ConvertMarkdownToHTML converts markdown content to HTML with syntax highlighting
func ConvertMarkdownToHTML(markdown []byte) (string, error) {
	log.Printf("ConvertMarkdownToHTML: markdown length=%d bytes", len(markdown))

	var buf bytes.Buffer
	if err := md.Convert(markdown, &buf); err != nil {
		return "", fmt.Errorf("failed to convert markdown to HTML: %w", err)
	}

	return buf.String(), nil
}

// GetCodeHighlightCSS returns the CSS for code syntax highlighting
func GetCodeHighlightCSS() (string, error) {
	log.Printf("GetCodeHighlightCSS")

	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}

	formatter := chromahtml.New(chromahtml.WithClasses(true))
	var buf bytes.Buffer
	if err := formatter.WriteCSS(&buf, style); err != nil {
		return "", fmt.Errorf("failed to generate CSS: %w", err)
	}

	return buf.String(), nil
}

// DetectLanguage attempts to detect the programming language from code
func DetectLanguage(code string) string {
	lexer := lexers.Analyse(code)
	if lexer != nil {
		return lexer.Config().Name
	}
	return "text"
}

