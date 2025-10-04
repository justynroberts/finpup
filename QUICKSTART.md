# Finton Quick Start

## Installation (2 steps)

```bash
cd /Users/justynroberts/work/finton
./install.sh
```

Or manually:
```bash
go build -o finton cmd/finton/main.go
sudo mv finton /usr/local/bin/
```

## Usage

```bash
# Edit a file
finton myfile.txt

# Create new file
finton newfile.go
```

## Essential Keys

- **Ctrl+S** - Save
- **Ctrl+Q** - Quit
- **Ctrl+C** - Copy line
- **Ctrl+V** - Paste
- **Ctrl+X** - Cut line
- **Ctrl+A** - AI prompt
- **Ctrl+H** - Change theme
- **Ctrl+F** - Format JSON

## AI Setup (Optional)

### Using Ollama (Local, Free)

```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# Pull a model
ollama pull llama3.2
```

### Enable AI

Edit `~/.finton.yaml`:

```yaml
ai:
  enabled: true
  provider: ollama
  base_url: http://localhost:11434
  model: llama3.2
```

### Using AI in Finton

1. Press **Ctrl+A**
2. Type prompt: "add error handling"
3. Press **Enter**
4. AI response appears below cursor

**Replace mode**: Type "replace: convert to uppercase"

## Themes

Press **Ctrl+H** to cycle through:
- Dark
- Light
- Monokai
- Solarized

## Tips

- Line numbers always shown on left
- Auto syntax highlighting based on file extension
- Status bar shows file name, cursor position, modifications
- Help bar at bottom shows all commands
- System clipboard integration (copies also to internal buffer as fallback)

## Supported Languages

Auto-highlighting for: Go, Python, JavaScript, TypeScript, JSON, YAML, Markdown, Shell

## Test It

```bash
finton test.txt
```

Try editing, then:
- Press Ctrl+H to change theme
- Press Ctrl+C to copy a line
- Press Ctrl+V to paste
- Press Ctrl+S to save
- Press Ctrl+Q to quit
