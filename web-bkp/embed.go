// Package ui handles the PocketBase Admin frontend embedding.
package web

import (
	"embed"
	"io/fs"

	"github.com/gin-contrib/static"
)

//go:embed all:dist
var distDir embed.FS

// DistDirFS contains the embedded dist directory files (without the "dist" prefix)
var DistFS, _ = fs.Sub(distDir, "dist")
var StaticDir = static.LocalFile("dist", true)
