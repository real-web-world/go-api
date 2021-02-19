package user

import (
	"github.com/alexedwards/argon2id"

	"github.com/real-web-world/go-web-api/global"
	"github.com/real-web-world/go-web-api/pkg/bdk"
	"github.com/real-web-world/go-web-api/services/cache"
)

func ValidUserPwd(inputPwd, hashPwd string) bool {
	match, _ := argon2id.ComparePasswordAndHash(inputPwd, hashPwd)
	return match
}
func CreateHashPwd(pwd string) (string, error) {
	return argon2id.CreateHash(pwd, argon2id.DefaultParams)
}
func Login(uid int) (token string, expire int, err error) {
	expire = global.Conf.Token.Expire * 60
	token = bdk.Bytes2Str(bdk.RandomAlphaNum(32))
	err = cache.SetToken(token, expire, uid)
	if err != nil {
		return "", 0, err
	}
	return
}
func Logout(token string) error {
	return cache.DelToken(token)
}
