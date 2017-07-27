package alog

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestLoggerLevels(t *testing.T) {
	buf := new(bytes.Buffer)
	lg := NewText(buf, DebugLevel, "", "happylog ")
	lg.Debug("at debug")
	lg.Info("at info")
	lg.Warn("at warn")
	lg.Error("at error")
	time.Sleep(100 * time.Millisecond)

	s := buf.String()
	if !strings.Contains(s, "DEBUG at debug") {
		t.Error("bad debug log")
	}
	if !strings.Contains(s, "INFO at info") {
		t.Error("bad info log")
	}
	if !strings.Contains(s, "WARN at warn") {
		t.Error("bad warn log")
	}
	if !strings.Contains(s, "ERROR at error") {
		t.Error("bad error log")
	}

	fmt.Println(s)
}

func TestFields(t *testing.T) {
	buf := new(bytes.Buffer)
	lg := NewText(buf, NoLevel, "", "")

	flg := lg.WithFields(Fields{
		"foo": "bar",
		"baz": "quz",
	})

	fflg := flg.WithFields(Fields{
		"July": 7,
	})
	fflg = fflg.WithField("December", 12)

	lg.Info("hello")
	time.Sleep(100 * time.Millisecond)
	s := string(buf.Next(4096))
	if strings.Contains(s, "(foo=bar) (baz=quz)") {
		t.Fatal("message should not contain fileds")
	}

	flg.Info("fields")
	time.Sleep(100 * time.Millisecond)
	s = string(buf.Next(4096))
	if !strings.Contains(s, "(foo=bar) (baz=quz)") && !strings.Contains(s, "(baz=quz) (foo=bar)") {
		t.Fatal("missing or badly formatter fields in message:", s)
	}
	if strings.Contains(s, "(July=7)") || strings.Contains(s, "(December=12)") {
		t.Fatal("log has fields that should not be there")
	}
	fmt.Print(s)

	fflg.Info("fields2")
	time.Sleep(100 * time.Millisecond)
	s = string(buf.Next(4096))
	if !strings.Contains(s, "(foo=bar)") || !strings.Contains(s, "(baz=quz)") || !strings.Contains(s, "(July=7)") || !strings.Contains(s, "(December=12)") {
		t.Fatal("missing fields in message:", s)
	}
	fmt.Print(s)

	lg.Info("byebye")
	time.Sleep(100 * time.Millisecond)
	s = string(buf.Next(4096))
	if strings.Contains(s, "(foo=bar) (baz=quz)") {
		t.Fatal("message should not contain fileds")
	}
}
