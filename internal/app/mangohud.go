package app

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/nik9play/deej/internal/config"
	"github.com/nik9play/deej/internal/integrations/mangohud"
	"github.com/nik9play/deej/internal/serial"
	"go.uber.org/zap"
)

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

// start subscribes to slider moves and writes each one to the feed.
func (h *hudFeed) start(sliders *serial.Client, settings func() *config.Settings, logger *zap.SugaredLogger) {
	ch := sliders.SubscribeToSliderMoveEvents()

	h.workers.Go(func() {
		for {
			select {
			case event := <-ch:
				s := settings()
				if !s.MangoHud.Enabled {
					continue
				}
				targets, ok := s.ActiveMapping().Get(event.SliderID)
				if !ok || len(targets) == 0 {
					continue
				}
				label := targets[0]
				if err := h.feed.Write(label, float64(event.PercentValue)); err != nil {
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
