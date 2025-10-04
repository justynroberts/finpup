package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/justynroberts/finpup/internal/buffer"
	"github.com/justynroberts/finpup/internal/highlight"
	"github.com/justynroberts/finpup/pkg/themes"
)

type UI struct {
	screen      tcell.Screen
	buffer      *buffer.Buffer
	theme       themes.Theme
	highlighter *highlight.Highlighter
	offsetY     int
	width       int
	height      int
	statusMsg   string
}

func New(buf *buffer.Buffer, theme themes.Theme) (*UI, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	width, height := screen.Size()

	ui := &UI{
		screen:      screen,
		buffer:      buf,
		theme:       theme,
		highlighter: highlight.New(buf.FilePath),
		offsetY:     0,
		width:       width,
		height:      height,
		statusMsg:   "",
	}

	screen.SetStyle(tcell.StyleDefault.
		Background(theme.Background).
		Foreground(theme.Foreground))
	screen.Clear()

	return ui, nil
}

func (ui *UI) Close() {
	ui.screen.Fini()
}

func (ui *UI) SetTheme(theme themes.Theme) {
	ui.theme = theme
	ui.screen.SetStyle(tcell.StyleDefault.
		Background(theme.Background).
		Foreground(theme.Foreground))
}

func (ui *UI) Draw() {
	ui.screen.Clear()
	ui.width, ui.height = ui.screen.Size()

	// Adjust vertical offset to keep cursor visible
	contentHeight := ui.height - 2 // Reserve space for status bars
	if ui.buffer.CursorY < ui.offsetY {
		ui.offsetY = ui.buffer.CursorY
	}
	if ui.buffer.CursorY >= ui.offsetY+contentHeight {
		ui.offsetY = ui.buffer.CursorY - contentHeight + 1
	}

	// Draw lines
	for i := 0; i < contentHeight; i++ {
		lineNum := ui.offsetY + i
		if lineNum >= len(ui.buffer.Lines) {
			break
		}

		ui.drawLine(i, lineNum)
	}

	// Draw status bar
	ui.drawStatusBar()

	// Draw help bar
	ui.drawHelpBar()

	// Position cursor
	screenY := ui.buffer.CursorY - ui.offsetY
	if screenY >= 0 && screenY < contentHeight {
		ui.screen.ShowCursor(ui.buffer.CursorX+4, screenY) // +4 for line numbers
	}

	ui.screen.Show()
}

func (ui *UI) drawLine(screenY, lineNum int) {
	lineNumStr := fmt.Sprintf("%3d ", lineNum+1)
	style := tcell.StyleDefault.
		Background(ui.theme.Background).
		Foreground(ui.theme.LineNumFG)

	for i, r := range lineNumStr {
		ui.screen.SetContent(i, screenY, r, nil, style)
	}

	line := ui.buffer.Lines[lineNum]
	styledRunes, err := ui.highlighter.HighlightLine(line)
	if err != nil || len(styledRunes) == 0 {
		// Fallback to plain text
		style := tcell.StyleDefault.
			Background(ui.theme.Background).
			Foreground(ui.theme.Foreground)
		for i, r := range line {
			if 4+i >= ui.width {
				break
			}
			ui.screen.SetContent(4+i, screenY, r, nil, style)
		}
		return
	}

	for i, sr := range styledRunes {
		if 4+i >= ui.width {
			break
		}
		style := tcell.StyleDefault.
			Background(ui.theme.Background).
			Foreground(sr.Color)
		ui.screen.SetContent(4+i, screenY, sr.Rune, nil, style)
	}
}

func (ui *UI) drawStatusBar() {
	y := ui.height - 2
	style := tcell.StyleDefault.
		Background(ui.theme.StatusBG).
		Foreground(ui.theme.StatusFG)

	// Clear status bar
	for x := 0; x < ui.width; x++ {
		ui.screen.SetContent(x, y, ' ', nil, style)
	}

	modFlag := ""
	if ui.buffer.Modified {
		modFlag = "[+] "
	}

	fileName := ui.buffer.FilePath
	if fileName == "" {
		fileName = "[No Name]"
	}

	status := fmt.Sprintf(" %s%s | Line %d/%d, Col %d",
		modFlag, fileName, ui.buffer.CursorY+1, len(ui.buffer.Lines), ui.buffer.CursorX+1)

	if ui.statusMsg != "" {
		status += " | " + ui.statusMsg
	}

	for i, r := range status {
		if i >= ui.width {
			break
		}
		ui.screen.SetContent(i, y, r, nil, style)
	}
}

func (ui *UI) drawHelpBar() {
	y := ui.height - 1
	style := tcell.StyleDefault.
		Background(ui.theme.Background).
		Foreground(ui.theme.Foreground)

	help := " ^S Save | ^Q Quit | ^K DelLine | ^Z Undo | ^T Top | ^B Bottom | ^I Insert/Ovr | ^G Goto | ^A AI | ^H Theme | ^F Format"

	for x := 0; x < ui.width; x++ {
		ui.screen.SetContent(x, y, ' ', nil, style)
	}

	for i, r := range help {
		if i >= ui.width {
			break
		}
		ui.screen.SetContent(i, y, r, nil, style)
	}
}

func (ui *UI) PollEvent() tcell.Event {
	return ui.screen.PollEvent()
}

func (ui *UI) SetStatus(msg string) {
	ui.statusMsg = msg
}

func (ui *UI) ShowPrompt(prompt string) (string, bool) {
	// Draw prompt popup
	ui.screen.Clear()

	midY := ui.height / 2
	midX := ui.width / 2

	// Draw popup box
	boxWidth := 60
	if boxWidth > ui.width-4 {
		boxWidth = ui.width - 4
	}
	startX := midX - boxWidth/2

	style := tcell.StyleDefault.
		Background(ui.theme.StatusBG).
		Foreground(ui.theme.StatusFG)

	// Draw prompt text
	for i, r := range prompt {
		if i >= boxWidth {
			break
		}
		ui.screen.SetContent(startX+i, midY-1, r, nil, style)
	}

	// Draw input area
	inputY := midY
	for x := startX; x < startX+boxWidth; x++ {
		ui.screen.SetContent(x, inputY, ' ', nil, style)
	}

	ui.screen.ShowCursor(startX, inputY)
	ui.screen.Show()

	// Get input
	var input strings.Builder
	cursorPos := 0

	for {
		ev := ui.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEnter:
				return input.String(), true
			case tcell.KeyEscape:
				return "", false
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if input.Len() > 0 && cursorPos > 0 {
					s := input.String()
					input.Reset()
					input.WriteString(s[:cursorPos-1] + s[cursorPos:])
					cursorPos--
				}
			case tcell.KeyRune:
				s := input.String()
				input.Reset()
				input.WriteString(s[:cursorPos])
				input.WriteRune(ev.Rune())
				input.WriteString(s[cursorPos:])
				cursorPos++
			}

			// Redraw input
			for x := startX; x < startX+boxWidth; x++ {
				ui.screen.SetContent(x, inputY, ' ', nil, style)
			}
			for i, r := range input.String() {
				if i >= boxWidth {
					break
				}
				ui.screen.SetContent(startX+i, inputY, r, nil, style)
			}
			ui.screen.ShowCursor(startX+cursorPos, inputY)
			ui.screen.Show()
		}
	}
}

// ShowAIPrompt shows a prompt that captures Ctrl+I, Ctrl+R, Ctrl+O for mode selection
func (ui *UI) ShowAIPrompt(prompt string) (string, string, bool) {
	// Draw prompt popup
	ui.screen.Clear()

	midY := ui.height / 2
	midX := ui.width / 2

	// Draw popup box
	boxWidth := 70
	if boxWidth > ui.width-4 {
		boxWidth = ui.width - 4
	}
	startX := midX - boxWidth/2

	style := tcell.StyleDefault.
		Background(ui.theme.StatusBG).
		Foreground(ui.theme.StatusFG)

	mode := "insert"
	modeText := "[INSERT]"

	// Get input
	var input strings.Builder
	cursorPos := 0

	for {
		// Draw prompt text with mode
		ui.screen.Clear()
		promptWithMode := prompt + " " + modeText
		for i, r := range promptWithMode {
			if i >= boxWidth {
				break
			}
			ui.screen.SetContent(startX+i, midY-1, r, nil, style)
		}

		// Draw input area
		inputY := midY
		for x := startX; x < startX+boxWidth; x++ {
			ui.screen.SetContent(x, inputY, ' ', nil, style)
		}
		for i, r := range input.String() {
			if i >= boxWidth {
				break
			}
			ui.screen.SetContent(startX+i, inputY, r, nil, style)
		}
		ui.screen.ShowCursor(startX+cursorPos, inputY)
		ui.screen.Show()

		ev := ui.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEnter:
				return input.String(), mode, true
			case tcell.KeyEscape:
				return "", "", false
			case tcell.KeyCtrlI:
				mode = "insert"
				modeText = "[INSERT]"
			case tcell.KeyCtrlR:
				mode = "replace"
				modeText = "[REPLACE]"
			case tcell.KeyCtrlO:
				mode = "overwrite"
				modeText = "[OVERWRITE]"
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if input.Len() > 0 && cursorPos > 0 {
					s := input.String()
					input.Reset()
					input.WriteString(s[:cursorPos-1] + s[cursorPos:])
					cursorPos--
				}
			case tcell.KeyRune:
				s := input.String()
				input.Reset()
				input.WriteString(s[:cursorPos])
				input.WriteRune(ev.Rune())
				input.WriteString(s[cursorPos:])
				cursorPos++
			}
		}
	}
}
