package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/real-web-world/go-api/conf"
)

const (
	LangZH         = "zh"
	LangEN         = "en"
	LangZHtw       = "zh_tw"
	PrometheusApi  = "/prom"
	DebugApiPrefix = "/debug"
)

var (
	DB            *gorm.DB
	RedisPool     *redis.Pool
	Conf          = &conf.AppConf{}
	TransMap      = make(map[string]*ut.Translator)
	ValidFuncList []func()
	Logger        *zap.SugaredLogger
)

func GetTrans(localeParam ...string) *ut.Translator {
	locale := Conf.Lang.DefaultLang
	if len(localeParam) > 0 {
		locale = localeParam[0]
	}
	var translator *ut.Translator
	var ok bool
	if translator, ok = TransMap[locale]; !ok {
		translator = TransMap[Conf.Lang.DefaultLang]
	}
	return translator
}
