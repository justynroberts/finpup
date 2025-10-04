# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

finpup is a terminal-based text editor written in Go, designed to be simpler than nano while offering modern features like syntax highlighting, AI integration, and color themes.

## Architecture

### Core Components

1. **Buffer System** (`internal/buffer/`)
   - Text buffer with line-based storage
   - Cursor management and text manipulation
   - File I/O operations
   - All text operations go through the buffer

2. **Editor** (`internal/editor/`)
   - Main editor loop and event handling
   - Keyboard input processing
   - Orchestrates buffer, UI, and AI components
   - Handles all Ctrl+[key] commands

3. **Terminal UI** (`internal/ui/`)
   - Screen rendering using tcell
   - Status bar and help bar display
   - Line number rendering (column 0-3)
   - Text content rendering (column 4+)
   - Modal prompt dialogs

4. **Syntax Highlighting** (`internal/highlight/`)
   - Uses Chroma library for tokenization
   - File type detection by extension
   - Per-line highlighting for performance
   - JSON pretty-printing

5. **AI Integration** (`internal/ai/`)
   - Ollama and OpenAI API support
   - Context-aware text generation
   - Append or replace modes

6. **Themes** (`pkg/themes/`)
   - Color scheme definitions
   - Theme switching logic
   - 4 built-in themes: Dark, Light, Monokai, Solarized

7. **Configuration** (`internal/config/`)
   - YAML-based config at `~/.finpup.yaml`
   - AI settings, theme preferences, editor settings
   - Auto-creates default config if missing

### Key Design Decisions

- **Line-based buffer**: Text stored as `[]string` for simple line operations
- **tcell for UI**: Cross-platform terminal handling with true color support
- **Stateless rendering**: UI redraws completely on each event
- **Vertical offset**: Automatic scrolling keeps cursor visible
- **Modal prompts**: Blocking dialogs for save-as and AI prompts
- **Internal + System clipboard**: Falls back gracefully if clipboard unavailable

## Common Development Tasks

### Building and Testing

```bash
# Build
make build

# Run tests
make test

# Test with coverage
make test-coverage

# Install globally
make install

# Clean
make clean
```

### Running a Single Test

```bash
go test -v ./internal/buffer -run TestInsertRune
```

### Adding a New File Type

1. Add extension mapping in `internal/highlight/DetectLanguage()`
2. Chroma lexer is auto-selected based on file extension

### Adding a New Theme

1. Define theme in `pkg/themes/themes.go`
2. Add to `AllThemes` slice
3. Theme automatically available via Ctrl+H

### Adding a New Keyboard Shortcut

1. Add handler in `internal/editor/handleEvent()`
2. Update help bar in `internal/ui/drawHelpBar()`
3. Update README.md key bindings table

### Adding a New AI Provider

1. Add provider in `internal/ai/GenerateText()` switch
2. Implement provider-specific method (e.g., `generateAnthropic()`)
3. Update config.yaml example and docs

## Testing Strategy

- **Buffer tests**: Core text operations (insert, delete, save, load)
- **Config tests**: Load defaults, verify settings
- **UI/Editor**: Manual testing required (terminal UI)
- **AI**: Mock-friendly design, test with local Ollama

## Key Files

- `cmd/finpup/main.go` - Entry point (30 lines)
- `internal/editor/editor.go` - Main editor logic (250 lines)
- `internal/buffer/buffer.go` - Text buffer (150 lines)
- `internal/ui/ui.go` - Terminal rendering (200 lines)
- `go.mod` - Dependencies (tcell, chroma, clipboard, yaml)

## Code Style

- Go standard formatting (gofmt)
- Keep functions focused and small
- Error handling: return errors, let caller decide
- Comments for exported functions
- Prefer composition over inheritance
- Use stdlib where possible

## Performance Considerations

- Line-by-line highlighting (not whole file)
- Lazy rendering (only visible lines)
- Efficient cursor bounds checking
- Minimal allocations in hot paths (event loop)

## Known Limitations

- No multi-line selection
- Single undo level (buffer state only)
- AI requires network/local LLM
- JSON-only formatting (extend for other languages)
- No UTF-8 width handling for CJK characters (tcell handles basics)
