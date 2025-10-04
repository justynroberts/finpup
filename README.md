# finpup - Ultra-Simple CLI Text Editor

finpup is a fast, simple, and powerful terminal-based text editor written in Go. It's designed to be even easier than nano while providing modern features like syntax highlighting, AI-assisted editing, and beautiful color themes.

## Features

- **Ultra-Simple Interface**: Clearer than nano with minimal key bindings
- **Fast Text Processing**: Efficient buffer management for instant response
- **Syntax Highlighting**: Automatic syntax highlighting for Go, Python, JavaScript, JSON, YAML, and more
- **Copy/Paste**: System clipboard integration with fallback to internal clipboard
- **AI Integration**: Built-in AI assistance for code generation and text manipulation (Ollama/OpenAI)
- **Color Themes**: 4 beautiful themes (Dark, Light, Monokai, Solarized) - switch on the fly
- **Pretty Print**: JSON formatting with Ctrl+F
- **Line Numbers**: Always visible for easy navigation

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/justynroberts/finpup.git
cd finpup

# Build
make build

# Install to /usr/local/bin (requires sudo)
make install
```

### Quick Build

```bash
go build -o finpup cmd/finpup/main.go
```

## Usage

```bash
# Edit existing file
finpup myfile.txt

# Create new file
finpup newfile.go

# Just explore
finpup
```

## Key Bindings

| Key       | Action                                    |
|-----------|-------------------------------------------|
| Ctrl+S    | Save file                                 |
| Ctrl+Q    | Quit (press twice if modified)            |
| Ctrl+C    | Copy current line to clipboard            |
| Ctrl+V    | Paste from clipboard                      |
| Ctrl+X    | Cut current line                          |
| Ctrl+K    | Delete current line                       |
| Ctrl+D    | Delete current line (duplicate)           |
| Ctrl+Z    | Undo (50 levels)                          |
| Ctrl+G    | Go to line number                         |
| Ctrl+T    | Jump to top of file                       |
| Ctrl+B    | Jump to bottom of file                    |
| Ctrl+I    | Toggle Insert/Overwrite mode              |
| Ctrl+A    | AI prompt (use ^I/^R/^O in dialog)        |
| Ctrl+H    | Toggle color theme                        |
| Ctrl+F    | Format (currently JSON only)              |
| Arrows    | Navigate                                  |
| Home/End  | Jump to line start/end                    |
| PgUp/PgDn | Page navigation                           |

## AI Integration

finpup supports AI-assisted editing through Ollama (local) or OpenAI-compatible APIs.

### Configuration

Edit `~/.finpup.yaml`:

```yaml
ai:
  enabled: true
  provider: ollama        # "ollama", "openai", or "openrouter"
  base_url: http://localhost:11434
  model: llama3.2
  api_key: ""            # For OpenAI/OpenRouter

theme:
  current: dark          # dark, light, monokai, solarized

editor:
  tab_size: 4
  show_line_numbers: true
  auto_indent: true
```

### OpenRouter Example

```yaml
ai:
  enabled: true
  provider: openrouter
  base_url: https://openrouter.ai
  model: anthropic/claude-3.5-sonnet
  api_key: sk-or-v1-your-api-key-here
```

### Using AI (Ctrl+A)

1. Press **Ctrl+A** to open AI prompt dialog
2. Select mode using Ctrl keys (defaults to INSERT):
   - **Ctrl+I** - INSERT mode: Inserts AI response after current line
   - **Ctrl+R** - REPLACE mode: Replaces current line with AI response
   - **Ctrl+O** - OVERWRITE mode: Overwrites entire buffer with AI response
3. Type your prompt (e.g., "add error handling", "write a hello world program")
4. Press **Enter** to execute or **Esc** to cancel
5. Last prompt is saved - press Ctrl+A then Enter to reuse it

**Example workflow:**
- Press Ctrl+A
- Press Ctrl+R to switch to REPLACE mode
- Type "convert to uppercase"
- Press Enter

### Ollama Setup

```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# Pull a model
ollama pull llama3.2

# Ollama runs on localhost:11434 by default
```

## Color Themes

Press **Ctrl+H** to cycle through themes:

- **Dark**: Classic dark theme with high contrast
- **Light**: Clean light theme for daytime coding
- **Monokai**: Popular colorful theme inspired by Sublime Text
- **Solarized**: Easy-on-the-eyes precision colors

## Development

### Project Structure

```
finpup/
├── cmd/finpup/          # Main entry point
├── internal/
│   ├── buffer/          # Text buffer management
│   ├── editor/          # Core editor logic
│   ├── ui/              # Terminal UI rendering
│   ├── highlight/       # Syntax highlighting
│   ├── ai/              # AI integration
│   └── config/          # Configuration management
└── pkg/
    └── themes/          # Color themes
```

### Running Tests

```bash
make test

# With coverage
make test-coverage
```

### Building

```bash
# Development build
make build

# Run locally
make run

# Clean build artifacts
make clean
```

## Supported File Types

Syntax highlighting is automatically enabled for:

- Go (`.go`)
- Python (`.py`)
- JavaScript/JSX (`.js`, `.jsx`)
- TypeScript/TSX (`.ts`, `.tsx`)
- JSON (`.json`)
- YAML (`.yaml`, `.yml`)
- Markdown (`.md`)
- Shell (`.sh`, `.bash`)

## Troubleshooting

### Clipboard not working

- macOS: Should work out of the box
- Linux: Install `xclip` or `xsel`
- Fallback: Internal clipboard always works

### AI not responding

1. Check Ollama is running: `curl http://localhost:11434`
2. Verify model is installed: `ollama list`
3. Check config: `cat ~/.finpup.yaml`
4. Enable AI: Set `ai.enabled: true`

### Terminal colors look wrong

Try different themes with **Ctrl+H** or set `TERM=xterm-256color`

## License

MIT License - see LICENSE file

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Submit a pull request

## Credits

Built with:
- [tcell](https://github.com/gdamore/tcell) - Terminal handling
- [chroma](https://github.com/alecthomas/chroma) - Syntax highlighting
- [clipboard](https://github.com/atotto/clipboard) - Cross-platform clipboard
