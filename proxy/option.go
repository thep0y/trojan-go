package proxy

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/thep0y/trojan-go/common"
	"github.com/thep0y/trojan-go/constant"
	"github.com/thep0y/trojan-go/option"
)

type Option struct {
	path *string
}

func (o *Option) Name() string {
	return Name
}

func detectAndReadConfig(file string) ([]byte, bool, error) {
	isJSON := false
	switch {
	case strings.HasSuffix(file, ".json"):
		isJSON = true
	case strings.HasSuffix(file, ".yaml"), strings.HasSuffix(file, ".yml"):
		isJSON = false
	default:
		log.Fatal().
			Str("file", file).
			Msg("unsupported config format, use .yaml or .json instead")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, false, err
	}
	return data, isJSON, nil
}

func (o *Option) Handle() error {
	defaultConfigPath := []string{
		"config.json",
		"config.yml",
		"config.yaml",
	}

	isJSON := false
	var data []byte
	var err error

	switch *o.path {
	case "":
		log.Warn().Msg("no specified config file, use default path to detect config file")
		for _, file := range defaultConfigPath {
			log.Warn().Str("file-path", file).Msg("try to load config from default path")
			data, isJSON, err = detectAndReadConfig(file)
			if err != nil {
				log.Warn().Err(err).Send()
				continue
			}
			break
		}
	default:
		data, isJSON, err = detectAndReadConfig(*o.path)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	}

	if data != nil {
		log.Info().
			Str("version", constant.Version).
			Msg("trojan-go is initializing")
		proxy, err := NewProxyFromConfigData(data, isJSON)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		err = proxy.Run()
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	}

	log.Fatal().Msg("no valid config")
	return nil
}

func (o *Option) Priority() int {
	return -1
}

func init() {
	option.RegisterHandler(&Option{
		path: flag.String("config", "", "Trojan-Go config filename (.yaml/.yml/.json)"),
	})
	option.RegisterHandler(&StdinOption{
		format: flag.String(
			"stdin-format",
			"disabled",
			"Read from standard input (yaml/json)",
		),
		suppressHint: flag.Bool("stdin-suppress-hint", false, "Suppress hint text"),
	})
}

type StdinOption struct {
	format       *string
	suppressHint *bool
}

func (o *StdinOption) Name() string {
	return Name + "_STDIN"
}

func (o *StdinOption) Handle() error {
	isJSON, e := o.isFormatJson()
	if e != nil {
		return e
	}

	if o.suppressHint == nil || !*o.suppressHint {
		fmt.Printf("Trojan-Go %s (%s/%s)\n", constant.Version, runtime.GOOS, runtime.GOARCH)
		if isJSON {
			fmt.Println("Reading JSON configuration from stdin.")
		} else {
			fmt.Println("Reading YAML configuration from stdin.")
		}
	}

	data, e := io.ReadAll(bufio.NewReader(os.Stdin))
	if e != nil {
		log.Fatal().Err(e).Msg("Failed to read from stdin")
	}

	proxy, err := NewProxyFromConfigData(data, isJSON)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	err = proxy.Run()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return nil
}

func (o *StdinOption) Priority() int {
	return 0
}

func (o *StdinOption) isFormatJson() (isJson bool, e error) {
	if o.format == nil {
		return false, common.NewError("format specifier is nil")
	}
	if *o.format == "disabled" {
		return false, common.NewError("reading from stdin is disabled")
	}
	return strings.ToLower(*o.format) == "json", nil
}
