package logging

import (
	"io"
	"os"
	"time"

	"router-practice/variable"

	"github.com/rs/zerolog"
)

var (
	Object zerolog.Logger
	F      *os.File
)

func SetupLogger() {
	zerolog.TimeFieldFormat = "20060102150405"
	zerolog.TimestampFieldName = "datetime"

	fname := variable.ProgramName + "-" + time.Now().Format("20060102") + ".log"
	F, _ = os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	w := io.MultiWriter(F)

	Object = zerolog.New(w).With().Logger()
}

func RenewLogger() {
	zerolog.TimeFieldFormat = "20060102150405"
	zerolog.TimestampFieldName = "datetime"

	fname := variable.ProgramName + "-" + time.Now().Format("20060102") + ".log"

	F.Close()

	F, _ = os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	w := io.MultiWriter(F)

	Object = zerolog.New(w).With().Logger()
}
