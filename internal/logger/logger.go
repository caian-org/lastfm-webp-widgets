package logger

import (
	// internal
	. "github.com/caian-org/lastfm-webp-widgets/internal/constants"
	"github.com/caian-org/lastfm-webp-widgets/pkg/log15-2.16.0"
)

var Log = log15.New("cmd", CMD_NAME)
