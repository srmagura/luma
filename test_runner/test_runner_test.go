package testrunner

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// To execute these tests, run test.sh in the repository root directory

// Prefix a test file with _ to make it not run
func TestAll(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err.Error())
	}

	for _, entry := range entries {
		if !entry.IsDir() &&
			filepath.Ext(entry.Name()) == ".luma" &&
			filepath.Base(entry.Name())[0] != '_' {
			t.Run(
				filepath.Base(entry.Name()),
				func(t *testing.T) {
					runTest(t, entry.Name())
				},
			)
		}
	}
}

func runTest(t *testing.T, srcPath string) {
	compilerOutput, err := exec.Command("../compiler/lumac", srcPath).Output()
	if err != nil {
		t.Fatalf("Compilation failed:\n%s", compilerOutput)
	}

	expectedPath := strings.ReplaceAll(srcPath, ".luma", ".expected")
	expectedBytes, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatal("Failed to read the expected output file.")
	}

	expectedOutput := strings.TrimRight(string(expectedBytes), "\n")
	actualOutput := ""
	// actualOutput := strings.TrimRight(actualOutputBuilder.String(), "\n")

	t.Logf("EXPECTED:\n%s\n", expectedOutput)
	t.Logf("ACTUAL:\n%s\n", actualOutput)

	if expectedOutput != actualOutput {
		t.Fatal("Output did not match")
	}
}
