package main

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"testing"
)

// expectedStdout is the exact byte sequence a successful run must write: the
// pinned greeting followed by a single newline. It is a literal on purpose — a
// comparison against the message constant would still pass after the greeting
// was mutated, which would defeat this test.
const expectedStdout = "Hello, World!\n"

// TestMainPrintsHelloWorld builds the program once into a temporary directory
// and runs the resulting binary from a working directory other than the module
// root, asserting the full output contract: exit status 0, empty stderr and
// exactly the pinned bytes on stdout.
func TestMainPrintsHelloWorld(t *testing.T) {
	goTool, err := exec.LookPath("go")
	if err != nil {
		t.Skipf("go toolchain not found in PATH, cannot build and run the program: %v", err)
	}

	bin := filepath.Join(t.TempDir(), "hello")
	build := exec.Command(goTool, "build", "-o", bin, ".")
	var buildStderr bytes.Buffer
	build.Stderr = &buildStderr
	if err := build.Run(); err != nil {
		t.Fatalf("go build failed: %v\nstderr:\n%s", err, buildStderr.String())
	}

	cmd := exec.Command(bin)
	cmd.Dir = t.TempDir() // exercise working-directory independence
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()
	if cmd.ProcessState == nil {
		t.Fatalf("could not start the built binary %q: %v", bin, runErr)
	}
	if code := cmd.ProcessState.ExitCode(); code != 0 {
		t.Fatalf("exit status = %d, want 0\nstderr:\n%s", code, stderr.String())
	}
	if got := stderr.String(); got != "" {
		t.Fatalf("stderr = %q, want empty", got)
	}
	if got := stdout.String(); got != expectedStdout {
		t.Fatalf("stdout = %q, want %q", got, expectedStdout)
	}
}
