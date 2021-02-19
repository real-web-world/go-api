package bdk

import (
	"math/rand"
	"os"
	"regexp"
	"time"
)

var (
	maskPhoneReg = regexp.MustCompile("(\\d{3})(\\d{4})(\\d{4})")
	wxUaReg      = regexp.MustCompile("MicroMessenger")
	phoneReg     = regexp.MustCompile("1\\d{10}")
)

func IsFile(filename string) bool {
	fd, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !fd.Mode().IsDir()
}
func RandomAlphaNum(lengthParam ...int) []byte {
	length := 16
	if len(lengthParam) > 0 {
		length = lengthParam[0]
	}
	bytes := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return result
}
func InArrayInt(a int, arr []int) bool {
	for _, v := range arr {
		if v == a {
			return true
		}
	}
	return false
}
func MaskPhone(phone string) string {
	return maskPhoneReg.ReplaceAllString(phone, "$1****$3")
}
func IsWx(ua string) bool {
	return wxUaReg.MatchString(ua)
}
func IsPhone(phone string) bool {
	return phoneReg.MatchString(phone)
}
