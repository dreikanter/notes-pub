package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreikanter/npub"
	"github.com/dreikanter/npub/internal/config"
)

func runInit(t *testing.T, args ...string) (string, error) {
	t.Helper()
	rootCmd.SetArgs(append([]string{"init"}, args...))
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
	})
	err := rootCmd.Execute()
	return out.String(), err
}

func TestInitCommandWritesSampleConfig(t *testing.T) {
	dir := t.TempDir()

	out, err := runInit(t, dir)
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	target := filepath.Join(dir, config.DefaultConfigFile)
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !bytes.Equal(got, npub.SampleConfig) {
		t.Errorf("written config does not match embedded sample")
	}
	if !strings.Contains(out, target) {
		t.Errorf("output %q does not mention target path %q", out, target)
	}
}

func TestInitCommandCreatesMissingDirectory(t *testing.T) {
	parent := t.TempDir()
	dir := filepath.Join(parent, "new", "project")

	if _, err := runInit(t, dir); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, config.DefaultConfigFile)); err != nil {
		t.Errorf("expected config file to exist: %v", err)
	}
}

func TestInitCommandRefusesToOverwrite(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, config.DefaultConfigFile)
	if err := os.WriteFile(target, []byte("existing: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := runInit(t, dir); err == nil {
		t.Fatal("expected error when config already exists, got nil")
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if string(got) != "existing: true\n" {
		t.Errorf("existing config was overwritten: got %q", got)
	}
}

func TestInitCommandDefaultsToCwd(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	if _, err := runInit(t); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, config.DefaultConfigFile)); err != nil {
		t.Errorf("expected config file in cwd: %v", err)
	}
}
