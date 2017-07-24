package alog

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoggerLevels(t *testing.T) {
	buf := new(bytes.Buffer)
	lg := New(buf, "happylog ", "")
	defer lg.Close()
	lg.Print("one")
	lg.Println("two")
	lg.Printf("three")
	time.Sleep(100 * time.Millisecond)

	s := buf.String()
	if !strings.Contains(s, "one") {
		t.Error("missng one")
	}
	if !strings.Contains(s, "two") {
		t.Error("missng two")
	}
	if !strings.Contains(s, "three") {
		t.Error("missng three")
	}
}

func TestCompat(t *testing.T) {
	logit := func(logger StdLogger) {
		logger.Print("one")
		logger.Println("two")
		logger.Printf("three")
	}

	lgStd := log.New(os.Stdout, "", log.LstdFlags)
	lg := New(os.Stdout, "", "")

	logit(lgStd)
	logit(lg)
	lg.Close()
}
