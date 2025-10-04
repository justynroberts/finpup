package highlight

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gdamore/tcell/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"gopkg.in/yaml.v3"
)

type Highlighter struct {
	lexer     chroma.Lexer
	formatter chroma.Formatter
	style     *chroma.Style
}

func New(filePath string) *Highlighter {
	lexer := lexers.Match(filePath)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	return &Highlighter{
		lexer:     lexer,
		formatter: formatters.TTY256,
		style:     styles.Get("monokai"),
	}
}

func (h *Highlighter) HighlightLine(line string) ([]StyledRune, error) {
	iterator, err := h.lexer.Tokenise(nil, line)
	if err != nil {
		return runesFromString(line, tcell.ColorWhite), nil
	}

	var result []StyledRune
	for token := iterator(); token != chroma.EOF; token = iterator() {
		style := h.style.Get(token.Type)
		color := tcell.ColorWhite

		if style.Colour.IsSet() {
			color = tcell.GetColor(style.Colour.String())
		}

		for _, r := range token.Value {
			result = append(result, StyledRune{
				Rune:  r,
				Color: color,
			})
		}
	}

	return result, nil
}

type StyledRune struct {
	Rune  rune
	Color tcell.Color
}

func runesFromString(s string, color tcell.Color) []StyledRune {
	result := make([]StyledRune, 0, len(s))
	for _, r := range s {
		result = append(result, StyledRune{Rune: r, Color: color})
	}
	return result
}

func FormatJSON(text string) (string, error) {
	// Simple JSON formatter
	var result strings.Builder
	indent := 0
	inString := false

	for i := 0; i < len(text); i++ {
		c := text[i]

		if c == '"' && (i == 0 || text[i-1] != '\\') {
			inString = !inString
		}

		if !inString {
			switch c {
			case '{', '[':
				result.WriteByte(c)
				result.WriteByte('\n')
				indent++
				result.WriteString(strings.Repeat("  ", indent))
			case '}', ']':
				result.WriteByte('\n')
				indent--
				result.WriteString(strings.Repeat("  ", indent))
				result.WriteByte(c)
			case ',':
				result.WriteByte(c)
				result.WriteByte('\n')
				result.WriteString(strings.Repeat("  ", indent))
			case ':':
				result.WriteString(": ")
			case ' ', '\n', '\r', '\t':
				// Skip whitespace
			default:
				result.WriteByte(c)
			}
		} else {
			result.WriteByte(c)
		}
	}

	return result.String(), nil
}

func FormatYAML(text string) (string, error) {
	var data interface{}
	if err := yaml.Unmarshal([]byte(text), &data); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(data); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

func FormatHCL(text string) (string, error) {
	formatted := hclwrite.Format([]byte(text))
	return string(formatted), nil
}

func DetectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".js", ".jsx":
		return "javascript"
	case ".ts", ".tsx":
		return "typescript"
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".md":
		return "markdown"
	case ".sh", ".bash":
		return "bash"
	case ".hcl", ".tf":
		return "hcl"
	default:
		return "text"
	}
}
