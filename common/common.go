package common

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type Runnable interface {
	Run() error
	Close() error
}

func SHA224String(password string) string {
	hash := sha256.New224()
	hash.Write([]byte(password))
	val := hash.Sum(nil)
	str := ""
	for _, v := range val {
		str += fmt.Sprintf("%02x", v)
	}
	return str
}

func GetProgramDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	return dir
}

func GetAssetLocation(file string) string {
	if filepath.IsAbs(file) {
		return file
	}
	if loc := os.Getenv("TROJAN_GO_LOCATION_ASSET"); loc != "" {
		absPath, err := filepath.Abs(loc)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		log.Debug().
			Str("TROJAN_GO_LOCATION_ASSET", absPath).
			Msg("env set")
		return filepath.Join(absPath, file)
	}
	return filepath.Join(GetProgramDir(), file)
}
