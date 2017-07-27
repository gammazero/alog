package alog

import (
	"fmt"
	"os"
	"time"
)

type fieldLogger struct {
	*logger
	fields Fields
}

func (f *fieldLogger) Print(v ...interface{}) {
	f.entChan <- &entry{ts: time.Now(), args: v, fields: f.fields}
}

func (f *fieldLogger) Println(v ...interface{}) {
	f.entChan <- &entry{ts: time.Now(), args: v, ln: true, fields: f.fields}
}

func (f *fieldLogger) Printf(format string, v ...interface{}) {
	f.entChan <- &entry{ts: time.Now(), format: format, args: v, fields: f.fields}
}

func (f *fieldLogger) WithFields(fields Fields) Logger {
	// Create new fieldLogger with parents and specified fields.
	newFields := make(Fields, len(f.fields)+len(fields))
	for k, v := range f.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}
	return &fieldLogger{
		logger: f.logger,
		fields: newFields,
	}
}

func (f *fieldLogger) WithField(key string, value interface{}) Logger {
	return f.WithFields(Fields{key: value})
}

// ---- Leveled log functions -----

func (f *fieldLogger) Panic(v ...interface{}) {
	f.log(f.fields, PanicLevel, v)
	f.Close()
	panic(fmt.Sprint(v...))
}
func (f *fieldLogger) Panicln(v ...interface{}) {
	f.logln(f.fields, PanicLevel, v)
	f.Close()
	panic(fmt.Sprint(v...))
}
func (f *fieldLogger) Panicf(format string, v ...interface{}) {
	f.logf(f.fields, PanicLevel, format, v)
	f.Close()
	panic(fmt.Sprintf(format, v...))
}

func (f *fieldLogger) Fatal(v ...interface{}) {
	f.log(f.fields, FatalLevel, v)
	f.Close()
	os.Exit(1)
}
func (f *fieldLogger) Fatalln(v ...interface{}) {
	f.logln(f.fields, FatalLevel, v)
	f.Close()
	os.Exit(1)
}
func (f *fieldLogger) Fatalf(format string, v ...interface{}) {
	f.logf(f.fields, FatalLevel, format, v)
	f.Close()
	os.Exit(1)
}

func (f *fieldLogger) Error(v ...interface{}) {
	f.log(f.fields, ErrorLevel, v)
}
func (f *fieldLogger) Errorln(v ...interface{}) {
	f.logln(f.fields, ErrorLevel, v)
}
func (f *fieldLogger) Errorf(format string, v ...interface{}) {
	f.logf(f.fields, ErrorLevel, format, v)
}

func (f *fieldLogger) Warn(v ...interface{}) {
	f.log(f.fields, WarnLevel, v)
}
func (f *fieldLogger) Warnln(v ...interface{}) {
	f.logln(f.fields, WarnLevel, v)
}
func (f *fieldLogger) Warnf(format string, v ...interface{}) {
	f.logf(f.fields, WarnLevel, format, v)
}

func (f *fieldLogger) Info(v ...interface{}) {
	f.log(f.fields, InfoLevel, v)
}
func (f *fieldLogger) Infoln(v ...interface{}) {
	f.logln(f.fields, InfoLevel, v)
}
func (f *fieldLogger) Infof(format string, v ...interface{}) {
	f.logf(f.fields, InfoLevel, format, v)
}

func (f *fieldLogger) Debug(v ...interface{}) {
	f.log(f.fields, DebugLevel, v)
}
func (f *fieldLogger) Debugln(v ...interface{}) {
	f.logln(f.fields, DebugLevel, v)
}
func (f *fieldLogger) Debugf(format string, v ...interface{}) {
	f.logf(f.fields, DebugLevel, format, v)
}
