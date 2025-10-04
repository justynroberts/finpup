# Finton Installation Guide

## Quick Install (Recommended)

```bash
cd /Users/justynroberts/work/finton
./install.sh
```

This will:
1. Build the `finton` binary
2. Install it to `/usr/local/bin/` (requires sudo)
3. Make it available system-wide

## Manual Installation

### Option 1: Build and Install

```bash
# Build
go build -o finton cmd/finton/main.go

# Install globally (requires sudo)
sudo mv finton /usr/local/bin/

# Or install locally
mkdir -p ~/bin
mv finton ~/bin/
# Add ~/bin to PATH if not already
```

### Option 2: Use Make

```bash
make build    # Build binary
make install  # Install to /usr/local/bin (requires sudo)
```

## Verify Installation

```bash
# Check finton is in PATH
which finton

# Run validation
cd /Users/justynroberts/work/finton
./validate.sh
```

## First Run

```bash
# Edit a file
finton test.txt

# Configuration file is auto-created at ~/.finton.yaml
cat ~/.finton.yaml
```

## AI Setup (Optional)

### Using Ollama (Local, Free)

```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# Pull a model
ollama pull llama3.2

# Verify Ollama is running
curl http://localhost:11434

# Enable AI in Finton
cat > ~/.finton.yaml << 'YAML'
ai:
  enabled: true
  provider: ollama
  base_url: http://localhost:11434
  model: llama3.2
  api_key: ""

theme:
  current: dark

editor:
  tab_size: 4
  show_line_numbers: true
  auto_indent: true
YAML
```

### Using OpenAI

```bash
# Edit config
nano ~/.finton.yaml

# Set:
# ai:
#   enabled: true
#   provider: openai
#   base_url: https://api.openai.com
#   model: gpt-4
#   api_key: "your-api-key-here"
```

## Usage

```bash
# Edit existing file
finton myfile.txt

# Create new file
finton newfile.go

# Edit with syntax highlighting
finton code.py
finton data.json
```

## Key Bindings

Once installed, use these shortcuts:

- **Ctrl+S** - Save
- **Ctrl+Q** - Quit
- **Ctrl+C** - Copy line
- **Ctrl+V** - Paste
- **Ctrl+X** - Cut line
- **Ctrl+A** - AI prompt (if enabled)
- **Ctrl+H** - Toggle theme
- **Ctrl+F** - Format JSON

## Troubleshooting

### Command not found

```bash
# Check if /usr/local/bin is in PATH
echo $PATH | grep /usr/local/bin

# If not, add to ~/.zshrc or ~/.bashrc:
export PATH="/usr/local/bin:$PATH"
source ~/.zshrc  # or ~/.bashrc
```

### Permission denied during install

```bash
# Use sudo for system-wide install
sudo ./install.sh

# Or install to user directory
mkdir -p ~/bin
go build -o ~/bin/finton cmd/finton/main.go
export PATH="$HOME/bin:$PATH"
```

### Binary won't run

```bash
# Make sure it's executable
chmod +x finton

# Check architecture
file finton
# Should show: Mach-O 64-bit executable arm64
```

## Uninstall

```bash
# Remove binary
sudo rm /usr/local/bin/finton

# Remove config (optional)
rm ~/.finton.yaml
```

## Next Steps

1. Read [QUICKSTART.md](QUICKSTART.md) for features
2. Read [README.md](README.md) for full documentation
3. Try the demo: `./demo.sh`
4. Edit a file: `finton test.txt`
