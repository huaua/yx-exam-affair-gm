package webres

import (
	"embed"
)

//go:embed html/index.html
var Html []byte

//go:embed all:html/assets
var Assets embed.FS

//go:embed all:html/static
var Static embed.FS
