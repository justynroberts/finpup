# finpup

A fast, simple terminal text editor with AI integration and syntax highlighting.

## Quick Start

### Installation

Download the latest release for your platform from [Releases](https://github.com/justynroberts/finpup/releases).

**macOS/Linux:**
```bash
curl -L https://github.com/justynroberts/finpup/releases/latest/download/finpup -o finpup
chmod +x finpup
sudo mv finpup /usr/local/bin/
```

**From source:**
```bash
git clone https://github.com/justynroberts/finpup.git
cd finpup
make build
sudo make install
```

### Usage

```bash
finpup filename.txt    # Edit a file
finpup                 # Start with empty buffer
```

## Features

- **Simple Interface**: Easier than nano with clear key bindings
- **Syntax Highlighting**: Automatic highlighting for Go, Python, JavaScript, JSON, YAML, and more
- **AI Integration**: Built-in AI assistance via Ollama or OpenAI-compatible APIs
- **Color Themes**: Dark, Light, Monokai, Solarized (switch with Ctrl+H)
- **Clipboard Support**: System clipboard integration with internal fallback
- **JSON Formatting**: Pretty-print JSON with Ctrl+F
- **Undo Support**: 50 levels of undo with Ctrl+Z

## Key Bindings

| Key       | Action                                    |
|-----------|-------------------------------------------|
| Ctrl+S    | Save file                                 |
| Ctrl+Q    | Quit (press twice if modified)            |
| Ctrl+C    | Copy current line                         |
| Ctrl+V    | Paste                                     |
| Ctrl+X    | Cut current line                          |
| Ctrl+K    | Delete current line                       |
| Ctrl+Z    | Undo                                      |
| Ctrl+G    | Go to line number                         |
| Ctrl+T    | Jump to top                               |
| Ctrl+B    | Jump to bottom                            |
| Ctrl+A    | AI prompt                                 |
| Ctrl+H    | Toggle theme                              |
| Ctrl+F    | Format JSON                               |
| Arrows    | Navigate                                  |

## AI Configuration

Create `~/.finpup.yaml`:

```yaml
ai:
  enabled: true
  provider: ollama              # "ollama", "openai", or "openrouter"
  base_url: http://localhost:11434
  model: llama3.2
  api_key: ""                   # Required for OpenAI/OpenRouter

theme:
  current: dark                 # dark, light, monokai, solarized

editor:
  tab_size: 4
  show_line_numbers: true
  auto_indent: true
```

### Ollama Setup

```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# Pull a model
ollama pull llama3.2
```

### Using AI (Ctrl+A)

1. Press **Ctrl+A** to open prompt
2. Select mode:
   - **Ctrl+I** - INSERT: Add after current line
   - **Ctrl+R** - REPLACE: Replace current line
   - **Ctrl+O** - OVERWRITE: Replace entire buffer
3. Type your prompt and press **Enter**
4. Press **Esc** to cancel

Example: Press Ctrl+A, then Ctrl+R, type "convert to uppercase", press Enter.

## Supported Languages

Syntax highlighting for: Go, Python, JavaScript, TypeScript, JSON, YAML, Markdown, Shell, JSX, TSX

## Development

### Build from Source

```bash
# Clone repository
git clone https://github.com/justynroberts/finpup.git
cd finpup

# Install dependencies
go mod download

# Build
make build

# Run tests
make test

# Install system-wide
sudo make install
```

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
│   └── config/          # Configuration
└── pkg/
    └── themes/          # Color themes
```

## Troubleshooting

**Clipboard not working:**
- macOS: Works out of the box
- Linux: Install `xclip` or `xsel`
- Fallback internal clipboard always available

**AI not responding:**
1. Check Ollama: `curl http://localhost:11434`
2. Verify model: `ollama list`
3. Check config: `cat ~/.finpup.yaml`
4. Ensure `ai.enabled: true`

**Wrong colors:**
- Try different themes with Ctrl+H
- Set `TERM=xterm-256color`

## License

MIT License

## Contributing

Contributions welcome! Please fork the repository, create a feature branch, write tests, and submit a pull request.
