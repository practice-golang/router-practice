package variable

import (
	"embed"

	"github.com/rs/zerolog"
)

var (
	ProgramName string = "router-practice"

	Content embed.FS
	Static  embed.FS

	Logger zerolog.Logger
)
