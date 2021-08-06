package api

import (
	"crypto/sha256"
)

func encodePwd(passwd string) string {
	h := sha256.New()
	h.Write([]byte(passwd))
	return string(h.Sum(nil))
}
