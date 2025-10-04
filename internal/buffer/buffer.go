package buffer

import (
	"bufio"
	"os"
	"strings"
)

type Buffer struct {
	Lines      []string
	FilePath   string
	Modified   bool
	CursorX    int
	CursorY    int
	SelectMode bool
	SelectX    int
	SelectY    int
}

func New(filePath string) (*Buffer, error) {
	b := &Buffer{
		Lines:    []string{""},
		FilePath: filePath,
		Modified: false,
		CursorX:  0,
		CursorY:  0,
	}

	if filePath != "" {
		if err := b.Load(); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}

	return b, nil
}

func (b *Buffer) Load() error {
	file, err := os.Open(b.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // New file
		}
		return err
	}
	defer file.Close()

	b.Lines = []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		b.Lines = append(b.Lines, scanner.Text())
	}

	if len(b.Lines) == 0 {
		b.Lines = []string{""}
	}

	return scanner.Err()
}

func (b *Buffer) Save() error {
	file, err := os.Create(b.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i, line := range b.Lines {
		if _, err := writer.WriteString(line); err != nil {
			return err
		}
		if i < len(b.Lines)-1 {
			if _, err := writer.WriteString("\n"); err != nil {
				return err
			}
		}
	}

	b.Modified = false
	return writer.Flush()
}

func (b *Buffer) InsertRune(r rune) {
	if b.CursorY >= len(b.Lines) {
		b.Lines = append(b.Lines, "")
		b.CursorY = len(b.Lines) - 1
	}

	line := b.Lines[b.CursorY]
	if b.CursorX > len(line) {
		b.CursorX = len(line)
	}

	b.Lines[b.CursorY] = line[:b.CursorX] + string(r) + line[b.CursorX:]
	b.CursorX++
	b.Modified = true
}

func (b *Buffer) OverwriteRune(r rune) {
	if b.CursorY >= len(b.Lines) {
		b.Lines = append(b.Lines, "")
		b.CursorY = len(b.Lines) - 1
	}

	line := b.Lines[b.CursorY]
	if b.CursorX > len(line) {
		b.CursorX = len(line)
	}

	// If at end of line, just insert
	if b.CursorX >= len(line) {
		b.Lines[b.CursorY] = line + string(r)
	} else {
		// Replace character at cursor position
		b.Lines[b.CursorY] = line[:b.CursorX] + string(r) + line[b.CursorX+1:]
	}
	b.CursorX++
	b.Modified = true
}

func (b *Buffer) InsertNewline() {
	if b.CursorY >= len(b.Lines) {
		b.Lines = append(b.Lines, "")
		b.CursorY = len(b.Lines) - 1
	}

	line := b.Lines[b.CursorY]
	if b.CursorX > len(line) {
		b.CursorX = len(line)
	}

	// Split current line
	before := line[:b.CursorX]
	after := line[b.CursorX:]

	b.Lines[b.CursorY] = before
	b.Lines = append(b.Lines[:b.CursorY+1], append([]string{after}, b.Lines[b.CursorY+1:]...)...)

	b.CursorY++
	b.CursorX = 0
	b.Modified = true
}

func (b *Buffer) DeleteRune() {
	if b.CursorY >= len(b.Lines) {
		return
	}

	line := b.Lines[b.CursorY]

	if b.CursorX > 0 && b.CursorX <= len(line) {
		// Delete character before cursor
		b.Lines[b.CursorY] = line[:b.CursorX-1] + line[b.CursorX:]
		b.CursorX--
		b.Modified = true
	} else if b.CursorX == 0 && b.CursorY > 0 {
		// Join with previous line
		prevLine := b.Lines[b.CursorY-1]
		b.Lines[b.CursorY-1] = prevLine + line
		b.Lines = append(b.Lines[:b.CursorY], b.Lines[b.CursorY+1:]...)
		b.CursorY--
		b.CursorX = len(prevLine)
		b.Modified = true
	}
}

func (b *Buffer) DeleteCurrentLine() string {
	if b.CursorY >= len(b.Lines) {
		return ""
	}

	deleted := b.Lines[b.CursorY]

	if len(b.Lines) == 1 {
		b.Lines[0] = ""
		b.CursorX = 0
	} else {
		b.Lines = append(b.Lines[:b.CursorY], b.Lines[b.CursorY+1:]...)
		if b.CursorY >= len(b.Lines) {
			b.CursorY = len(b.Lines) - 1
		}
		b.CursorX = 0
	}

	b.Modified = true
	return deleted
}

func (b *Buffer) GetCurrentLine() string {
	if b.CursorY < len(b.Lines) {
		return b.Lines[b.CursorY]
	}
	return ""
}

func (b *Buffer) InsertText(text string) {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if i > 0 {
			b.InsertNewline()
		}
		for _, r := range line {
			b.InsertRune(r)
		}
	}
}

func (b *Buffer) ReplaceCurrentLine(text string) {
	if b.CursorY < len(b.Lines) {
		b.Lines[b.CursorY] = text
		b.CursorX = len(text)
		b.Modified = true
	}
}

func (b *Buffer) GetAllText() string {
	return strings.Join(b.Lines, "\n")
}

func (b *Buffer) ToggleSelection() {
	if b.SelectMode {
		b.SelectMode = false
	} else {
		b.SelectMode = true
		b.SelectX = b.CursorX
		b.SelectY = b.CursorY
	}
}

func (b *Buffer) HasSelection() bool {
	return b.SelectMode && (b.SelectX != b.CursorX || b.SelectY != b.CursorY)
}

func (b *Buffer) GetSelection() string {
	if !b.HasSelection() {
		return ""
	}

	startY, endY := b.SelectY, b.CursorY
	startX, endX := b.SelectX, b.CursorX

	if startY > endY || (startY == endY && startX > endX) {
		startY, endY = endY, startY
		startX, endX = endX, startX
	}

	if startY == endY {
		return b.Lines[startY][startX:endX]
	}

	var result strings.Builder
	for y := startY; y <= endY; y++ {
		if y == startY {
			result.WriteString(b.Lines[y][startX:])
		} else if y == endY {
			result.WriteString(b.Lines[y][:endX])
		} else {
			result.WriteString(b.Lines[y])
		}
		if y < endY {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func (b *Buffer) ReplaceSelection(text string) {
	if !b.HasSelection() {
		return
	}

	startY, endY := b.SelectY, b.CursorY
	startX, endX := b.SelectX, b.CursorX

	if startY > endY || (startY == endY && startX > endX) {
		startY, endY = endY, startY
		startX, endX = endX, startX
	}

	newLines := strings.Split(text, "\n")

	if startY == endY {
		b.Lines[startY] = b.Lines[startY][:startX] + text + b.Lines[startY][endX:]
		b.CursorY = startY
		if len(newLines) == 1 {
			b.CursorX = startX + len(text)
		} else {
			b.CursorY = startY + len(newLines) - 1
			b.CursorX = len(newLines[len(newLines)-1])
		}
	} else {
		before := b.Lines[startY][:startX]
		after := b.Lines[endY][endX:]

		var newLinesFull []string
		newLinesFull = append(newLinesFull, b.Lines[:startY]...)

		if len(newLines) == 1 {
			newLinesFull = append(newLinesFull, before+text+after)
		} else {
			newLinesFull = append(newLinesFull, before+newLines[0])
			if len(newLines) > 2 {
				newLinesFull = append(newLinesFull, newLines[1:len(newLines)-1]...)
			}
			newLinesFull = append(newLinesFull, newLines[len(newLines)-1]+after)
		}

		newLinesFull = append(newLinesFull, b.Lines[endY+1:]...)
		b.Lines = newLinesFull

		b.CursorY = startY + len(newLines) - 1
		b.CursorX = len(newLines[len(newLines)-1])
	}

	b.SelectMode = false
	b.Modified = true
}
