package testrunner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/srmagura/luma/compiler"
	"github.com/srmagura/luma/runtime"
)

// Prefix a test file with _ to make it not run
func TestAll(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err.Error())
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".luma" && filepath.Base(entry.Name())[0] != '_' {
			runTest(t, entry.Name())
		}
	}
}

func runTest(t *testing.T, srcPath string) {
	t.Logf("==== %s ====", filepath.Base(srcPath))

	srcBytes, err := os.ReadFile(srcPath)
	if err != nil {
		t.Fatal("Failed to read the source file.")
	}

	src := string(srcBytes)

	ast, err := compiler.Compile(src)
	if err != nil {
		t.Fatal(err.Error())
	}

	var actualOutputBuilder strings.Builder
	runtime.Execute(ast, &actualOutputBuilder)

	outPath := strings.ReplaceAll(srcPath, ".luma", ".out")
	outBytes, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal("Failed to read the output file.")
	}

	expectedOutput := strings.TrimRight(string(outBytes), "\n")
	actualOutput := strings.TrimRight(actualOutputBuilder.String(), "\n")

	t.Logf("EXPECTED:\n%s\n", expectedOutput)
	t.Logf("ACTUAL:\n%s\n", actualOutput)

	if expectedOutput != actualOutput {
		t.Fatal("Output did not match")
	}
}
