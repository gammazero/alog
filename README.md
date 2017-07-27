# alog
Asynchronous structured leveled logging

[![GoDoc](https://godoc.org/github.com/gammazero/alog?status.svg)](https://godoc.org/github.com/gammazero/alog)

The alog package provides structured logging with JSON output and with semi-structured text-based logging.  Logging is asynchronous, where work and time to format log messages and write them to I/O is done by an asynchronous goroutine, minimizing the time the application waits for a log call to return.  Decisions to log in a normal, especially a performance critical, code path should still be
carefully considered.

A Logger instance generates lines of output to an io.Writer.  Each logging operation makes a single call to the Writer's Write method.  A Logger can be used simultaneously from multiple goroutines; it guarantees serialized access to the Writer.

The use of structured and leveled logging is completely optional.  Leveled logging is disabled by specifying `alog.NoLevel` when creating a new logger instance.  If considering using leveled logging, please read this discussion: 

https://dave.cheney.net/2015/11/05/lets-talk-about-logging

## Example

Using `alog` is simple: create an instance and start logging messages.

```go
package main

import (
    "github.com/gammazero/alog"
)

func main() {
    logger := alog.NewText(nil, alog.InfoLevel, "", "")
    logger.LevelPrint(InfoLevel, "hello world")
    logger.WithFields(map[string]interface{}{
        "hero": "rick",
        "sidekick": "morty",
    }).Error("Portal malfunction")
    logger.Close()
}
```

By default, logs are written to stdout (when a `nil` `io.Writer` is specified).

## Log Format
The default logging output is semi-structured text based.  The above produces this output:

```text
Jul 22 02:43:36 INFO hello world 
Jul 22 02:43:36 ERROR Portal malfunction (hero=rick) (sidekick=morty)
```

You can output log data as structured JSON:

```go
    logger := alog.NewJSON(os.Stderr, alog.DebugLevel, "")
    ...
```

Using that in the first block of code will produce the output:

```json
{"level":"info","msg":"hello world","time":"Jul 22 02:43:36"}
{"hero":"rick","level":"error","msg":"Portal malfunction","sidekick":"morty","time":"Jul 22 02:43:36"}
```

## Default Logger

Using alog requires creating a logger instance.  There is no default logger since the asynchronous logging must run a separate goroutine.  To use alog in a manner similar to the default logger create a global alog instance named `log`:

```go
    var log = alog.NewText(os.Stdout, alog.NoLevel, time.RFC3339,"")

    func main() {
        defer log.Close()
    }
```

## Compatibility

For compatibility with the stdlib logger, a alog instance provides a limited set of compatible logging functions.
