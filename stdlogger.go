package alog

// StdLogger is an interface implemented by both alog.Logger and stdlib
// log.Logger.  This allows both implementations to be used interchangeably.
type StdLogger interface {
	// Print, Println, and Printf log a message.  Arguments are handled in the
	// manner of fmt.Print, fmt.Println, and fmt.Printf respectively.
	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})

	// Fatal, Fatalln, and Fatalf log a message and then call os.Exit(1).
	// Arguments are handled in the manner of fmt.Print, fmt.Println, and
	// fmt.Printf respectively.  These are provided to allow alog.Logger to
	// serve as a replacement for stdlib log.Logger.
	Fatal(v ...interface{})
	Fatalln(v ...interface{})
	Fatalf(format string, v ...interface{})

	// Panic, Panicln, and Panicf log a message and then call panic() with the
	// message.  Arguments are handled in the manner of fmt.Print, fmt.Println,
	// and fmt.Printf respectively.  These are provided to allow alog.Logger to
	// serve as a replacement for stdlib log.Logger.
	Panic(v ...interface{})
	Panicln(v ...interface{})
	Panicf(format string, v ...interface{})
}
