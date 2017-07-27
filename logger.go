/*
Package alog provides asynchronous structured leveled log functionality.

The alog package provides structured logging with JSON output and with
semi-structured text-based logging.  Logging is asynchronous, where work and
time to format log messages and write them to I/O is done by an asynchronous
goroutine, minimizing the time the application waits for a log call to return.
Decisions to log in a normal, especially a performance critical, code path
should still be carefully considered.

A Logger instance generates lines of output to an io.Writer.  Each logging
operation makes a single call to the Writer's Write method.  A Logger can be
used simultaneously from multiple goroutines; it guarantees serialized access
to the Writer.

The use of structured and leveled logging is completely optional.  Leveled
logging is disabled by specifying `alog.NoLevel` when creating a new logger
instance.  If considering using leveled logging, please read this discussion:

https://dave.cheney.net/2015/11/05/lets-talk-about-logging

Using alog requires creating a logger instance.  There is no default logger
since the asynchronous logging requires a separate goroutine.

*/
package alog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

type Fields map[string]interface{}

type Logger interface {
	// Print, Println, and Printf log a message without a level label and
	// without regard to any configured log level, if using leveled logging.
	// Arguments are handled in the manner of fmt.Print, fmt.Println, and
	// fmt.Printf respectively.
	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})

	// ---- Leveled Logging ----
	//
	// The following functions are for optional leveled logging.

	// Panic, Panicln, and Panicf log a message at PanicLevel and then call
	// panic() with the message.  Arguments are handled in the manner of
	// fmt.Print, fmt.Println, and fmt.Printf respectively.
	Panic(v ...interface{})
	Panicln(v ...interface{})
	Panicf(format string, v ...interface{})

	// Fatal, Fatalln, and Fatalf log a message at FatalLevel and then call
	// os.Exit(1).  Arguments are handled in the manner of fmt.Print,
	// fmt.Println, and fmt.Printf respectively.
	Fatal(v ...interface{})
	Fatalln(v ...interface{})
	Fatalf(format string, v ...interface{})

	// Error, Errorln, and Errorf log a message at ErrorLevel.  Arguments are
	// handled in the manner of fmt.Print, fmt.Println, and fmt.Printf
	// respectively.
	Error(v ...interface{})
	Errorln(v ...interface{})
	Errorf(format string, v ...interface{})

	// Warn, Warnln, and Warnf log a message at WarnLevel.  Arguments are
	// handled in the manner of fmt.Print, fmt.Println, and fmt.Printf
	// respectively.
	Warn(v ...interface{})
	Warnln(v ...interface{})
	Warnf(format string, v ...interface{})

	// Info, Infoln, and Infof log a message at InfoLevel.  Arguments are
	// handled in the manner of fmt.Print, fmt.Println, and fmt.Printf
	// respectively.
	Info(v ...interface{})
	Infoln(v ...interface{})
	Infof(format string, v ...interface{})

	// Debug, Debugln, and Debugf log a message at DebugLevel.  Arguments are
	// handled in the manner of fmt.Print, fmt.Println, and fmt.Printf
	// respectively.
	Debug(v ...interface{})
	Debugln(v ...interface{})
	Debugf(format string, v ...interface{})

	// WithFields creates a wrapper for the Logger that outputs each log
	// message with the specified fields included as part of the message.
	WithFields(fields Fields) Logger

	// WithField calls WithField for a single entry.
	WithField(key string, value interface{}) Logger

	// Close stops asynchronous logging and waits for any unwritten entries to
	// be written to the io.Writer.  This does not close the log's io.Writer,
	// and doing so it the caller's responsibility.  Do not call Close() while
	// there are goroutines that may be writing to the log.
	Close()
}

// New creates a new Logger instance that outputs log entries as text.
//
// The out variable sets the destination to which log data is written.
//
// Set level to NoLevel to choose not to do leveled logging.  Otherwise, set to
// the severity level to log at.
//
// The timeLayout defines the timestamp format according to time.Format.  If
// not specified, defaults to "Jan 02 15:04:05".  To disable timestamp output,
// specify a TimeLayout string consisting on one or more spaces. The prefix
// appears at the beginning of each generated log line.
func NewText(out io.Writer, level Level, timeLayout, prefix string) Logger {
	a := newLogger(out, level, timeLayout)
	a.writeFunc = a.writeText
	a.prefix = prefix
	go a.run()
	return a
}

// NewJSON creates a new Logger instance that outputs log entries as JSON.
//
// The out variable sets the destination to which log data is written.
//
// Set level to NoLevel to choose not to do leveled logging.  Otherwise, set to
// the severity level to log at.
//
// The timeLayout defines the timestamp format according to time.Format.  If
// not specified, defaults to "Jan 02 15:04:05".  To disable timestamp output,
// specify a TimeLayout string consisting on one or more spaces.
func NewJSON(out io.Writer, level Level, timeLayout string) Logger {
	a := newLogger(out, level, timeLayout)
	a.writeFunc = a.writeJSON
	go a.run()
	return a
}

const (
	// ErrorField is exported so it can be used explicitly as a field name.
	ErrorField = "error"

	levelField    = "level"
	extLevelField = "fields.level"
	msgField      = "msg"
	extMsgField   = "fields.msg"
	timeField     = "time"
	extTimeField  = "fields.time"
)

const defaultTimeLayout = "Jan 02 15:04:05"

func newLogger(out io.Writer, level Level, timeLayout string) *logger {
	if out == nil {
		out = os.Stdout
	}
	if timeLayout == "" {
		timeLayout = defaultTimeLayout
	} else {
		timeLayout = strings.TrimSpace(timeLayout)
	}
	if level < NoLevel {
		level = NoLevel
	} else if level > DebugLevel {
		level = DebugLevel
	}
	a := &logger{
		out:      out,
		entChan:  make(chan *entry, 64),
		doneChan: make(chan struct{}),
		level:    level,
		tsLayout: timeLayout,
	}
	runtime.SetFinalizer(a, closeLogger)
	return a
}

func closeLogger(a *logger) { close(a.entChan) }

// entry represents a single log entry that has not yet been written
type entry struct {
	ts     time.Time
	level  Level
	format string
	args   []interface{}
	fields Fields
	ln     bool
}

type logger struct {
	buf       []byte
	entChan   chan *entry
	doneChan  chan struct{}
	writeFunc func(*entry)
	out       io.Writer
	level     Level
	tsLayout  string
	prefix    string
}

func (a *logger) Print(v ...interface{}) {
	a.entChan <- &entry{ts: time.Now(), args: v}
}

func (a *logger) Println(v ...interface{}) {
	a.entChan <- &entry{ts: time.Now(), args: v, ln: true}
}

func (a *logger) Printf(format string, v ...interface{}) {
	a.entChan <- &entry{ts: time.Now(), format: format, args: v}
}

func (a *logger) WithFields(fields Fields) Logger {
	return &fieldLogger{
		logger: a,
		fields: fields,
	}
}

func (a *logger) WithField(key string, value interface{}) Logger {
	return a.WithFields(Fields{key: value})
}

func (a *logger) WithError(err error) Logger {
	return a.WithFields(Fields{ErrorField: err})
}

func (a *logger) Close() {
	close(a.entChan)
	<-a.doneChan
}

func (a *logger) run() {
	for ent := range a.entChan {
		a.writeFunc(ent)
	}
	close(a.doneChan)
}

func (a *logger) writeText(ent *entry) {
	a.buf = a.buf[:0]
	if a.prefix != "" {
		a.buf = append(a.buf, a.prefix...)
	}
	if a.tsLayout != "" {
		a.buf = append(a.buf, ent.ts.Format(a.tsLayout)...)
	}
	if a.level != NoLevel && ent.level != NoLevel {
		a.buf = append(a.buf, levelNamesText[int(ent.level)]...)
	} else {
		a.buf = append(a.buf, ' ')
	}
	if ent.format != "" {
		a.buf = append(a.buf, fmt.Sprintf(ent.format, ent.args...)...)
	} else if ent.ln {
		a.buf = append(a.buf, fmt.Sprintln(ent.args...)...)
		a.buf = a.buf[:len(a.buf)-1]
	} else {
		a.buf = append(a.buf, fmt.Sprint(ent.args...)...)
	}
	for k, v := range ent.fields {
		a.buf = append(a.buf, " ("...)
		a.buf = append(a.buf, k...)
		a.buf = append(a.buf, '=')
		a.buf = append(a.buf, fmt.Sprint(v)...)
		a.buf = append(a.buf, ')')
	}

	a.buf = append(a.buf, '\n')
	a.out.Write(a.buf)
}

func (a *logger) writeJSON(ent *entry) {
	if ent.fields == nil {
		ent.fields = make(map[string]interface{}, 3)
	}
	// Convert any error types to string.
	for k, v := range ent.fields {
		if err, ok := v.(error); ok {
			ent.fields[k] = err.Error()
		}
	}
	if a.tsLayout != "" {
		if v, ok := ent.fields[timeField]; ok {
			ent.fields[extTimeField] = v
		}
		ent.fields[timeField] = ent.ts.Format(a.tsLayout)
	}
	if a.level != NoLevel && ent.level != NoLevel {
		if v, ok := ent.fields[levelField]; ok {
			ent.fields[extLevelField] = v
		}
		ent.fields[levelField] = ent.level.String()
	}
	if ent.format != "" {
		if v, ok := ent.fields[msgField]; ok {
			ent.fields[extMsgField] = v
		}
		ent.fields[msgField] = fmt.Sprintf(ent.format, ent.args...)
	} else {
		if v, ok := ent.fields[msgField]; ok {
			ent.fields[extMsgField] = v
		}
		ent.fields[msgField] = fmt.Sprint(ent.args...)
	}
	encoded, err := json.Marshal(ent.fields)
	if err != nil {
		fmt.Println("Failed to marshal fields to JSON:", err)
		return
	}
	a.out.Write(append(encoded, '\n'))
}

// ---- Leveled log functions -----

// Log severity levels
const (
	// Leveled logging is disabled.  Messages are not filtered by level, and no
	// level label or field appears in the log messages.  This is the default
	// behavior.
	NoLevel Level = iota

	// PanicLevel is most severe level; it logs a message and calls panic.
	// This severity level indicates an unrecoverable condition caused by a
	// defect in programming logic or other situation not allowed by the
	// system or program.
	PanicLevel

	// FatalLevel logs and then calls os.Exit(1).  This severity level
	// indicates an unrecoverable condition that requires the termination of
	// the program.
	FatalLevel

	// ErrorLevel is used for error conditions or failures, usually
	// sufficiently critical to prevent the program from executing one or more
	// intended tasks.
	ErrorLevel

	// WarnLevel indicates undesirable conditions that should not normally
	// occur during proper execution with ideal configuration, but that are not
	// critical enough to stop the program from executing intended tasks.
	WarnLevel

	// InfoLevel indicated data that is informative, but its record is not
	// crucial under normal conditions.
	InfoLevel

	// DebugLevel labels detailed information the is intended for debugging
	// program logic or configuration.
	DebugLevel
)

// Level is the severity value for log entries.
type Level int

var levelNames = [DebugLevel + 1]string{
	"", "panic", "fatal", "error", "warn", "info", "debug"}

// String converts a Level value to a string containing the name of the level.
func (lvl Level) String() string { return levelNames[int(lvl)] }

var levelNamesText = [DebugLevel + 1]string{
	"", " PANIC ", " FATAL ", " ERROR ", " WARN ", " INFO ", " DEBUG "}

func (a *logger) Panic(v ...interface{}) {
	a.log(nil, PanicLevel, v)
	a.Close()
	panic(fmt.Sprint(v...))
}
func (a *logger) Panicln(v ...interface{}) {
	a.logln(nil, PanicLevel, v)
	a.Close()
	panic(fmt.Sprint(v...))
}
func (a *logger) Panicf(format string, v ...interface{}) {
	a.logf(nil, PanicLevel, format, v)
	a.Close()
	panic(fmt.Sprintf(format, v...))
}

func (a *logger) Fatal(v ...interface{}) {
	a.log(nil, FatalLevel, v)
	a.Close()
	os.Exit(1)
}
func (a *logger) Fatalln(v ...interface{}) {
	a.logln(nil, FatalLevel, v)
	a.Close()
	os.Exit(1)
}
func (a *logger) Fatalf(format string, v ...interface{}) {
	a.logf(nil, FatalLevel, format, v)
	a.Close()
	os.Exit(1)
}

func (a *logger) Error(v ...interface{})   { a.log(nil, ErrorLevel, v) }
func (a *logger) Errorln(v ...interface{}) { a.logln(nil, ErrorLevel, v) }
func (a *logger) Errorf(format string, v ...interface{}) {
	a.logf(nil, ErrorLevel, format, v)
}

func (a *logger) Warn(v ...interface{})   { a.log(nil, WarnLevel, v) }
func (a *logger) Warnln(v ...interface{}) { a.logln(nil, WarnLevel, v) }
func (a *logger) Warnf(format string, v ...interface{}) {
	a.logf(nil, WarnLevel, format, v)
}

func (a *logger) Info(v ...interface{})   { a.log(nil, InfoLevel, v) }
func (a *logger) Infoln(v ...interface{}) { a.logln(nil, InfoLevel, v) }
func (a *logger) Infof(format string, v ...interface{}) {
	a.logf(nil, InfoLevel, format, v)
}

func (a *logger) Debug(v ...interface{})   { a.log(nil, DebugLevel, v) }
func (a *logger) Debugln(v ...interface{}) { a.logln(nil, DebugLevel, v) }
func (a *logger) Debugf(format string, v ...interface{}) {
	a.logf(nil, DebugLevel, format, v)
}

func (a *logger) log(fields Fields, level Level, v []interface{}) {
	if !a.LogableAt(level) {
		return
	}
	a.entChan <- &entry{
		ts:     time.Now(),
		level:  level,
		args:   v,
		fields: fields,
	}
}
func (a *logger) logln(fields Fields, level Level, v []interface{}) {
	if !a.LogableAt(level) {
		return
	}
	a.entChan <- &entry{
		ts:     time.Now(),
		level:  level,
		args:   v,
		fields: fields,
		ln:     true,
	}
}
func (a *logger) logf(fields Fields, level Level, format string, v []interface{}) {
	if !a.LogableAt(level) {
		return
	}
	a.entChan <- &entry{
		ts:     time.Now(),
		level:  level,
		format: format,
		args:   v,
		fields: fields,
	}
}

func (a *logger) LogableAt(level Level) bool {
	if a.level != NoLevel && a.level < level {
		return false
	}
	return true
}
