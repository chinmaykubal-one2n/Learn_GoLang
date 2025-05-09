package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type grepTestCase struct {
	name            string
	searchString    string
	input           string
	wantMatches     []string
	caseInsensitive bool
	before          int
	after           int
	conuntOnly      bool
}

var grepTestCases = []grepTestCase{
	{
		name:         "Zero matches",
		searchString: "hello",
		input:        "this is a test\nanother line\nno match here\n",
		wantMatches:  []string{},
	},
	{
		name:         "One match",
		searchString: "apple",
		input:        "orange\nbanana\napple\n",
		wantMatches:  []string{"apple"},
	},
	{
		name:         "Many matches",
		searchString: "cat",
		input:        "cat is here\ncat again\nno match\nwildcat\n",
		wantMatches:  []string{"cat is here", "cat again", "wildcat"},
	},
	{
		name:            "Case insensitive match",
		searchString:    "hello",
		input:           "line1\nHello World!\nline3\n",
		wantMatches:     []string{"Hello World!"},
		caseInsensitive: true,
	},
	{
		name:         "Before context",
		searchString: "cat",
		input:        "line1\nline2\nline3\ncat is here\nline4\n",
		wantMatches:  []string{"line2", "line3", "cat is here"},
		before:       2,
	},
	{
		name:         "After context",
		searchString: "cat",
		input:        "line1\nline2\nline3\ncat is here\nline4\nline5\n",
		wantMatches:  []string{"cat is here", "line4", "line5"},
		after:        2,
	},
	{
		name:         "Before and after context",
		searchString: "cat",
		input:        "line1\nline2\nline3\ncat is here\nline4\nline5\n",
		wantMatches:  []string{"line2", "line3", "cat is here", "line4", "line5"},
		before:       2,
		after:        2,
	},
	{
		name:         "Count only",
		searchString: "cat",
		input:        "line1\nline2\ncat is here\ncat again\nno match\nwildcat\n",
		wantMatches:  []string{"3"},
		conuntOnly:   true,
	},
}

func TestGrepReader(t *testing.T) {
	for _, grepTestCase := range grepTestCases {
		countOnly = grepTestCase.conuntOnly
		before = grepTestCase.before
		after = grepTestCase.after
		caseInsensitive = grepTestCase.caseInsensitive

		reader := strings.NewReader(grepTestCase.input)
		gotMatches, _ := grepReader(grepTestCase.searchString, reader)

		for i := range gotMatches {
			if gotMatches[i] != grepTestCase.wantMatches[i] {
				t.Errorf("grepReader() mismatch at index %d: got %s, want %s", i, gotMatches[i], grepTestCase.wantMatches[i])
			}
		}
	}
}

func TestWriteToFile(t *testing.T) {
	var lines []string
	var expectedContent string

	for _, grepTestCase := range grepTestCases {
		lines = append(lines, grepTestCase.input)
	}

	for _, line := range lines {
		expectedContent += string(line) + "\n"
	}

	writeToFile("test_output.txt", lines)

	fileContent, err := os.ReadFile("test_output.txt")
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
	}
	if string(fileContent) != expectedContent {
		t.Errorf("File content mismatch: got %s, want %s", string(fileContent), expectedContent)
	}

	// checking for file exsists error
	fileExsistsErr := writeToFile("test_output.txt", lines)
	if fileExsistsErr == nil || !strings.Contains(fileExsistsErr.Error(), "already exists") {
		t.Errorf("Expected file exists error, got: %v", fileExsistsErr)
	}

	err = os.Remove("test_output.txt")
	if err != nil {
		t.Errorf("Failed to remove test file: %v", err)
	}

	// checking for invalid file path
	invalidFilePathErr := writeToFile("/invalid/path/to/file.txt", []string{"data"})
	if invalidFilePathErr == nil {
		t.Errorf("Expected error on invalid path, got %v", invalidFilePathErr)
	}
}

func TestValidateFile(t *testing.T) {
	file, err := os.Create("test.txt")

	if err != nil {
		t.Errorf("Failed to create test file: %v", err)
	}
	defer file.Close()
	defer os.Remove("test.txt")

	// Test valid file
	properFile, err := validateFile(file.Name())
	if err != nil {
		t.Errorf("validateFile() error = %v", err)
	}
	if properFile.Name() != file.Name() {
		t.Errorf("validateFile() got = %v, want %v", properFile.Name(), file.Name())
	}

	// Test non-existent file
	_, nonExistentErr := validateFile("nonexistent.txt")
	if nonExistentErr == nil || !strings.Contains(nonExistentErr.Error(), "open: No such file or directory") {
		t.Errorf("expected 'No such file' error, got %v", nonExistentErr)
	}

	// Test directory
	err = os.Mkdir("testdir", 0755)
	if err != nil {
		t.Errorf("Failed to create test directory: %v", err)
	}
	defer os.Remove("testdir")
	_, dirErr := validateFile("testdir")
	if dirErr == nil || !strings.Contains(dirErr.Error(), "read: Is a directory") {
		t.Errorf("expected directory error, got %v", dirErr)
	}

	// Test permission denied
	privateFile, err := os.Create("private.txt")
	if err != nil {
		t.Errorf("Failed to create private file: %v", err)
	}
	defer os.Chmod(privateFile.Name(), 0600)
	defer os.Remove(privateFile.Name())
	os.Chmod(privateFile.Name(), 000)

	_, permissionErr := validateFile(privateFile.Name())
	if permissionErr == nil || !strings.Contains(permissionErr.Error(), "Permission denied") {
		t.Errorf("expected directory error, got %v", permissionErr)
	}
}

func TestRecursiveSearch(t *testing.T) {
	countOnly = false
	before = 2
	after = 2
	caseInsensitive = true

	tmp := t.TempDir()

	dir1 := filepath.Join(tmp, "dir1")
	dir2 := filepath.Join(tmp, "dir2")
	dir3 := filepath.Join(tmp, "dir3")
	os.MkdirAll(dir1, 0755)
	os.MkdirAll(dir2, 0755)
	os.MkdirAll(dir3, 0755)

	file1 := filepath.Join(dir1, "file.txt")
	os.WriteFile(file1, []byte("Hello from dir1"), 0644)

	file2 := filepath.Join(dir2, "file.txt")
	os.WriteFile(file2, []byte("Hello from dir2"), 0644)

	var buf bytes.Buffer
	recursiveSearch("hello", tmp, &buf)
	got := buf.String()

	want1 := fmt.Sprintf("%s:Hello from dir1\n", file1)
	want2 := fmt.Sprintf("%s:Hello from dir2\n", file2)

	if !strings.Contains(got, want1) {
		t.Errorf("Expected output %s not found.\nGot:\n%s", want1, got)
	}

	if !strings.Contains(got, want2) {
		t.Errorf("Expected output %s not found.\nGot:\n%s", want2, got)
	}
}

func TestWriteStdout(t *testing.T) {
	lines := []string{"Hello", "World"}
	var buf bytes.Buffer

	writeStdout(lines, &buf)

	got := buf.String()
	want := "Hello\nWorld\n"

	if got != want {
		t.Errorf("Expected:\n%s\nGot:\n%s", want, got)
	}

}
