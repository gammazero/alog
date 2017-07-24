package alog

import (
	"io/ioutil"
	"log"
	"testing"
)

type junk struct {
	foo int
	bar string
	baz float64
}

func BenchmarkLoggerVsLog(b *testing.B) {
	j := junk{42, "xyzzy", 3.14159}
	lg := New(ioutil.Discard, "", "")
	lg2 := New(ioutil.Discard, "", "")

	words := []string{"for", "score", "and", "seven", "years", "ago"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lg.Print("Animal Types", "shark", "=", "fish", "guppy", "=", "fish",
			"batray", "=", "fish")
		lg.Print(words)
		lg.Println(4, 20, "and", 7, []byte("years"), "ago")
		lg.Printf("junk: %+v", j)

		lg2.Print("Hello world")
		lg2.Printf("Numbers: %d %d %d", 1, 2, 3)
		lg2.Println("for", "score", "and", "seven", "years", "ago")
	}
	lg.Close()
}

func BenchmarkLog(b *testing.B) {
	j := junk{42, "xyzzy", 3.14159}
	lg := log.New(ioutil.Discard, "stdlog", log.LstdFlags)
	lg2 := log.New(ioutil.Discard, "stdlog", log.LstdFlags)

	words := []string{"for", "score", "and", "seven", "years", "ago"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lg.Print("Animal Types", "shark", "=", "fish", "guppy", "=", "fish",
			"batray", "=", "fish")
		lg.Print(words)
		lg.Println(4, 20, "and", 7, []byte("years"), "ago")
		lg.Printf("junk: %+v", j)

		lg2.Print("Hello world")
		lg2.Printf("Numbers: %d %d %d", 1, 2, 3)
		lg2.Println("for", "score", "and", "seven", "years", "ago")
	}
}
