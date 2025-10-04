package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/justynroberts/finpup/internal/buffer"
	"github.com/justynroberts/finpup/internal/highlight"
)

type UI struct {
	screen      tcell.Screen
	buffer      *buffer.Buffer
	highlighter *highlight.Highlighter
	offsetY     int
	width       int
	height      int
	statusMsg   string
}

func New(buf *buffer.Buffer) (*UI, error) {
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
		highlighter: highlight.New(buf.FilePath),
		offsetY:     0,
		width:       width,
		height:      height,
		statusMsg:   "",
	}

	screen.SetStyle(tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite))
	screen.Clear()

	return ui, nil
}

func (ui *UI) Close() {
	ui.screen.Fini()
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
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorGray)

	for i, r := range lineNumStr {
		ui.screen.SetContent(i, screenY, r, nil, style)
	}

	line := ui.buffer.Lines[lineNum]
	styledRunes, err := ui.highlighter.HighlightLine(line)
	if err != nil || len(styledRunes) == 0 {
		// Fallback to plain text
		style := tcell.StyleDefault.
			Background(tcell.ColorBlack).
			Foreground(tcell.ColorWhite)
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
			Background(tcell.ColorBlack).
			Foreground(sr.Color)
		ui.screen.SetContent(4+i, screenY, sr.Rune, nil, style)
	}
}

func (ui *UI) drawStatusBar() {
	y := ui.height - 2
	style := tcell.StyleDefault.
		Background(tcell.ColorBlue).
		Foreground(tcell.ColorWhite)

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
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)

	help := " ^S Save | ^Q Quit | ^K Del | ^Z Undo | ^T Top | ^B Bottom | ^W Select | ^A AI | ^E Emoji | ^F Format"

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
		Background(tcell.ColorBlue).
		Foreground(tcell.ColorWhite)

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
		Background(tcell.ColorBlue).
		Foreground(tcell.ColorWhite)

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

// ShowEmojiPicker shows a scrollable emoji picker
func (ui *UI) ShowEmojiPicker() (string, bool) {
	emojis := []struct {
		emoji string
		name  string
	}{
		{"😀", "grinning"}, {"😃", "smiley"}, {"😄", "smile"}, {"😁", "grin"},
		{"😆", "laughing"}, {"😅", "sweat_smile"}, {"🤣", "rofl"}, {"😂", "joy"},
		{"🙂", "slightly_smiling"}, {"🙃", "upside_down"}, {"😉", "wink"}, {"😊", "blush"},
		{"😇", "innocent"}, {"🥰", "smiling_hearts"}, {"😍", "heart_eyes"}, {"🤩", "star_struck"},
		{"😘", "kissing_heart"}, {"😗", "kissing"}, {"😚", "kissing_closed_eyes"}, {"😙", "kissing_smiling_eyes"},
		{"🥲", "smiling_tear"}, {"😋", "yum"}, {"😛", "stuck_out_tongue"}, {"😜", "stuck_out_tongue_winking"},
		{"🤪", "zany"}, {"😝", "stuck_out_tongue_closed_eyes"}, {"🤑", "money_mouth"}, {"🤗", "hugs"},
		{"🤭", "hand_over_mouth"}, {"🤫", "shushing"}, {"🤔", "thinking"}, {"🤐", "zipper_mouth"},
		{"🤨", "raised_eyebrow"}, {"😐", "neutral"}, {"😑", "expressionless"}, {"😶", "no_mouth"},
		{"😏", "smirk"}, {"😒", "unamused"}, {"🙄", "roll_eyes"}, {"😬", "grimacing"},
		{"🤥", "lying"}, {"😌", "relieved"}, {"😔", "pensive"}, {"😪", "sleepy"},
		{"🤤", "drooling"}, {"😴", "sleeping"}, {"😷", "mask"}, {"🤒", "thermometer"},
		{"🤕", "head_bandage"}, {"🤢", "nauseated"}, {"🤮", "vomiting"}, {"🤧", "sneezing"},
		{"🥵", "hot"}, {"🥶", "cold"}, {"😵", "dizzy"}, {"🤯", "exploding_head"},
		{"😎", "sunglasses"}, {"🤓", "nerd"}, {"🧐", "monocle"}, {"😕", "confused"},
		{"😟", "worried"}, {"🙁", "frowning"}, {"☹️", "frowning2"}, {"😮", "open_mouth"},
		{"😯", "hushed"}, {"😲", "astonished"}, {"😳", "flushed"}, {"🥺", "pleading"},
		{"😦", "frowning_open"}, {"😧", "anguished"}, {"😨", "fearful"}, {"😰", "anxious_sweat"},
		{"😥", "sad_sweat"}, {"😢", "cry"}, {"😭", "sob"}, {"😱", "scream"},
		{"😖", "confounded"}, {"😣", "persevere"}, {"😞", "disappointed"}, {"😓", "sweat"},
		{"😩", "weary"}, {"😫", "tired"}, {"🥱", "yawn"}, {"😤", "triumph"},
		{"😡", "rage"}, {"😠", "angry"}, {"🤬", "cursing"}, {"👍", "thumbsup"},
		{"👎", "thumbsdown"}, {"👌", "ok_hand"}, {"✌️", "victory"}, {"🤞", "crossed_fingers"},
		{"🤟", "love_you"}, {"🤘", "metal"}, {"👋", "wave"}, {"🤚", "raised_back_hand"},
		{"👏", "clap"}, {"🙌", "raised_hands"}, {"👐", "open_hands"}, {"🤲", "palms_up"},
		{"🙏", "pray"}, {"✍️", "writing"}, {"💪", "muscle"}, {"🦾", "mechanical_arm"},
		{"❤️", "heart"}, {"🧡", "orange_heart"}, {"💛", "yellow_heart"}, {"💚", "green_heart"},
		{"💙", "blue_heart"}, {"💜", "purple_heart"}, {"🖤", "black_heart"}, {"🤍", "white_heart"},
		{"💔", "broken_heart"}, {"❤️‍🔥", "heart_on_fire"}, {"💯", "100"}, {"💢", "anger"},
		{"💥", "boom"}, {"💫", "dizzy_symbol"}, {"💦", "sweat_drops"}, {"💨", "dash"},
		{"🔥", "fire"}, {"✨", "sparkles"}, {"⭐", "star"}, {"🌟", "star2"},
		{"💤", "zzz"}, {"🚀", "rocket"}, {"🎉", "tada"}, {"🎊", "confetti"},
		{"✅", "check"}, {"❌", "x"}, {"⚠️", "warning"}, {"🔔", "bell"},
		{"📌", "pin"}, {"📍", "location"}, {"💡", "bulb"}, {"🔒", "lock"},
		{"🔓", "unlock"}, {"🔑", "key"}, {"🎯", "dart"}, {"💰", "moneybag"},
	}

	selected := 0
	offset := 0

	for {
		// Clear with black background
		defaultStyle := tcell.StyleDefault.
			Background(tcell.ColorBlack).
			Foreground(tcell.ColorWhite)
		for y := 0; y < ui.height; y++ {
			for x := 0; x < ui.width; x++ {
				ui.screen.SetContent(x, y, ' ', nil, defaultStyle)
			}
		}

		midY := ui.height / 2
		midX := ui.width / 2

		boxWidth := 70
		boxHeight := 15
		if boxWidth > ui.width-4 {
			boxWidth = ui.width - 4
		}
		if boxHeight > ui.height-4 {
			boxHeight = ui.height - 4
		}
		startX := midX - boxWidth/2
		startY := midY - boxHeight/2

		style := tcell.StyleDefault.
			Background(tcell.ColorBlack).
			Foreground(tcell.ColorWhite)

		selectedStyle := tcell.StyleDefault.
			Background(tcell.ColorWhite).
			Foreground(tcell.ColorBlack)

		// Title
		title := " Emoji Picker (↑↓ navigate, Enter select, Esc cancel) "
		for i, r := range title {
			if i >= boxWidth {
				break
			}
			ui.screen.SetContent(startX+i, startY, r, nil, style)
		}

		// Display emojis
		perRow := (boxWidth - 2) / 8
		visibleRows := boxHeight - 2

		if offset > selected/perRow {
			offset = selected / perRow
		}
		if selected/perRow >= offset+visibleRows {
			offset = selected/perRow - visibleRows + 1
		}

		for i := 0; i < visibleRows*perRow && offset*perRow+i < len(emojis); i++ {
			idx := offset*perRow + i
			if idx >= len(emojis) {
				break
			}

			row := i / perRow
			col := i % perRow

			emojiStr := fmt.Sprintf(" %s ", emojis[idx].emoji)

			currentStyle := style
			if idx == selected {
				currentStyle = selectedStyle
			}

			x := startX + 1 + col*8
			y := startY + 1 + row

			for j, r := range emojiStr {
				ui.screen.SetContent(x+j, y, r, nil, currentStyle)
			}
		}

		ui.screen.Show()

		ev := ui.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEnter:
				return emojis[selected].emoji, true
			case tcell.KeyEscape:
				return "", false
			case tcell.KeyUp:
				if selected >= perRow {
					selected -= perRow
				}
			case tcell.KeyDown:
				if selected+perRow < len(emojis) {
					selected += perRow
				}
			case tcell.KeyLeft:
				if selected > 0 {
					selected--
				}
			case tcell.KeyRight:
				if selected < len(emojis)-1 {
					selected++
				}
			}
		}
	}
}
