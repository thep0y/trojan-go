//go:build server || full || mini
// +build server full mini

package build

import (
	_ "github.com/thep0y/trojan-go/proxy/server"
)
