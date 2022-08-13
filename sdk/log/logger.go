package log

import (
	"fmt"
	"github.com/danielcosme/curious-ape/sdk/colors"
	"github.com/danielcosme/curious-ape/sdk/errors"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// Level represents the severity level for a log entry.
type Level int8

const (
	LevelOff   Level = iota // 0
	LevelTrace              // 1 ... etc
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "TRAC"
	case LevelDebug:
		return "DEBU"
	case LevelInfo:
		return "INFO"
	case LevelWarning:
		return "WARN"
	case LevelError:
		return "ERRO"
	case LevelFatal:
		return "FATL"
	default:
		return ""
	}
}

func (l Level) StringColor() string {
	switch l {
	case LevelTrace:
		return colors.Gray(l.String())
	case LevelDebug:
		return colors.White(l.String())
	case LevelInfo:
		return colors.Blue(l.String())
	case LevelWarning:
		return colors.Yellow(l.String())
	case LevelError:
		return colors.Red(l.String())
	case LevelFatal:
		return colors.Red(l.String())
	default:
		return ""
	}
}

type Logger struct {
	out        io.Writer // output destination
	minLevel   Level
	mu         sync.Mutex
	timeFormat string
}

var DefaultLogger = &defaultLogger
var defaultLogger Logger

func init() {
	defaultLogger = Logger{
		out:        os.Stdout,
		minLevel:   0,
		timeFormat: time.RFC3339,
	}
}

type Prop map[string]string

func New(out io.Writer, minLevel Level, timeFormat string) *Logger {
	return &Logger{
		out:        out,
		minLevel:   minLevel,
		timeFormat: timeFormat,
	}
}

func (l *Logger) Trace(args ...any) {
	l.print(LevelTrace, args...)
}

func (l *Logger) TraceP(message string, properties map[string]string) {
	l.printP(LevelTrace, message, properties, nil)
}

func (l *Logger) Tracef(message string, args ...any) {
	l.printf(LevelTrace, message, args...)
}

func (l *Logger) Debug(args ...any) {
	l.print(LevelDebug, args...)
}

func (l *Logger) DebugP(message string, properties map[string]string) {
	l.printP(LevelDebug, message, properties, nil)
}

func (l *Logger) Debugf(message string, args ...any) {
	l.printf(LevelDebug, message, args...)
}

func (l *Logger) Info(args ...any) {
	l.print(LevelInfo, args...)
}

func (l *Logger) InfoP(message string, properties map[string]string) {
	l.printP(LevelInfo, message, properties, nil)
}

func (l *Logger) Infof(message string, args ...any) {
	l.printf(LevelInfo, message, args...)
}

func (l *Logger) Warning(args ...any) {
	l.print(LevelWarning, args...)
}

func (l *Logger) WarningP(message string, properties map[string]string) {
	l.printP(LevelWarning, message, properties, nil)
}

func (l *Logger) Warningf(message string, args ...any) {
	l.printf(LevelWarning, message, args...)
}

func (l *Logger) Error(err error) {
	l.printP(LevelError, err.Error(), nil, getStack(err))
}

func (l *Logger) ErrorP(err error, properties map[string]string) {
	l.printP(LevelError, err.Error(), properties, getStack(err))
}

func (l *Logger) Fatal(err error) {
	l.printP(LevelFatal, err.Error(), nil, getStack(err))
	os.Exit(1) // For entries at the FATAL level, we also terminate the application.
}

func (l *Logger) FatalP(err error, properties map[string]string) {
	l.printP(LevelFatal, err.Error(), properties, getStack(err))
	os.Exit(1) // For entries at the FATAL level, we also terminate the application.
}

func (l *Logger) printf(level Level, message string, args ...any) {
	l.printP(level, fmt.Sprintf(message, args...), nil, nil)
}

func (l *Logger) print(level Level, args ...any) {
	l.printP(level, fmt.Sprint(args...), nil, nil)
}

func (l *Logger) printP(level Level, message string, properties map[string]string, stackTrace []byte) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}
	// TODO don't use colors if sending logs to file or external system
	// TODO Implement JSON formatting
	// options
	// 		- json
	// 		- text with colors
	// 		- text without colors

	aux := fields{
		Level:      level.String(),
		Time:       time.Now().Format(l.timeFormat), // TODO make this in UTC
		Message:    message,
		Properties: properties,
	}
	// TODO send with colors only if the writer is a tty

	keys := make([]string, 0, len(properties))
	for k := range properties {
		keys = append(keys, k)
	}
	// TODO only sort when writing to a tty
	if !sort.StringsAreSorted(keys) {
		sort.Strings(keys)
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("%s", level.StringColor()))
	lines = append(lines, colors.Purple(aux.Time))
	if message != "" {
		lines = append(lines, message)
	}
	for _, k := range keys {
		lines = append(lines, fmt.Sprintf("%s: %s", colors.Green(k), properties[k]))
	}

	if level >= LevelError && len(stackTrace) > 0 {
		aux.Trace = string(stackTrace)
		lines = append(lines, aux.Trace)
	}
	line := strings.Join(lines, colors.Cyan(" - "))

	// TODO - Enable JSON encoding
	// line, err := json.Marshal(aux)
	// if err != nil {
	// 	line = []byte(LevelError.String() + ": unable to marshal log message: " + err.Error())
	// }

	// Lock the mutex so that no two writes to the output destination can happen
	// concurrently. If we don't do this, it's possible that the text for two or more
	// log entries will be intermingled in the output.
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append([]byte(line), '\n'))
}

func (l *Logger) Write(message []byte) (n int, err error) {
	return l.printP(LevelError, string(message), nil, nil)
}

type fields struct {
	Level      string            `json:"level"`
	Time       string            `json:"time"`
	Message    string            `json:"message"`
	Properties map[string]string `json:"properties,omitempty"`
	Trace      string            `json:"trace,omitempty"`
}

func getStack(err error) []byte {
	stack := []byte{}
	if e, ok := err.(*errors.Error); ok {
		stack = e.Stack
	}
	return stack
}
