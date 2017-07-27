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

func BenchmarkLogger(b *testing.B) {
	lg := NewText(ioutil.Discard, NoLevel, "", "")
	lg2 := NewText(ioutil.Discard, NoLevel, "", "")
	j := junk{42, "xyzzy", 3.14159}

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

func BenchmarkStdlibLog(b *testing.B) {
	lg := log.New(ioutil.Discard, "stdlog", log.LstdFlags)
	lg2 := log.New(ioutil.Discard, "stdlog", log.LstdFlags)
	j := junk{42, "xyzzy", 3.14159}

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
