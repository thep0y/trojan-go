package version

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/thep0y/trojan-go/common"
	"github.com/thep0y/trojan-go/constant"
	"github.com/thep0y/trojan-go/option"
)

type versionOption struct {
	flag *bool
}

func (*versionOption) Name() string {
	return "version"
}

func (*versionOption) Priority() int {
	return 10
}

func (c *versionOption) Handle() error {
	if *c.flag {
		fmt.Println("Trojan-Go", constant.Version)
		fmt.Println("Go Version:", runtime.Version())
		fmt.Println("OS/Arch:", runtime.GOOS+"/"+runtime.GOARCH)
		fmt.Println("Git Commit:", constant.Commit)
		fmt.Println("")
		fmt.Println("Developed by PageFault (thep0y)")
		fmt.Println("Licensed under GNU General Public License version 3")
		fmt.Println("GitHub Repository:\thttps://github.com/thep0y/trojan-go")
		fmt.Println("Trojan-Go Documents:\thttps://thep0y.github.io/trojan-go/")
		return nil
	}
	return common.NewError("not set")
}

func init() {
	option.RegisterHandler(&versionOption{
		flag: flag.Bool("version", false, "Display version and help info"),
	})
}
