package config

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
)

// Watcher watches configuration file changes
type Watcher struct {
	configPath string
	watcher    *fsnotify.Watcher
	reloadChan chan bool
	stopChan   chan bool
	logger     zerolog.Logger
	mu         sync.Mutex
	lastReload time.Time
	debounce   time.Duration
	stopped    bool
	stopOnce   sync.Once
}

// NewWatcher creates a new configuration file watcher
func NewWatcher(configPath string, logger zerolog.Logger) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	// Watch the directory instead of the file to handle editor rewrites
	dir := filepath.Dir(absPath)
	if err := watcher.Add(dir); err != nil {
		watcher.Close()
		return nil, err
	}

	w := &Watcher{
		configPath: absPath,
		watcher:    watcher,
		reloadChan: make(chan bool, 1),
		stopChan:   make(chan bool),
		logger:     logger,
		debounce:   1 * time.Second, // Debounce to avoid multiple reloads
	}

	return w, nil
}

// Start starts watching for configuration file changes
func (w *Watcher) Start() {
	go w.watch()
}

// watch monitors file system events
func (w *Watcher) watch() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			// Only process events for our config file
			if filepath.Clean(event.Name) != w.configPath {
				continue
			}

			// Process write and create events
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				w.handleChange()
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			w.logger.Error().Err(err).Msg("File watcher error")

		case <-w.stopChan:
			return
		}
	}
}

// handleChange handles configuration file changes with debouncing
func (w *Watcher) handleChange() {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Debounce: ignore changes within debounce duration
	now := time.Now()
	if now.Sub(w.lastReload) < w.debounce {
		w.logger.Debug().Msg("Config change ignored (debounce)")
		return
	}

	w.lastReload = now
	w.logger.Info().Str("file", w.configPath).Msg("Config file changed")

	// Send reload signal (non-blocking)
	select {
	case w.reloadChan <- true:
	default:
		w.logger.Debug().Msg("Reload already pending")
	}
}

// ReloadChan returns the channel that receives reload signals
func (w *Watcher) ReloadChan() <-chan bool {
	return w.reloadChan
}

// Stop stops the watcher
func (w *Watcher) Stop() error {
	var err error
	w.stopOnce.Do(func() {
		w.mu.Lock()
		w.stopped = true
		w.mu.Unlock()
		
		close(w.stopChan)
		err = w.watcher.Close()
		w.logger.Debug().Msg("Config watcher stopped")
	})
	return err
}
