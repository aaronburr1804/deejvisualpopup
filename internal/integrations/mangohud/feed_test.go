package mangohud

import (
	"os"
	"path/filepath"
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
	err := feed.Write("label", 0.5)
	if err == nil {
		t.Fatalf("Write() expected error, got nil")
	}
}
