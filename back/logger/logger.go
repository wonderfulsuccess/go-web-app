package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorPurple = "\033[35m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

var (
	// DefaultLogger is the default logger instance
	DefaultLogger = New()
)

// Logger represents a logging object
type Logger struct {
	mu     sync.Mutex
	logger *slog.Logger
	debug  bool
}

// New creates a new logger instance
func New() *Logger {
	return &Logger{
		logger: slog.New(newColoredHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
		debug: false,
	}
}

// coloredHandler is a custom handler that adds colors to log output
type coloredHandler struct {
	handler slog.Handler
	output  io.Writer
}

func newColoredHandler(output io.Writer, opts *slog.HandlerOptions) *coloredHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &coloredHandler{
		handler: slog.NewTextHandler(output, opts),
		output:  output,
	}
}

func (h *coloredHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *coloredHandler) Handle(ctx context.Context, r slog.Record) error {
	var color string
	switch r.Level {
	case slog.LevelInfo:
		color = colorGreen
	case slog.LevelWarn:
		color = colorYellow
	case slog.LevelError:
		color = colorRed
	default:
		if r.Level >= slog.LevelError+1 { // Fatal level
			color = colorPurple
		}
	}

	// Format the time
	timeStr := r.Time.Format("02-15:04:05")

	// Format the log message
	msg := fmt.Sprintf("%s[%s]%s%s %s", 
		color,
		r.Level.String(),
		colorReset,
		timeStr,
		r.Message,
	)

	// Add attributes if any
	r.Attrs(func(attr slog.Attr) bool {
		msg += fmt.Sprintf(" %s=%v", attr.Key, attr.Value)
		return true
	})

	// Add newline and write to output
	msg += "\n"
	_, err := fmt.Fprint(h.output, msg)
	return err
}

func (h *coloredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &coloredHandler{
		handler: h.handler.WithAttrs(attrs),
		output:  h.output,
	}
}

func (h *coloredHandler) WithGroup(name string) slog.Handler {
	return &coloredHandler{
		handler: h.handler.WithGroup(name),
		output:  h.output,
	}
}

// SetOutput sets the output destination
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger = slog.New(newColoredHandler(w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

// SetPrefix sets the output prefix
func (l *Logger) SetPrefix(p string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger = l.logger.With("prefix", p)
}

// SetDebug sets the debug mode
func (l *Logger) SetDebug(debug bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.debug = debug
	if debug {
		l.logger = slog.New(newColoredHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}
}

// output prints the log message
func (l *Logger) output(ctx context.Context, level slog.Level, msg string, v ...interface{}) {
	if !l.logger.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	// Skip 3 frames to get the correct caller:
	// 0: runtime.Callers
	// 1: runtime.Callers
	// 2: logger.(*Logger).output
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), level, fmt.Sprintf(msg, v...), pcs[0])
	_ = l.logger.Handler().Handle(ctx, r)
}

// Debugf logs a debug message
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.output(context.Background(), slog.LevelDebug, format, v...)
}

// Infof logs an info message
func (l *Logger) Infof(format string, v ...interface{}) {
	l.output(context.Background(), slog.LevelInfo, format, v...)
}

// Warnf logs a warning message
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.output(context.Background(), slog.LevelWarn, format, v...)
}

// Errorf logs an error message
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.output(context.Background(), slog.LevelError, format, v...)
}

// Fatalf logs a fatal message and exits
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.output(context.Background(), slog.LevelError+1, format, v...)
	os.Exit(1)
}

// Package-level functions

// Debug logs a debug message using the default logger
func Debug(format string, v ...interface{}) {
	DefaultLogger.Debugf(format, v...)
}

// Debugf logs a debug message using the default logger
func Debugf(format string, v ...interface{}) {
	DefaultLogger.Debugf(format, v...)
}

// Debugln logs a debug message with a newline
func Debugln(v ...interface{}) {
	DefaultLogger.Debugf("%s", v...)
}

// Info logs an info message using the default logger
func Info(format string, v ...interface{}) {
	DefaultLogger.Infof(format, v...)
}

// Infof logs an info message using the default logger
func Infof(format string, v ...interface{}) {
	DefaultLogger.Infof(format, v...)
}

// Infoln logs an info message with a newline
func Infoln(v ...interface{}) {
	DefaultLogger.Infof("%s", v...)
}

// Warning logs a warning message using the default logger
func Warning(format string, v ...interface{}) {
	DefaultLogger.Warnf(format, v...)
}

// Warningf logs a warning message using the default logger
func Warningf(format string, v ...interface{}) {
	DefaultLogger.Warnf(format, v...)
}

// Warningln logs a warning message with a newline
func Warningln(v ...interface{}) {
	DefaultLogger.Warnf("%s", v...)
}

// Error logs an error message using the default logger
func Error(format string, v ...interface{}) {
	DefaultLogger.Errorf(format, v...)
}

// Errorf logs an error message using the default logger
func Errorf(format string, v ...interface{}) {
	DefaultLogger.Errorf(format, v...)
}

// Errorln logs an error message with a newline
func Errorln(v ...interface{}) {
	DefaultLogger.Errorf("%s", v...)
}

// Fatal logs a fatal message and exits using the default logger
func Fatal(format string, v ...interface{}) {
	DefaultLogger.Fatalf(format, v...)
}

// Fatalf logs a fatal message and exits using the default logger
func Fatalf(format string, v ...interface{}) {
	DefaultLogger.Fatalf(format, v...)
}

// Fatalln logs a fatal message with a newline and exits
func Fatalln(v ...interface{}) {
	DefaultLogger.Fatalf("%s", v...)
}