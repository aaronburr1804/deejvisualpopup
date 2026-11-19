package mangohud

import (
	"fmt"
	"github.com/nik9play/deej/internal/atomicfile"
	"math"
)

type Feed struct {
	path string
}
// Write writes a label and a setting to the feed file. The setting is clamped to [0, 1] and converted to a percentage.	
func (f *Feed) Write(label string, setting float64) error {
	setting = math.Max(0, math.Min(1, setting))
	percent := math.Round(setting * 100)
	return atomicfile.Write(f.path, []byte(fmt.Sprintf("%s: %.0f%%\n", label, percent)))
}
// New creates a new Feed that writes to the given path. The path should be the same as the one used by the HUD script.
func New(path string) *Feed {
	return &Feed{
		path: path,
	}
}
