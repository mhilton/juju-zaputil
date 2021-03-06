package zaputil

import (
	"github.com/juju/loggo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggoToZap = map[loggo.Level]zapcore.Level{
	loggo.TRACE:    zap.DebugLevel, // There's no zap equivalent to TRACE.
	loggo.DEBUG:    zap.DebugLevel,
	loggo.INFO:     zap.InfoLevel,
	loggo.WARNING:  zap.WarnLevel,
	loggo.ERROR:    zap.ErrorLevel,
	loggo.CRITICAL: zap.ErrorLevel, // There's no zap equivalent to CRITICAL.
}

// NewLoggoWriter returns a loggo.Writer that writes to the
// given zap logger.
func NewLoggoWriter(logger *zap.Logger) loggo.Writer {
	return zapLoggoWriter{
		logger: logger,
	}
}

// zapLoggoWriter implements a loggo.Writer by writing to a zap.Logger,
// so can be used as an adaptor from loggo to zap.
type zapLoggoWriter struct {
	logger *zap.Logger
}

// zapLoggoWriter implements loggo.Writer.Write by writing the entry
// to w.logger. It ignores entry.Timestamp because zap will affix its
// own timestamp.
func (w zapLoggoWriter) Write(entry loggo.Entry) {
	if ce := w.logger.Check(loggoToZap[entry.Level], entry.Message); ce != nil {
		ce.Write(zap.String("module", entry.Module), zap.String("file", entry.Filename), zap.Int("line", entry.Line))
	}
}
