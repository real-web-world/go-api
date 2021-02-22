package bootstrap

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	en2 "github.com/go-playground/locales/en"
	zh2 "github.com/go-playground/locales/zh"
	zhTw2 "github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	zhTwTrans "github.com/go-playground/validator/v10/translations/zh_tw"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/real-web-world/go-api/conf"
	"github.com/real-web-world/go-api/global"
	"github.com/real-web-world/go-api/middleware"
	"github.com/real-web-world/go-api/pkg/bdk"
	"github.com/real-web-world/go-api/pkg/gin"
	"github.com/real-web-world/go-api/pkg/logger"
	"github.com/real-web-world/go-api/pkg/valid"
)

func initConf() {
	_ = godotenv.Load(".env")
	if bdk.IsFile(".env.local") {
		_ = godotenv.Overload(".env.local")
	}
	err := configor.Load(global.Conf)
	if err != nil {
		panic(err)
	}
}

func initLog(cfg *conf.LogConf) {
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filepath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  true,
	})
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeDuration = zapcore.StringDurationEncoder
	level, err := logger.Str2ZapLevel(cfg.Level)
	if err != nil {
		panic("zap level is Incorrect")
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		ws,
		zap.NewAtomicLevelAt(level),
	)
	global.Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()

}
func initDB() {
	cfg := global.Conf.DB
	var err error
	global.DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.UserName, cfg.Pwd, cfg.Host, cfg.Port, cfg.Database, cfg.Charset),
		DefaultStringSize:         256,
		DisableDatetimePrecision:  false,
		DontSupportRenameIndex:    false,
		DontSupportRenameColumn:   false,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.Prefix,
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.GormLogger,
	})
	if err != nil {
		panic(errors.Wrap(err, "init db connect failed"))
	}
	sqlDB, _ := global.DB.DB()
	sqlDB.SetMaxIdleConns(cfg.MaxIDleConn)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifeTime * time.Minute)

}
func initCache() {
	cfg := global.Conf.Redis
	global.RedisPool = &redis.Pool{
		MaxIdle:     cfg.ConnPool,
		MaxActive:   1000,
		Wait:        true,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
				redis.DialDatabase(0), redis.DialPassword(cfg.Pwd))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// nolint gocritic
func logFormatter(p gin.LogFormatterParams) string {
	isDevApi := p.Request.RequestURI == global.PrometheusApi || strings.Index(p.Request.
		RequestURI, global.DebugApiPrefix) == 0
	if isDevApi || p.StatusCode == http.StatusNotFound {
		return ""
	}
	reqTime := p.TimeStamp.Format("2006-01-02 15:04:05")
	path := p.Request.URL.Path
	method := p.Request.Method
	code := p.StatusCode
	clientIp := p.ClientIP
	errMsg := p.ErrorMessage
	processTime := p.Latency
	return fmt.Sprintf("API: %s %d %s %s %s %v %s\n", reqTime, code, clientIp, path, method, processTime,
		errMsg)
}
func initTrans() {
	zh := zh2.New()
	zhTw := zhTw2.New()
	en := en2.New()
	Uni := ut.New(zh, zhTw, en)
	enT, _ := Uni.GetTranslator(global.LangEN)
	valid2 := binding.Validator.Engine().(*validator.Validate)
	_ = enTrans.RegisterDefaultTranslations(valid2, enT)
	zhT, _ := Uni.GetTranslator(global.LangZH)
	_ = zhTrans.RegisterDefaultTranslations(valid2, zhT)
	zhTwT, _ := Uni.GetTranslator(global.LangZHtw)
	_ = zhTwTrans.RegisterDefaultTranslations(valid2, zhTwT)
	global.TransMap[global.LangZH] = &zhT
	global.TransMap[global.LangEN] = &enT
	global.TransMap[global.LangZHtw] = &zhTwT
}
func initRegValidates() {
	for _, fn := range global.ValidFuncList {
		fn()
	}
}
func initEngine() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.LoggerWithFormatter(logFormatter))
	defaultValid := &valid.DefaultValidator{}
	binding.Validator = defaultValid
	initTrans()
	defaultValid.InitTrans(global.TransMap)
	initRegValidates()
	if global.Conf.PProf.Enable {
		pprof.Register(engine)
	}
	return engine
}
func initMiddleware(e *gin.Engine) {
	Conf := global.Conf
	// sentry
	e.Use(middleware.GenerateSentryMiddleware(Conf))
	e.Use(ginApp.PrepareProc)
	e.Use(middleware.Prometheus)
	// if Conf.ReqRateLimit.Enable {
	// 	lmt := tollbooth.NewLimiter(Conf.ReqRateLimit.Max, nil)
	// 	e.Use(middleware.RateLimit(lmt))
	// }
	// 解析header中的token
	e.Use(middleware.Auth)
	// e.Use(middleware.APIAuth)
	// trace
	e.Use(middleware.HTTPTrace)
	// csrf
	e.Use(middleware.Cors(Conf))
	// 绑定参数错误提示本地化
	e.Use(middleware.Locale(Conf))
	e.Use(middleware.Logger)
	// 报错时返回错误信息给客户端
	e.Use(middleware.ErrHandler)
}
func InitApp() *gin.Engine {
	initConf()
	initLog(&global.Conf.Log)
	initDB()
	initCache()
	e := initEngine()
	initMiddleware(e)
	gin.SetMode(global.Conf.Mode)
	return e

}

func InitMigration() {
	initConf()
	initDB()
}
