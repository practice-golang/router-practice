package logging

import (
	"io"
	"os"
	"router-practice/variable"
	"time"

	"github.com/rs/zerolog"
)

var (
	Object zerolog.Logger
)

func SetupLogger() {
	zerolog.TimeFieldFormat = "20060102150405"
	zerolog.TimestampFieldName = "datetime"

	fname := variable.ProgramName + "-" + time.Now().Format("20060102") + ".log"
	f, _ := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	// w := io.MultiWriter(os.Stdout, f)
	w := io.MultiWriter(f)
	// w := io.Writer(f)

	Object = zerolog.New(w).With().Logger()
}