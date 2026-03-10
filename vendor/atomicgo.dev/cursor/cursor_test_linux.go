package cursor

import (
	"log"
	"os"
	"testing"
)

// TestCustomIOWriter tests the cursor functions with a custom Writer.
func TestCustomIOWriter(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testingTmpFile-")
	defer os.Remove(tmpFile.Name())

	if err != nil {
		log.Fatal(err)
	}

	w := tmpFile
	SetTarget(w)

	Up(2)

	expected := "\x1b[2A"
	actual := getFileContent(t, w.Name())

	if expected != actual {
		t.Errorf("wanted: %v, got %v", expected, actual)
	}

	clearFile(t, w)
	Down(2)

	expected = "\x1b[2B"
	actual = getFileContent(t, w.Name())

	if expected != actual {
		t.Errorf("wanted: %v, got %v", expected, actual)
	}

	clearFile(t, w)
	Right(2)

	expected = "\x1b[2C"
	actual = getFileContent(t, w.Name())

	if expected != actual {
		t.Errorf("wanted: %v, got %v", expected, actual)
	}

	clearFile(t, w)
	Left(2)

	expected = "\x1b[2D"
	actual = getFileContent(t, w.Name())

	if expected != actual {
		t.Errorf("wanted: %v, got %v", expected, actual)
	}

	clearFile(t, w)
	Hide()

	expected = "\x1b[?25l"
	actual = getFileContent(t, w.Name())

	if expected != actual {
		t.Errorf("wanted: %v, got %v", expected, actual)
	}

	clearFile(t, w)
	Show()

	expected = "\x1b[?25h"
	actual = getFileContent(t, w.Name())

	if expected != actual {
		t.Errorf("wanted: %v, got %v", expected, actual)
	}

	clearFile(t, w)
	ClearLine()

	expected = "\x1b[2K"
	actual = getFileContent(t, w.Name())

	if expected != actual {
		t.Errorf("wanted: %v, got %v", expected, actual)
	}

	clearFile(t, w)
	HorizontalAbsolute(3)

	expected = "\x1b[4G"
	actual = getFileContent(t, w.Name())

	if expected != actual {
		t.Errorf("wanted: %v, got %v", expected, actual)
	}
}

func getFileContent(t *testing.T, fileName string) string {
	t.Helper()

	content, err := os.ReadFile(fileName)
	if err != nil {
		t.Errorf("failed to read file contents: %s", err)

		return ""
	}

	return string(content)
}

func clearFile(t *testing.T, file *os.File) {
	t.Helper()

	err := file.Truncate(0)
	if err != nil {
		t.Errorf("failed to clear file")

		return
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		t.Errorf("failed to clear file")

		return
	}
}
