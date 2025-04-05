package rand

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"math/rand"
	"strconv"
	"time"
)

func RandMD5(var1 string) []byte { // 16
	hash := md5.Sum(generateUniqueByteArr(var1))
	return hash[:]
}

func RandSha1(var1 string) []byte { // 20
	hash := sha1.Sum(generateUniqueByteArr(var1))
	return hash[:]
}

func RandSha256(var1 string) []byte { // 32
	hash := sha256.Sum256(generateUniqueByteArr(var1))
	return hash[:]
}

func RandSha512(var1 string) []byte { // 64
	hash := sha512.Sum512(generateUniqueByteArr(var1))
	return hash[:]
}

func generateUniqueByteArr(var1 string) []byte {
	r1 := time.Now().UnixNano()
	r2 := rand.Int()
	r3 := rand.Int()
	return []byte(strconv.FormatInt(r1, 16) + var1 + strconv.Itoa(r2) + strconv.Itoa(r3))
}
