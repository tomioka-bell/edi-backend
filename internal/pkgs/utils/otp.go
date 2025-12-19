package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func New6DigitCode() (plain string, hash string, err error) {
	// สุ่ม 0-999999 แล้ว zero-pad เป็น 6 หลักแบบไม่ bias
	var b [4]byte
	if _, err = rand.Read(b[:]); err != nil {
		return "", "", err
	}
	n := (int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])) & 0x7fffffff
	code := n % 1000000
	plain = fmt.Sprintf("%06d", code)

	h := sha256.Sum256([]byte(plain))
	hash = hex.EncodeToString(h[:])
	return
}

func HashCode(plain string) string {
	h := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(h[:])
}
