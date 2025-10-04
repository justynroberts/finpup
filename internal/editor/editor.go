package editor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/justynroberts/finpup/internal/ai"
	"github.com/justynroberts/finpup/internal/buffer"
	"github.com/justynroberts/finpup/internal/config"
	"github.com/justynroberts/finpup/internal/highlight"
	"github.com/justynroberts/finpup/internal/ui"
)

type Editor struct {
	buffer       *buffer.Buffer
	ui           *ui.UI
	config       *config.Config
	aiClient     *ai.Client
	clipboard    string
	running      bool
	undoStack    [][]string
	aiPromptHistory []string
	lastAIPrompt string
	insertMode   bool
}

func New(filePath string) (*Editor, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	buf, err := buffer.New(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create buffer: %w", err)
	}

	ui, err := ui.New(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create UI: %w", err)
	}

	return &Editor{
		buffer:       buf,
		ui:           ui,
		config:       cfg,
		aiClient:     ai.New(&cfg.AI),
		clipboard:    "",
		running:      true,
		undoStack:    make([][]string, 0, 50),
		aiPromptHistory: make([]string, 0, 20),
		lastAIPrompt: "",
		insertMode:   true,
	}, nil
}

func (e *Editor) Run() error {
	defer e.ui.Close()

	e.ui.Draw()

	for e.running {
		ev := e.ui.PollEvent()
		e.handleEvent(ev)
		e.ui.Draw()
	}

	return nil
}

func (e *Editor) handleEvent(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventResize:
		e.ui.Draw()

	case *tcell.EventKey:
		if ev.Key() == tcell.KeyCtrlS {
			e.handleSave()
		} else if ev.Key() == tcell.KeyCtrlQ {
			e.handleQuit()
		} else if ev.Key() == tcell.KeyCtrlC {
			e.handleCopy()
		} else if ev.Key() == tcell.KeyCtrlV {
			e.handlePaste()
		} else if ev.Key() == tcell.KeyCtrlX {
			e.handleCut()
		} else if ev.Key() == tcell.KeyCtrlA {
			e.handleAI()
		} else if ev.Key() == tcell.KeyCtrlE {
			e.handleEmojiPicker()
		} else if ev.Key() == tcell.KeyCtrlW {
			e.handleToggleSelection()
		} else if ev.Key() == tcell.KeyCtrlF {
			e.handleFormat()
		} else if ev.Key() == tcell.KeyCtrlD {
			e.handleDeleteLine()
		} else if ev.Key() == tcell.KeyCtrlG {
			e.handleGoToLine()
		} else if ev.Key() == tcell.KeyCtrlZ {
			e.handleUndo()
		} else if ev.Key() == tcell.KeyCtrlK {
			e.handleDeleteLine()
		} else if ev.Key() == tcell.KeyCtrlT {
			e.handleJumpToTop()
		} else if ev.Key() == tcell.KeyCtrlB {
			e.handleJumpToBottom()
		} else if ev.Key() == tcell.KeyCtrlI {
			e.toggleInsertMode()
		} else if ev.Key() == tcell.KeyEnter {
			e.saveUndo()
			e.buffer.InsertNewline()
		} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
			e.saveUndo()
			e.buffer.DeleteRune()
		} else if ev.Key() == tcell.KeyUp {
			e.moveCursorUp()
		} else if ev.Key() == tcell.KeyDown {
			e.moveCursorDown()
		} else if ev.Key() == tcell.KeyLeft {
			e.moveCursorLeft()
		} else if ev.Key() == tcell.KeyRight {
			e.moveCursorRight()
		} else if ev.Key() == tcell.KeyHome {
			e.buffer.CursorX = 0
		} else if ev.Key() == tcell.KeyEnd {
			e.buffer.CursorX = len(e.buffer.GetCurrentLine())
		} else if ev.Key() == tcell.KeyPgUp {
			e.pageUp()
		} else if ev.Key() == tcell.KeyPgDn {
			e.pageDown()
		} else if ev.Key() == tcell.KeyRune {
			e.saveUndo()
			if e.insertMode {
				e.buffer.InsertRune(ev.Rune())
			} else {
				e.buffer.OverwriteRune(ev.Rune())
			}
		}
	}
}

func (e *Editor) handleSave() {
	if e.buffer.FilePath == "" {
		prompt, ok := e.ui.ShowPrompt("Save as: ")
		if !ok || prompt == "" {
			e.ui.SetStatus("Save cancelled")
			return
		}
		e.buffer.FilePath = prompt
	}

	if err := e.buffer.Save(); err != nil {
		e.ui.SetStatus(fmt.Sprintf("Error saving: %v", err))
		return
	}

	e.ui.SetStatus(fmt.Sprintf("Saved to %s", e.buffer.FilePath))
}

func (e *Editor) handleQuit() {
	if e.buffer.Modified {
		// Force quit on second Ctrl+Q
		e.buffer.Modified = false
		e.ui.SetStatus("File modified! Press Ctrl+Q again to force quit or Ctrl+S to save")
	} else {
		e.running = false
	}
}

func (e *Editor) handleCopy() {
	line := e.buffer.GetCurrentLine()
	e.clipboard = line

	// Try to copy to system clipboard
	if err := clipboard.WriteAll(line); err == nil {
		e.ui.SetStatus("Copied to clipboard")
	} else {
		e.ui.SetStatus("Copied to internal clipboard")
	}
}

func (e *Editor) handlePaste() {
	// Try system clipboard first
	text, err := clipboard.ReadAll()
	if err == nil && text != "" {
		e.buffer.InsertText(text)
		e.ui.SetStatus("Pasted from clipboard")
		return
	}

	// Fall back to internal clipboard
	if e.clipboard != "" {
		e.buffer.InsertText(e.clipboard)
		e.ui.SetStatus("Pasted from internal clipboard")
	}
}

func (e *Editor) handleCut() {
	line := e.buffer.DeleteCurrentLine()
	e.clipboard = line

	if err := clipboard.WriteAll(line); err == nil {
		e.ui.SetStatus("Cut to clipboard")
	} else {
		e.ui.SetStatus("Cut to internal clipboard")
	}
}

func (e *Editor) handleAI() {
	if !e.config.AI.Enabled {
		e.ui.SetStatus("AI disabled. Edit ~/.finpup.yaml to enable")
		return
	}

	// Determine context: selection, or whole document
	var context string
	var contextType string
	if e.buffer.HasSelection() {
		context = e.buffer.GetSelection()
		contextType = "selection"
	} else {
		context = e.buffer.GetAllText()
		contextType = "document"
	}

	// Show prompt with mode selection
	promptHint := fmt.Sprintf("AI Prompt [%s] (^I=insert ^R=replace ^O=overwrite):", contextType)

	prompt, mode, ok := e.ui.ShowAIPrompt(promptHint)
	if !ok {
		e.ui.SetStatus("AI cancelled")
		return
	}

	// If just pressed enter with no input, use last prompt
	if prompt == "" && e.lastAIPrompt != "" {
		prompt = e.lastAIPrompt
	}

	if prompt == "" {
		e.ui.SetStatus("AI cancelled - no prompt")
		return
	}

	// Save to history
	e.aiPromptHistory = append(e.aiPromptHistory, prompt)
	if len(e.aiPromptHistory) > 20 {
		e.aiPromptHistory = e.aiPromptHistory[1:]
	}
	e.lastAIPrompt = prompt

	e.ui.SetStatus("Generating AI response...")
	e.ui.Draw()

	result, err := e.aiClient.GenerateText(prompt, context)
	if err != nil {
		e.ui.SetStatus(fmt.Sprintf("AI error: %v", err))
		return
	}

	e.saveUndo()

	switch mode {
	case "replace":
		if e.buffer.HasSelection() {
			e.buffer.ReplaceSelection(result)
			e.ui.SetStatus("Selection replaced with AI response")
		} else {
			// Replace entire buffer
			e.buffer.Lines = strings.Split(result, "\n")
			e.buffer.CursorY = 0
			e.buffer.CursorX = 0
			e.buffer.Modified = true
			e.ui.SetStatus("Document replaced with AI response")
		}
	case "overwrite":
		// Clear entire buffer and insert AI response
		e.buffer.Lines = strings.Split(result, "\n")
		e.buffer.CursorY = 0
		e.buffer.CursorX = 0
		e.buffer.Modified = true
		e.ui.SetStatus("Buffer overwritten with AI response")
	case "insert":
		e.buffer.InsertText("\n" + result)
		e.ui.SetStatus("AI response inserted")
	}

	// Clear selection after AI operation
	e.buffer.SelectMode = false
}

func (e *Editor) handleEmojiPicker() {
	emoji, ok := e.ui.ShowEmojiPicker()
	if !ok || emoji == "" {
		e.ui.SetStatus("Emoji cancelled")
		return
	}

	e.saveUndo()
	for _, r := range emoji {
		e.buffer.InsertRune(r)
	}
	e.ui.SetStatus(fmt.Sprintf("Inserted emoji: %s", emoji))
}

func (e *Editor) handleFormat() {
	e.saveUndo()
	ext := filepath.Ext(e.buffer.FilePath)
	text := e.buffer.GetAllText()
	var formatted string
	var err error

	switch ext {
	case ".json":
		formatted, err = highlight.FormatJSON(text)
		if err != nil {
			e.ui.SetStatus(fmt.Sprintf("JSON format error: %v", err))
			return
		}
		e.ui.SetStatus("Formatted JSON")
	case ".yaml", ".yml":
		formatted, err = highlight.FormatYAML(text)
		if err != nil {
			e.ui.SetStatus(fmt.Sprintf("YAML format error: %v", err))
			return
		}
		e.ui.SetStatus("Formatted YAML")
	case ".hcl", ".tf":
		formatted, err = highlight.FormatHCL(text)
		if err != nil {
			e.ui.SetStatus(fmt.Sprintf("HCL format error: %v", err))
			return
		}
		e.ui.SetStatus("Formatted HCL")
	default:
		e.ui.SetStatus("Format supports JSON, YAML, and HCL/Terraform files")
		return
	}

	lines := strings.Split(formatted, "\n")
	e.buffer.Lines = lines
	e.buffer.Modified = true
}

func (e *Editor) moveCursorUp() {
	if e.buffer.CursorY > 0 {
		e.buffer.CursorY--
		e.adjustCursorX()
	}
}

func (e *Editor) moveCursorDown() {
	if e.buffer.CursorY < len(e.buffer.Lines)-1 {
		e.buffer.CursorY++
		e.adjustCursorX()
	}
}

func (e *Editor) moveCursorLeft() {
	if e.buffer.CursorX > 0 {
		e.buffer.CursorX--
	} else if e.buffer.CursorY > 0 {
		e.buffer.CursorY--
		e.buffer.CursorX = len(e.buffer.GetCurrentLine())
	}
}

func (e *Editor) moveCursorRight() {
	lineLen := len(e.buffer.GetCurrentLine())
	if e.buffer.CursorX < lineLen {
		e.buffer.CursorX++
	} else if e.buffer.CursorY < len(e.buffer.Lines)-1 {
		e.buffer.CursorY++
		e.buffer.CursorX = 0
	}
}

func (e *Editor) adjustCursorX() {
	lineLen := len(e.buffer.GetCurrentLine())
	if e.buffer.CursorX > lineLen {
		e.buffer.CursorX = lineLen
	}
}

func (e *Editor) pageUp() {
	e.buffer.CursorY -= 10
	if e.buffer.CursorY < 0 {
		e.buffer.CursorY = 0
	}
	e.adjustCursorX()
}

func (e *Editor) pageDown() {
	e.buffer.CursorY += 10
	if e.buffer.CursorY >= len(e.buffer.Lines) {
		e.buffer.CursorY = len(e.buffer.Lines) - 1
	}
	e.adjustCursorX()
}

func (e *Editor) handleDeleteLine() {
	e.saveUndo()
	line := e.buffer.DeleteCurrentLine()
	e.clipboard = line
	e.ui.SetStatus("Line deleted (in clipboard)")
}

func (e *Editor) handleGoToLine() {
	input, ok := e.ui.ShowPrompt("Go to line: ")
	if !ok || input == "" {
		return
	}

	var lineNum int
	if _, err := fmt.Sscanf(input, "%d", &lineNum); err != nil {
		e.ui.SetStatus("Invalid line number")
		return
	}

	lineNum-- // Convert to 0-based
	if lineNum < 0 {
		lineNum = 0
	}
	if lineNum >= len(e.buffer.Lines) {
		lineNum = len(e.buffer.Lines) - 1
	}

	e.buffer.CursorY = lineNum
	e.buffer.CursorX = 0
	e.ui.SetStatus(fmt.Sprintf("Jumped to line %d", lineNum+1))
}

func (e *Editor) saveUndo() {
	// Deep copy current buffer state
	snapshot := make([]string, len(e.buffer.Lines))
	copy(snapshot, e.buffer.Lines)

	e.undoStack = append(e.undoStack, snapshot)

	// Limit undo stack size
	if len(e.undoStack) > 50 {
		e.undoStack = e.undoStack[1:]
	}
}

func (e *Editor) handleUndo() {
	if len(e.undoStack) == 0 {
		e.ui.SetStatus("Nothing to undo")
		return
	}

	// Pop last state
	lastState := e.undoStack[len(e.undoStack)-1]
	e.undoStack = e.undoStack[:len(e.undoStack)-1]

	e.buffer.Lines = lastState
	e.buffer.Modified = true

	// Adjust cursor if needed
	if e.buffer.CursorY >= len(e.buffer.Lines) {
		e.buffer.CursorY = len(e.buffer.Lines) - 1
	}
	e.adjustCursorX()

	e.ui.SetStatus("Undo successful")
}

func (e *Editor) handleJumpToTop() {
	e.buffer.CursorY = 0
	e.buffer.CursorX = 0
	e.ui.SetStatus("Jumped to top")
}

func (e *Editor) handleJumpToBottom() {
	e.buffer.CursorY = len(e.buffer.Lines) - 1
	e.buffer.CursorX = 0
	e.ui.SetStatus("Jumped to bottom")
}

func (e *Editor) toggleInsertMode() {
	e.insertMode = !e.insertMode
	if e.insertMode {
		e.ui.SetStatus("INSERT mode")
	} else {
		e.ui.SetStatus("OVERWRITE mode")
	}
}

func (e *Editor) handleToggleSelection() {
	e.buffer.ToggleSelection()
	if e.buffer.SelectMode {
		e.ui.SetStatus("Selection mode ON - move cursor to select")
	} else {
		e.ui.SetStatus("Selection mode OFF")
	}
}
