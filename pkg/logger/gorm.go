package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	// color
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"

	// ctx key
	KeyRecordSqlFn = "recordSqlFn"

	// log level
	Silent gormLogger.LogLevel = iota + 1
	Error
	Warn
	Info
)

type (
	Writer interface {
		Printf(string, ...interface{})
	}
	Config struct {
		SlowThreshold time.Duration
		Colorful      bool
		LogLevel      gormLogger.LogLevel
	}
	SqlRecord struct {
		Sql        string    `json:"sql"`
		Level      string    `json:"level"`
		BeginTime  time.Time `json:"beginTime"`
		Source     string    `json:"source"`
		ExecTime   string    `json:"execTime"`
		AffectRows int64     `json:"affectRows"`
	}
	logger struct {
		Writer
		Config
		infoStr, warnStr, errStr            string
		traceStr, traceErrStr, traceWarnStr string
	}
)

var (
	GormLogger = NewGormLogger(log.New(os.Stdout, "\r\n", log.LstdFlags), Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      Info,
		Colorful:      true,
	})
)

func NewGormLogger(writer Writer, config Config) gormLogger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = Green + "%s\n" + Reset + Green + "[info] " + Reset
		warnStr = BlueBold + "%s\n" + Reset + Magenta + "[warn] " + Reset
		errStr = Magenta + "%s\n" + Reset + Red + "[error] " + Reset
		traceStr = Green + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
		traceWarnStr = Green + "%s " + Yellow + "%s\n" + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
		traceErrStr = RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	}
	return &logger{
		Writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// LogMode log mode
func (l *logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	l.LogLevel = level
	return l
}
func (l logger) Info(_ context.Context, _ string, _ ...interface{}) {

}
func (l logger) Warn(_ context.Context, _ string, _ ...interface{}) {

}
func (l logger) Error(_ context.Context, _ string, _ ...interface{}) {

}
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > Silent {
		elapsed := time.Since(begin)
		level := "info"
		source := utils.FileWithLineNum()
		execTime := float64(elapsed.Nanoseconds()) / 1e6
		sql, rows := fc()
		switch {
		case err != nil:
			level = "error"
			if rows == -1 {
				l.Printf(l.traceErrStr, source, err, execTime, "-", sql)
			} else {
				l.Printf(l.traceErrStr, source, err, execTime, rows, sql)
			}
			sql = "[error] : " + err.Error() + " [sql]: " + sql
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0:
			level = "warn"
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				l.Printf(l.traceWarnStr, source, slowLog, execTime, "-", sql)
			} else {
				l.Printf(l.traceWarnStr, source, slowLog, execTime, rows, sql)
			}
		default:

		}
		if f, ok := ctx.Value(KeyRecordSqlFn).(func(record SqlRecord)); ok {
			go f(SqlRecord{
				Sql:        sql,
				Level:      level,
				BeginTime:  begin,
				Source:     source,
				ExecTime:   elapsed.String(),
				AffectRows: rows,
			})
		}
	}
}
