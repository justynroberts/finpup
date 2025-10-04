package buffer

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	b, err := New("")
	if err != nil {
		t.Fatalf("Failed to create buffer: %v", err)
	}

	if len(b.Lines) != 1 || b.Lines[0] != "" {
		t.Errorf("Expected empty buffer with one empty line, got %v", b.Lines)
	}
}

func TestInsertRune(t *testing.T) {
	b, _ := New("")
	b.InsertRune('H')
	b.InsertRune('i')

	if b.Lines[0] != "Hi" {
		t.Errorf("Expected 'Hi', got '%s'", b.Lines[0])
	}

	if !b.Modified {
		t.Error("Buffer should be marked as modified")
	}
}

func TestInsertNewline(t *testing.T) {
	b, _ := New("")
	b.InsertRune('H')
	b.InsertRune('i')
	b.InsertNewline()
	b.InsertRune('B')
	b.InsertRune('y')
	b.InsertRune('e')

	if len(b.Lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(b.Lines))
	}

	if b.Lines[0] != "Hi" || b.Lines[1] != "Bye" {
		t.Errorf("Expected ['Hi', 'Bye'], got %v", b.Lines)
	}
}

func TestDeleteRune(t *testing.T) {
	b, _ := New("")
	b.InsertRune('H')
	b.InsertRune('i')
	b.DeleteRune()

	if b.Lines[0] != "H" {
		t.Errorf("Expected 'H', got '%s'", b.Lines[0])
	}
}

func TestDeleteCurrentLine(t *testing.T) {
	b, _ := New("")
	b.InsertRune('H')
	b.InsertRune('i')
	b.InsertNewline()
	b.InsertRune('B')
	b.InsertRune('y')
	b.InsertRune('e')

	deleted := b.DeleteCurrentLine()
	if deleted != "Bye" {
		t.Errorf("Expected deleted line 'Bye', got '%s'", deleted)
	}

	if len(b.Lines) != 1 {
		t.Errorf("Expected 1 line remaining, got %d", len(b.Lines))
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpFile := "/tmp/finton_test.txt"
	defer os.Remove(tmpFile)

	// Create and save
	b1, _ := New(tmpFile)
	b1.InsertRune('H')
	b1.InsertRune('i')
	b1.InsertNewline()
	b1.InsertRune('B')
	b1.InsertRune('y')
	b1.InsertRune('e')

	if err := b1.Save(); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Load and verify
	b2, err := New(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}

	if len(b2.Lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(b2.Lines))
	}

	if b2.Lines[0] != "Hi" || b2.Lines[1] != "Bye" {
		t.Errorf("Expected ['Hi', 'Bye'], got %v", b2.Lines)
	}
}

func TestInsertText(t *testing.T) {
	b, _ := New("")
	b.InsertText("Hello\nWorld")

	if len(b.Lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(b.Lines))
	}

	if b.Lines[0] != "Hello" || b.Lines[1] != "World" {
		t.Errorf("Expected ['Hello', 'World'], got %v", b.Lines)
	}
}

func TestReplaceCurrentLine(t *testing.T) {
	b, _ := New("")
	b.InsertRune('O')
	b.InsertRune('l')
	b.InsertRune('d')
	b.ReplaceCurrentLine("New")

	if b.Lines[0] != "New" {
		t.Errorf("Expected 'New', got '%s'", b.Lines[0])
	}
}
