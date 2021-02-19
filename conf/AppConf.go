package conf

import (
	"time"

	"github.com/real-web-world/go-web-api/pkg/logger"
)

type AppConf struct {
	HTTPHost string `default:"0.0.0.0"`
	HTTPPort int    `default:"8898"`
	Mode     string `default:"release" env:"mode"`
	Host     string `required:"true" env:"webHost"`
	DB       DBConf
	Redis    RedisConf
	Sentry   SentryConf
	Csrf     CsrfConf
	Cdn      CdnConf
	Oss      OssConf
	Lang     LangConf
	Token    TokenConf
	PProf    PProfConf `json:"pprof"`
	Log      LogConf
}
type TokenConf struct {
	// token 过期时间 1440分钟 一天
	Expire int `default:"1440" env:"tokenExpire"`
}
type DBConf struct {
	Type     string `default:"mysql"`
	Host     string `default:"127.0.0.1" env:"dbHost"`
	Port     int    `default:"3306" env:"dbPort"`
	UserName string `required:"true" env:"dbUserName"`
	Pwd      string `required:"true" env:"dbPwd"`
	Charset  string `default:"utf8mb4"`
	Database string `required:"true" env:"dbDatabase"`
	// 表前缀
	Prefix string `required:"true" env:"dbTableNamePrefix"`
	// 最大闲置连接
	MaxIDleConn int `default:"10"`
	// 最大打开连接
	MaxOpenConn int `default:"100"`
	// 一个连接的最大生命时长(分钟)
	MaxLifeTime time.Duration `default:"60"`
}
type RedisConf struct {
	Host           string `default:"127.0.0.1" env:"redisHost"`
	Port           int    `default:"6379" env:"redisPort"`
	CollectionName string `env:"redisCollectionName"`
	Pwd            string `env:"redisPwd"`
	ConnPool       int    `default:"50"`
}
type SentryConf struct {
	Dsn string
}
type CsrfConf struct {
	AllowOrigins string `env:"allowOrigins"`
}
type CdnConf struct {
	Host string `env:"cdnHost"`
}
type OssConf struct {
	EndPoint         string `default:"oss-cn-shanghai.aliyuncs.com" env:"ossEndPoint"`
	InternalEndPoint string `default:"oss-cn-shanghai-internal.aliyuncs.com"`
	WebRootPath      string
	Bucket           string `required:"true" env:"ossBucket"`
	AccessKeyID      string `required:"true" env:"ossAccessKeyID"`
	AccessKeySecret  string `required:"true" env:"ossAccessKeySecret"`
	Host             string `required:"true" env:"ossHost"`
	CallbackURL      string `default:"/notify/oss" env:"ossCallbackURL"`
	UploadDir        string `default:"uploads/" env:"ossUploadDir"`
	TokenExpire      int64  `default:"120" env:"ossTokenExpire"`
}
type LangConf struct {
	DefaultLang string `default:"zh" env:"defaultLang"`
}
type PProfConf struct {
	Enable bool `default:"false" env:"enablePProf"`
}
type LogConf struct {
	Level      logger.LogLevelStr `default:"info" env:"logLevel"`
	Filepath   string             `default:"logs/main.log" env:"logFilepath"`
	MaxSize    int                `default:"1024" env:"logMaxSize"` // log file max size unit is mb
	MaxBackups int                `default:"7" env:"logMaxBackups"`
	MaxAge     int                `default:"7" env:"logMaxAge"`
	Compress   bool               `default:"true" env:"logCompress"`
}
