package fingerprint

import (
	"crypto/tls"

	"github.com/rs/zerolog/log"
)

func ParseCipher(s []string) []uint16 {
	all := tls.CipherSuites()
	var result []uint16
	for _, p := range s {
		found := true
		for _, q := range all {
			if q.Name == p {
				result = append(result, q.ID)
				break
			}
			if !found {
				log.Warn().Str("p", p).Msg("skipped invalid cipher suite")
			}
		}
	}
	return result
}
