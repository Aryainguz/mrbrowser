// Package telemetry provides a lightweight, colorful structured logger for Mr. Browser.
// It is intentionally simple — no external infrastructure required.
package telemetry

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// Level represents the log verbosity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorGray    = "\033[90m"
	colorGreen   = "\033[32m"
	colorCyan    = "\033[36m"
	colorYellow  = "\033[33m"
	colorRed     = "\033[31m"
	colorBold    = "\033[1m"
	colorMagenta = "\033[35m"
	colorBlue    = "\033[34m"
)

// Logger is the Mr. Browser structured logger.
type Logger struct {
	mu        sync.Mutex
	out       io.Writer
	level     Level
	component string
	noColor   bool
}

// Field is a key-value pair for structured logging.
type Field struct {
	Key   string
	Value interface{}
}

// F creates a new log field.
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

var defaultLogger = &Logger{
	out:     os.Stdout,
	level:   LevelInfo,
	noColor: os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb",
}

// New creates a named component logger.
func New(component string) *Logger {
	return &Logger{
		out:       os.Stdout,
		level:     defaultLogger.level,
		component: component,
		noColor:   defaultLogger.noColor,
	}
}

// SetLevel sets the global log level.
func SetLevel(l Level) {
	defaultLogger.level = l
}

// SetLevelFromString parses and sets the log level.
func SetLevelFromString(s string) {
	switch strings.ToLower(s) {
	case "debug":
		SetLevel(LevelDebug)
	case "warn", "warning":
		SetLevel(LevelWarn)
	case "error":
		SetLevel(LevelError)
	default:
		SetLevel(LevelInfo)
	}
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, fields ...Field) {
	l.log(LevelDebug, msg, fields...)
}

// Info logs an informational message.
func (l *Logger) Info(msg string, fields ...Field) {
	l.log(LevelInfo, msg, fields...)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, fields ...Field) {
	l.log(LevelWarn, msg, fields...)
}

// Error logs an error message.
func (l *Logger) Error(msg string, fields ...Field) {
	l.log(LevelError, msg, fields...)
}

// Success logs a success message (info-level with green checkmark).
func (l *Logger) Success(msg string, fields ...Field) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level > LevelInfo {
		return
	}
	l.write(colorGreen, "✓", msg, fields...)
}

// Step logs a workflow step (info-level with arrow).
func (l *Logger) Step(msg string, fields ...Field) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level > LevelInfo {
		return
	}
	l.write(colorCyan, "→", msg, fields...)
}

// Recover logs a self-healing recovery event.
func (l *Logger) Recover(msg string, fields ...Field) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level > LevelInfo {
		return
	}
	l.write(colorMagenta, "⟳", msg, fields...)
}

func (l *Logger) log(level Level, msg string, fields ...Field) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	var icon, clr string
	switch level {
	case LevelDebug:
		icon, clr = "·", colorGray
	case LevelInfo:
		icon, clr = "ℹ", colorBlue
	case LevelWarn:
		icon, clr = "⚠", colorYellow
	case LevelError:
		icon, clr = "✗", colorRed
	}

	l.write(clr, icon, msg, fields...)
}

func (l *Logger) write(clr, icon, msg string, fields ...Field) {
	now := time.Now().Format("15:04:05.000")

	var sb strings.Builder

	if l.noColor {
		sb.WriteString(fmt.Sprintf("%s %s ", now, icon))
	} else {
		sb.WriteString(fmt.Sprintf("%s%s%s %s%s%s ", colorGray, now, colorReset, clr, icon, colorReset))
	}

	if l.component != "" {
		if l.noColor {
			sb.WriteString(fmt.Sprintf("[%s] ", l.component))
		} else {
			sb.WriteString(fmt.Sprintf("%s[%s]%s ", colorBold+colorCyan, l.component, colorReset))
		}
	}

	sb.WriteString(msg)

	for _, f := range fields {
		if l.noColor {
			sb.WriteString(fmt.Sprintf("  %s=%v", f.Key, f.Value))
		} else {
			sb.WriteString(fmt.Sprintf("  %s%s%s=%v", colorGray, f.Key, colorReset, f.Value))
		}
	}

	sb.WriteByte('\n')
	_, _ = fmt.Fprint(l.out, sb.String())
}

// Package-level convenience functions using the default logger.

// Debug logs at debug level on the default logger.
func Debug(msg string, fields ...Field) { defaultLogger.log(LevelDebug, msg, fields...) }

// Info logs at info level on the default logger.
func Info(msg string, fields ...Field) { defaultLogger.log(LevelInfo, msg, fields...) }

// Warn logs at warn level on the default logger.
func Warn(msg string, fields ...Field) { defaultLogger.log(LevelWarn, msg, fields...) }

// Error logs at error level on the default logger.
func Error(msg string, fields ...Field) { defaultLogger.log(LevelError, msg, fields...) }
