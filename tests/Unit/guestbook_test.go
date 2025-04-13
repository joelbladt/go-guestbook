package guestbook_test

import (
	"os"
	"strings"
	"testing"

	"github.com/joelbladt/go-guestbook/src/guestbook"
)

// setupTestFile creates a temporary guestbook file for isolated tests.
func setupTestFile(content string) (string, func()) {
	tmpfile, err := os.CreateTemp("", "guestbook_test_*.txt")
	if err != nil {
		panic(err)
	}
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		panic(err)
	}
	tmpfile.Close()

	// Override global guestbook file path
	oldFile := guestbook.GuestbookFile
	guestbook.GuestbookFile = tmpfile.Name()

	cleanup := func() {
		os.Remove(tmpfile.Name())
		guestbook.GuestbookFile = oldFile
	}

	return tmpfile.Name(), cleanup
}

// TestSaveAndShow tests both Save and Show with a full cycle.
func TestSaveAndShow(t *testing.T) {
	_, cleanup := setupTestFile("")
	defer cleanup()

	err := guestbook.Save("TestUser", "Hello **World**!\nNew line 😊")
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	entries, err := guestbook.Show()
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}

	entry := entries[0]

	if entry.Name != "TestUser" {
		t.Errorf("Expected name 'TestUser', got '%s'", entry.Name)
	}

	if !strings.Contains(string(entry.Message), "<strong>World</strong>") {
		t.Errorf("Expected markdown to be rendered, got: %s", entry.Message)
	}

	if !strings.Contains(string(entry.Message), "<br>") {
		t.Errorf("Expected line break to be converted, got: %s", entry.Message)
	}
}

// TestEmptyNameOrMessage ensures empty name or message results in no entry.
func TestEmptyNameOrMessage(t *testing.T) {
	_, cleanup := setupTestFile("")
	defer cleanup()

	err := guestbook.Save("", "Test message")
	if err != nil {
		t.Fatalf("Save failed unexpectedly: %v", err)
	}

	err = guestbook.Save("TestUser", "")
	if err != nil {
		t.Fatalf("Save failed unexpectedly: %v", err)
	}

	entries, err := guestbook.Show()
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if len(entries) != 0 {
		t.Fatalf("Expected 0 entries, got %d", len(entries))
	}
}

// TestRenderMarkdown verifies the renderMarkdown function independently.
func TestRenderMarkdown(t *testing.T) {
	raw := "Hello **World**!\nLine 2"
	html := guestbook.RenderMarkdown(raw)

	if !strings.Contains(string(html), "<strong>World</strong>") {
		t.Errorf("Expected bold tag, got: %s", html)
	}

	if !strings.Contains(string(html), "<br>") {
		t.Errorf("Expected line break, got: %s", html)
	}
}

// TestShowWithCorruptData checks that corrupt entries are ignored
func TestShowWithCorruptData(t *testing.T) {
	badContent := `
---ENTRY---
invalid_timestamp
Someone
Message

---ENTRY---
2025-04-13T22:22:00+02:00
Joel
Valid entry!
`
	_, cleanup := setupTestFile(badContent)
	defer cleanup()

	entries, err := guestbook.Show()
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 valid entry, got %d", len(entries))
	}

	if entries[0].Name != "Joel" {
		t.Errorf("Expected entry by Joel, got: %s", entries[0].Name)
	}
}

// TestAddLineBreaks checks custom linebreak replacement inside <p>
func TestAddLineBreaks(t *testing.T) {
	input := "<p>Hello\nWorld</p>" // <-- real row
	expected := "<p>Hello<br>\nWorld</p>"
	result := guestbook.AddLineBreaks(input)
	if result != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, result)
	}
}
