package cache

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/real-web-world/go-web-api/global"
	"github.com/real-web-world/go-web-api/pkg/bdk"
	"github.com/real-web-world/go-web-api/pkg/fastcurd"
)

var (
	apiPrefix                   = "apiCache"
	captchaPrefix               = "captcha"
	tokenPrefix                 = "token"
	commonPhoneVerifyCodePrefix = "commonPhoneVerifyCode"
	captchaExpire               = 15 * 60
	commonPhoneVerifyCodeExpire = 15 * 60
)

func fmtKey(key string) string {
	return global.Conf.Redis.CollectionName + ":" + key
}
func fmtCaptchaKey(reqIP, capID string) string {
	return fmtKey(captchaPrefix + ":" + reqIP + "-" + capID)
}
func fmtApiKey(key string) string {
	return fmtKey(apiPrefix + ":" + key)
}
func fmtCommonPhoneVerifyCodeKey(reqIP, uuid, phone string) string {
	return fmtKey(commonPhoneVerifyCodePrefix + ":" + reqIP + "-" + uuid + phone)
}
func fmtTokenKey(token string) string {
	return fmtKey(tokenPrefix + ":" + token)
}
func SetCaptcha(reqIP, capID, code string) error {
	code = strings.ToLower(code)
	ce := global.RedisPool.Get()
	defer ce.Close()
	_, err := ce.Do("setex", fmtCaptchaKey(reqIP, capID), captchaExpire, code)
	if err != nil {
		return err
	}
	return nil
}
func VerifyCaptcha(reqIP, capID, code string) bool {
	// todo dev
	if code == "4396" {
		return true
	}
	code = strings.ToLower(code)
	key := fmtCaptchaKey(reqIP, capID)
	ce := global.RedisPool.Get()
	defer ce.Close()
	val, err := redis.String(ce.Do("get", key))
	if err != nil {
		return false
	}
	ok := val == code
	if ok {
		go func() { _, _ = ce.Do("del", key) }()
	}
	return ok
}
func SetAPICache(key string, resp *fastcurd.RetJSON, expireParam ...time.Duration) error {
	ce := global.RedisPool.Get()
	defer ce.Close()
	respBts, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	expire := 0
	if len(expireParam) == 1 {
		expire = int(expireParam[0] / time.Second)
	}
	if key == "" {
		return errors.New("缓存键不能为空")
	}
	// todo review string() to bts2str
	_, err = ce.Do("setex", fmtApiKey(key),
		expire, bdk.Bytes2Str(respBts))
	return err
}
func GetAPICache(key string) *fastcurd.RetJSON {
	ce := global.RedisPool.Get()
	defer ce.Close()
	var resp fastcurd.RetJSON
	jsonStr, err := redis.String(ce.Do("get", fmtApiKey(key)))
	if err != nil {
		return nil
	}
	if json.Unmarshal([]byte(jsonStr), &resp) != nil {
		return nil
	}
	return &resp
}
func SetPhoneCommonVerifyCode(phone, uuid, reqIP string, code int) error {
	ce := global.RedisPool.Get()
	defer ce.Close()
	key := fmtCommonPhoneVerifyCodeKey(reqIP, uuid, phone)
	_, err := ce.Do("setex", key, commonPhoneVerifyCodeExpire, code)
	if err != nil {
		return err
	}
	return nil
}
func VerifyPhoneCommonVerifyCode(phone, uuid, reqIP, code string) bool {
	code = strings.ToLower(code)
	key := fmtCommonPhoneVerifyCodeKey(reqIP, uuid, phone)
	ce := global.RedisPool.Get()
	defer ce.Close()
	val, err := redis.String(ce.Do("get", key))
	if err != nil {
		return false
	}
	ok := val == code
	if ok {
		go func() { _, _ = ce.Do("del", key) }()
	}
	return ok
}
func GetUIDByToken(token string) (uid int, err error) {
	ce := global.RedisPool.Get()
	defer ce.Close()
	return redis.Int(ce.Do("get", fmtTokenKey(token)))
}
func SetToken(token string, expire, uid int) error {
	ce := global.RedisPool.Get()
	defer ce.Close()
	_, err := ce.Do("setex", fmtTokenKey(token), expire, uid)
	return err
}
func DelToken(token string) error {
	ce := global.RedisPool.Get()
	defer ce.Close()
	_, err := ce.Do("del", fmtTokenKey(token))
	return err
}
func UpdateTokenExpire(token string, expire int) error {
	ce := global.RedisPool.Get()
	defer ce.Close()
	_, err := ce.Do("expire", fmtTokenKey(token), expire)
	return err
}
