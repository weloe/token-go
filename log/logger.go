package log

import "github.com/weloe/token-go/persist"

type Logger interface {
	persist.Watcher

	// Enable turn on or off
	Enable(bool bool)

	// IsEnabled return if logger is enabled
	IsEnabled() bool

	// StartCleanTimer log after start clean timer
	StartCleanTimer(period int64)
}
