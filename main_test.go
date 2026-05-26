package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDelExactMatch(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := run([]string{"add", "milk"}); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"add", "soy milk"}); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"del", "milk"}); err != nil {
		t.Fatal(err)
	}

	got := readRegisterFile(t, home, "default")
	if got != "soy milk\n" {
		t.Fatalf("register = %q, want %q", got, "soy milk\n")
	}
}

func TestDelRegex(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := run([]string{"add", "milk"}); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"add", "soy milk"}); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"del", "-r", "milk"}); err != nil {
		t.Fatal(err)
	}

	got := readRegisterFile(t, home, "default")
	if got != "" {
		t.Fatalf("register = %q, want empty", got)
	}
}

func TestDelDryRun(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := run([]string{"add", "milk"}); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"add", "soy milk"}); err != nil {
		t.Fatal(err)
	}

	out := captureStdout(t, func() {
		if err := run([]string{"del", "-d", "-r", "milk"}); err != nil {
			t.Fatal(err)
		}
	})

	if out != "milk\nsoy milk\n" {
		t.Fatalf("dry run output = %q, want %q", out, "milk\nsoy milk\n")
	}
	got := readRegisterFile(t, home, "default")
	if got != "milk\nsoy milk\n" {
		t.Fatalf("register = %q, want unchanged", got)
	}
}

func readRegisterFile(t *testing.T, home, name string) string {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(home, registerDir, name+fileExt))
	if os.IsNotExist(err) {
		return ""
	}
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	fn()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	return strings.ReplaceAll(buf.String(), "\r\n", "\n")
}
