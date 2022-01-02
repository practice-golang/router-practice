package router

import "embed"

var (
	Content    embed.FS
	Static     embed.FS
	AllMethods = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
)
