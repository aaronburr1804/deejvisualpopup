package app

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/nik9play/deej/internal/audio"
	"github.com/nik9play/deej/internal/config"
	"github.com/nik9play/deej/internal/integrations/mangohud"
	"github.com/nik9play/deej/internal/serial"
	"go.uber.org/zap"
)

// hudWriteDelay is the minimum time between feed writes during a slider sweep
const hudWriteDelay = 40 * time.Millisecond

// hudFeed is a MangoHud feed that can be used to display information on the screen.
type hudFeed struct {
	feed    *mangohud.Feed
	done    chan struct{}
	workers sync.WaitGroup
}

// newHUDFeed returns a new HUD feed or nil if the feed cannot be created.
func newHUDFeed(logger *zap.SugaredLogger) *hudFeed {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		logger.Warnw("MangoHud feed disabled: failed to get user cache directory", "error", err)
		return nil
	}
	path := filepath.Join(cacheDir, "deej", "deej-hud.txt")
	return &hudFeed{
		feed: mangohud.New(path),
		done: make(chan struct{}),
	}
}

// start subscribes to slider moves and writes the latest real move to the feed, at most once per hudWriteDelay
func (h *hudFeed) start(sliders *serial.Client, settings func() *config.Settings, logger *zap.SugaredLogger) {
	ch := sliders.SubscribeToSliderMoveEvents()

	h.workers.Go(func() {
		seen := make(sliderMemory)
		var pendingLabel string
		var pendingValue float64
		var flush <-chan time.Time

		for {
			select {
			case event := <-ch:
				if !seen.moved(event.SliderID, event.PercentValue) {
					continue
				}
				s := settings()
				if !s.MangoHud.Enabled {
					continue
				}
				targets, ok := s.ActiveMapping().Get(event.SliderID)
				if !ok || len(targets) == 0 {
					continue
				}
				pendingLabel = hudLabel(targets, s.MangoHud.UnmappedRename)
				pendingValue = float64(event.PercentValue)
				if flush == nil {
					flush = time.After(hudWriteDelay)
				}
			case <-flush:
				flush = nil
				if err := h.feed.Write(pendingLabel, pendingValue); err != nil {
					logger.Debugw("MangoHud feed write failed", "error", err)
				}
			case <-h.done:
				return
			}
		}
	})
}

// stop stops the feed and waits for all workers to finish.
func (h *hudFeed) stop() {
	close(h.done)
	h.workers.Wait()
}

// targetLabel turns yaml config target into its HUD label.
func targetLabel(target, unmappedRename string) string {
	switch {
	case strings.EqualFold(target, audio.MasterSessionName):
		return "Master"
	case strings.EqualFold(target, audio.InputSessionName):
		return "Mic"
	case strings.EqualFold(target, "deej.unmapped") && unmappedRename != "":
		return unmappedRename
	case strings.HasSuffix(strings.ToLower(target), ".exe"):
		return target[:len(target)-4]
	default:
		return target
	}
}

// hudLabel joins the labels for all of a slider's targets.
func hudLabel(targets []string, unmappedRename string) string {
	labels := make([]string, 0, len(targets))
	for _, target := range targets {
		labels = append(labels, targetLabel(target, unmappedRename))
	}
	return strings.Join(labels, " + ")
}

// sliderMemory remembers each slider's last value, to tell real moves from replays.
type sliderMemory map[int]float32

// moved records v and reports whether it's a real change from a known previous value
func (m sliderMemory) moved(id int, v float32) bool {
	prev, seen := m[id]
	m[id] = v
	return seen && prev != v
}
