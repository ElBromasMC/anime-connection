//go:build !dev

package assets

import (
	"embed"
)

//go:embed static
var Assets embed.FS
