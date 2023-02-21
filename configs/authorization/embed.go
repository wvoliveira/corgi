package authorization

import (
	"embed"
)

//go:embed model.conf policy.csv
var DistFiles embed.FS
