package log

import (
	"time"
)

const (
	defaultLogFormat       = "%emoji%  %msg%"
	defaultTimestampFormat = time.RFC3339
)

// PwcCtlLogFormat is a custom logrus formatter.
type PwcCtlLogFormat struct {
	TimestampFormat string
	LogFormat       string
	Color           bool
}
