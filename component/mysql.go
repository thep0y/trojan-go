//go:build mysql || full || mini
// +build mysql full mini

package build

import (
	_ "github.com/thep0y/trojan-go/statistic/mysql"
)
