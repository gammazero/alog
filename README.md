# alog
Asynchronous logging

[![GoDoc](https://godoc.org/github.com/gammazero/alog?status.svg)](https://godoc.org/github.com/gammazero/alog)

The alog package provides simple, fast asynchronous logging.  The work and time to format log messages and write them to I/O is done by an asynchronous goroutine, allowing the calling application to continue operating without waiting for logging. 


## Example

Using `alog` is a easy as using the stdlib logger: create an instance and start logging messages.

```go
package main

import (
    "github.com/gammazero/alog"
)

func main() {
    logger := alog.New(os.Stderr, "", "")
    logger.Print("hello world")
    logger.Close()
}
```

## Default Logger

Using `alog` requires creating a `Logger` instance.  There is no default logger since the asynchronous logging must run a separate goroutine.

## Compatibility

The `alog.StdLogger` is implemented by both `alog.Logger` and the stdlib `log.Logger`.  This allows code to that uses `alog.StdLogger` to work with either implementation interchangeably.

It is recommended that packages that want to use asynchronous logging create a global instance of alog.StdLogger, which can be assigned an instance of `alog.Logger` or stdlib `log.Logger`.

```go
package main

import (
    "github.com/gammazero/alog"
)

// Log can be reassigned a stdlib *Logger.
var log = alog.New(os.Stdout, "", "")

func main() {
    defer func() {
        if logger, ok := log.(*alog.Logger); ok {
            logger.Close()
        }
    }()
}
```

## Design

The alog package was designed to be:

- Compatible with the stdlib `log` package.
- Simple (few options, minimal interface).
- Fast (asynchronous, no complicated steps to generating a log entry).


### Philosophy

The design philosophy also follows this discussion:
https://dave.cheney.net/2015/11/05/lets-talk-about-logging

#### No separate levels (paraphrasing above)

'''No Warning Level'''

"Eliminate the warning level, it's either an informational message, or an error condition."

'''No Fatal or Panic Level'''

"In effect, log.Fatal is a less verbose than, but semantically equivalent to, panic. It is commonly accepted that libraries should not use panic...  Don't log at fatal level, prefer instead to return an error to the caller."

'''No Error Level'''

"You should never be logging anything at error level because you should either handle the error, or pass it back to the caller."

"If you choose to handle the error by logging it, by definition it's not an error any more - you handled it. The act of logging an error handles the error, hence it is no longer appropriate to log it as an error."

'''No Debug Level'''

Debug log messages are things that developers care about when they are developing or debugging software.  Generating these messages, including extra work to gather/calculate content for these message, should not be part of normal program execution.  Therefore, enabling or disabling debug should be done at the application level as it may need to enable/disable more then simply logging the messages.  So, the logger has no real need for debug level.

'''No need filter by severity level'''

There should not be an option to turn informational logging off as the user should only be told things which are useful for them.

Specific log data can be searched for by any number of tools.  If some form of categorization is needed, this is best left to the application creating log content, to include some identifier in the log message, instead of the log library trying to guess how the caller may want to include such information.
