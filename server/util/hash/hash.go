package hash

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

func RandSha1(var1 string) []byte {
	r1 := time.Now().UnixNano()
	r2 := rand.Int()
	r3 := rand.Int()

	h := sha1.New()
	h.Write([]byte(strconv.FormatInt(r1, 16) + var1 + strconv.Itoa(r2) + strconv.Itoa(r3)))
	return []byte(hex.EncodeToString(h.Sum(nil)))
}
