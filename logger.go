/*
Package alog provides an asynchronous logger.

The alog package provides simple, fast asynchronous logging.  The work and time
to format log messages and write them to I/O is done by an asynchronous
goroutine, allowing the calling application to continue operating without
waiting for logging.

The alog package provides a StdLogger interface which is implemented by both a
logger as well as the stdlib log.  Libraries that perform logging can use this
interface so that either the alog or the stdlib *Logger can be used
interchangeably.

Using alog requires creating a logger instance.  There is no default logger
since the asynchronous logging requires a separate goroutine.

For compatibility with the stdlib logger, create a global alog instance
named `log`:

    // Log can be reassigned a stdlib *Logger.
    var log = alog.New(os.Stdout, "", "")

    func main() {
        defer func() {
            if logger, ok := log.(*alog.Logger); ok {
                logger.Close()
            }
        }()
    }
*/
package alog

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// StdLogger is an interface implemented by both alog.Logger and stdlib log.Logger.  This allows both implementations to be used interchangeably.
type StdLogger interface {
	// Panic, Panicln, and Panicf log a message and then call panic() with the
	// message.  Arguments are handled in the manner of fmt.Print, fmt.Println,
	// and fmt.Printf respectively.
	Panic(v ...interface{})
	Panicln(v ...interface{})
	Panicf(format string, v ...interface{})

	// Fatal, Fatalln, and Fatalf log a message and then call os.Exit(1).
	// Arguments are handled in the manner of fmt.Print, fmt.Println, and
	// fmt.Printf respectively.
	Fatal(v ...interface{})
	Fatalln(v ...interface{})
	Fatalf(format string, v ...interface{})

	// Print, Println, and Printf log a message.  Arguments are handled in the
	// manner of fmt.Print, fmt.Println, and fmt.Printf respectively.
	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

const defaultTimeLayout = "Jan 02 15:04:05"

// New creates a new Logger.  The out variable sets the destination to which
// log data will be written.  The prefix appears at the beginning of each
// generated log line.  The timeLayout defines the timestamp format according
// to time.Format.  If not specified, defaults to "Jan 02 15:04:05".  To
// disable timestamp output, specify a TimeLayout string consisting on one or
// more spaces.
func New(out io.Writer, prefix, timeLayout string) *Logger {
	if out == nil {
		out = os.Stderr
	}

	if timeLayout == "" {
		timeLayout = defaultTimeLayout
	} else {
		timeLayout = strings.TrimSpace(timeLayout)
	}

	a := &Logger{
		out:      out,
		entChan:  make(chan *entry, 64),
		doneChan: make(chan struct{}),
		tsLayout: timeLayout,
		prefix:   prefix,
	}
	go a.run()
	return a
}

// entry represents a single log entry that has not yet been written
type entry struct {
	ts     time.Time
	format string
	args   []interface{}
	ln     bool
}

// A Logger represents an active logging object that generates lines of output
// to an io.Writer.  Each logging operation makes a single call to the Writer's
// Write method.  A Logger can be used simultaneously from multiple goroutines;
// it guarantees to serialize access to the Writer.
type Logger struct {
	out      io.Writer
	buf      []byte
	entChan  chan *entry
	doneChan chan struct{}
	tsLayout string
	prefix   string
}

// Panic, logs a message and then calls panic() with the message.  Arguments
// are handled in the manner of fmt.Print.
func (a *Logger) Panic(v ...interface{}) {
	a.Print(v...)
	a.Close()
	panic(fmt.Sprint(v...))
}

// Panicln, logs a message and then calls panic() with the message.  Arguments
// are handled in the manner of fmt.Println.
func (a *Logger) Panicln(v ...interface{}) {
	a.Println(v...)
	a.Close()
	panic(fmt.Sprint(v...))
}

// Panicf, logs a message and then calls panic() with the message.  Arguments
// are handled in the manner of fmt.Printf.
func (a *Logger) Panicf(format string, v ...interface{}) {
	a.Printf(format, v...)
	a.Close()
	panic(fmt.Sprintf(format, v...))
}

// Fatal logs a message and then calls os.Exit(1).  Arguments are handled in
// the manner of fmt.Print.
func (a *Logger) Fatal(v ...interface{}) {
	a.Print(v...)
	a.Close()
	os.Exit(1)
}

// Fatalln logs a message and then calls os.Exit(1).  Arguments are handled in
// the manner of fmt.Println.
func (a *Logger) Fatalln(v ...interface{}) {
	a.Println(v...)
	a.Close()
	os.Exit(1)
}

// Fatalf logs a message and then calls os.Exit(1).  Arguments are handled in
// the manner of fmt.Printf.
func (a *Logger) Fatalf(format string, v ...interface{}) {
	a.Printf(format, v...)
	a.Close()
	os.Exit(1)
}

// Print logs a message.  Arguments are handled in the manner of fmt.Print.
func (a *Logger) Print(v ...interface{}) {
	a.entChan <- &entry{
		ts:   time.Now(),
		args: v,
	}
}

// Println logs a message.  Arguments are handled in the manner of fmt.Println.
func (a *Logger) Println(v ...interface{}) {
	a.entChan <- &entry{
		ts:   time.Now(),
		args: v,
		ln:   true,
	}
}

// Printf logs a message.  Arguments are handled in the manner of fmt.Printf.
func (a *Logger) Printf(format string, v ...interface{}) {
	a.entChan <- &entry{
		ts:     time.Now(),
		format: format,
		args:   v,
	}
}

func (a *Logger) run() {
	for ent := range a.entChan {
		if ent == nil {
			break
		}
		a.output(ent)
	}
	close(a.doneChan)
}

// Close stops asynchronous logging and waits for any buffered data to be
// flushed.
func (a *Logger) Close() {
	// Exit async logging loop without closing channel.
	a.entChan <- nil
	<-a.doneChan
}

func (a *Logger) output(ent *entry) {
	a.buf = a.buf[:0]
	if a.prefix != "" {
		a.buf = append(a.buf, a.prefix...)
	}
	if a.tsLayout != "" {
		a.buf = append(a.buf, ent.ts.Format(a.tsLayout)...)
		a.buf = append(a.buf, ' ')
	}
	if ent.format != "" {
		a.buf = append(a.buf, fmt.Sprintf(ent.format, ent.args...)...)
		a.buf = append(a.buf, '\n')
	} else if ent.ln {
		a.buf = append(a.buf, fmt.Sprintln(ent.args...)...)
	} else {
		a.buf = append(a.buf, fmt.Sprint(ent.args...)...)
		a.buf = append(a.buf, '\n')
	}
	a.out.Write(a.buf)
}
