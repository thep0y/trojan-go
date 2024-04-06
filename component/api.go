//go:build api || full
// +build api full

package build

import (
	_ "github.com/thep0y/trojan-go/api/control"
	_ "github.com/thep0y/trojan-go/api/service"
)
