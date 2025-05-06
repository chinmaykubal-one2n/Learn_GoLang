package cmd

import (
	"os"
	"strings"
	"testing"
)

type grepTestCase struct {
	name         string
	searchString string
	input        string
	wantMatches  []string
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
}

func TestGrepReader(t *testing.T) {
	for _, tt := range grepTestCases {
		reader := strings.NewReader(tt.input)
		gotMatches, err := grepReader(tt.searchString, reader)

		if err != nil {
			t.Errorf("grepReader() error = %v", err)
		}

		for i := range gotMatches {
			if gotMatches[i] != tt.wantMatches[i] {
				t.Errorf("grepReader() mismatch at index %d: got %s, want %s", i, gotMatches[i], tt.wantMatches[i])
			}
		}
	}
}

// do we even need the below test?
// func TestWriteStdout(t *testing.T) {
// 	var lines []string
// 	for _, grepTestCase := range grepTestCases {
// 		lines = append(lines, grepTestCase.input)
// 	}
// 	writeStdout(lines)
// }

func TestWriteToFile(t *testing.T) {
	var lines []string
	var expectedContent string

	for _, grepTestCase := range grepTestCases {
		lines = append(lines, grepTestCase.input)
	}

	for _, line := range lines {
		expectedContent += string(line) + "\n"
	}

	err := writeToFile("test_output.txt", lines)
	if err != nil {
		t.Errorf("writeToFile() error = %v", err)
	}

	fileContent, err := os.ReadFile("test_output.txt")
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
	}
	if string(fileContent) != expectedContent {
		t.Errorf("File content mismatch: got %s, want %s", string(fileContent), expectedContent)
	}

	err = os.Remove("test_output.txt")
	if err != nil {
		t.Errorf("Failed to remove test file: %v", err)
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
