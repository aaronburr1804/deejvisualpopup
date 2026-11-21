package mangohud

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testCase(name string, input float64, want string) struct {
	name  string
	input float64
	want  string
} {
	return struct {
		name  string
		input float64
		want  string
	}{
		name:  name,
		input: input,
		want:  want,
	}
}
func TestFeed_Write(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  string
	}{
		testCase("clamp below 0", -1, "label: 0%\n"),
		testCase("minimum", 0, "label: 0%\n"),
		testCase("maximum", 1, "label: 100%\n"),
		testCase("middle", 0.5, "label: 50%\n"),
		testCase("clamp above 1", 2, "label: 100%\n"),
		testCase("near tie", 0.605, "label: 61%\n"),
		testCase("round down", 0.123, "label: 12%\n"),
		testCase("round up", 0.567, "label: 57%\n"),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feed := New(filepath.Join(t.TempDir(), "feed.txt"))
			if err := feed.Write("label", tt.input); err != nil {
				t.Fatalf("Write() error = %v", err)
			}
			got, err := os.ReadFile(feed.path)
			if err != nil {
				t.Fatalf("ReadFile() error = %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("ReadFile() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFeed_Write_MissingDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "feed.txt")
	feed := New(path)

	if err := feed.Write("label", 0.5); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(got) != "label: 50%\n" {
		t.Errorf("ReadFile() = %q, want %q", got, "label: 50%\n")
	}
}

func TestFeed_Write_DirCreateFails(t *testing.T) {
	dir := t.TempDir()
	blocker := filepath.Join(dir, "blocker")

	if err := os.WriteFile(blocker, []byte("blocker"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	path := filepath.Join(blocker, "feed.txt")
	feed := New(path)
	err := feed.Write("label", 0.5)

	if err == nil {
		t.Fatalf("Write() expected error, got nil")
	}

	if !strings.Contains(err.Error(), "create directory") {
		t.Errorf("Write() error = %v, does not contain 'create directory'", err)
	}
}
