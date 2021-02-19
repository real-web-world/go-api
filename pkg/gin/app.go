package ginApp

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"

	"github.com/real-web-world/go-web-api/pkg/bdk"
	"github.com/real-web-world/go-web-api/pkg/dto/retcode"
	"github.com/real-web-world/go-web-api/pkg/fastcurd"
	"github.com/real-web-world/go-web-api/pkg/logger"
	"github.com/real-web-world/go-web-api/services/cache"
)

const (
	RootID       = 1
	BoolStrTrue  = "true"
	BoolStrFalse = "false"
	// level
	LevelSuper   UserLevel = "超级管理员"
	LevelAdmin   UserLevel = "管理员"
	LevelAuthor  UserLevel = "作者"
	LevelGeneral UserLevel = "普通用户"
	// gender
	GenderMan     Gender = "男"
	GenderWoman   Gender = "女"
	GenderUnknown Gender = "未知"
	// head field
	HeadLocale      = "locale"
	HeadToken       = "token"
	HeadUserAgent   = "User-Agent"
	HeadContentType = "Content-Type"
	// ctx key
	KeyTrans          = "trans"
	KeyApp            = "app"
	KeyAuthUser       = "authUser"
	KeyInitOnce       = "initOnce"
	KeyProcBeginTime  = "procBeginTime"
	KeyNotSaveResp    = "notSaveResp"
	KeyResp           = "resp"
	KeyReqID          = "reqID"
	KeyApiCacheKey    = "apiCacheKey"
	KeyApiCacheExpire = "apiCacheExpire"
	KeyStatusCode     = "statusCode"
	KeyProcTime       = "procTime"
	KeyRecordSqlFn    = logger.KeyRecordSqlFn
)

var (
	respBadReq       = &fastcurd.RetJSON{Code: retcode.BadReq, Msg: "bad request"}
	respNoChange     = &fastcurd.RetJSON{Code: retcode.DefaultError, Msg: "no change"}
	respNoAuth       = &fastcurd.RetJSON{Code: retcode.NoAuth, Msg: "no auth"}
	respNoLogin      = &fastcurd.RetJSON{Code: retcode.NoLogin, Msg: "please login"}
	respReqFrequency = &fastcurd.RetJSON{Code: retcode.RateLimitError, Msg: "The request frequency is too high"}
)

type (
	UserLevel string
	Gender    string
	AuthUser  struct {
		IsLogin bool
		UID     int
		Level   UserLevel
		Name    string
		User    interface{}
	}
	App struct {
		IsLogin   bool
		IsAdmin   bool
		IsAuthor  bool
		IsSuper   bool
		IsGeneral bool
		C         *gin.Context
		AuthUser  *AuthUser
		mu        sync.Mutex
		sqls      []logger.SqlRecord
	}
	ServerInfo struct {
		Timestamp int64 `json:"timestamp"`
	}
)

func GetApp(c *gin.Context) *App {
	initOnce, ok := c.Get(KeyInitOnce)
	if !ok {
		panic("ctx must set initOnce")
	}
	initOnce.(*sync.Once).Do(func() {
		c.Set(KeyApp, newApp(c))
	})
	app, _ := c.Get(KeyApp)
	return app.(*App)
}
func newApp(c *gin.Context) *App {
	app := &App{
		C:    c,
		sqls: make([]logger.SqlRecord, 0, 4),
	}
	app.setRecordSqlFn()
	return app
}
func (app *App) SetUser(u *AuthUser) {
	app.AuthUser = u
	if app.AuthUser.IsLogin {
		app.IsLogin = true
		switch app.AuthUser.Level {
		case LevelAdmin:
			app.IsAdmin = true
		case LevelSuper:
			app.IsAdmin = true
			app.IsSuper = true
		case LevelAuthor:
			app.IsAuthor = true
		case LevelGeneral:
			app.IsGeneral = true
		default:
		}
	}
}
func (app *App) GetSqls() []logger.SqlRecord {
	return app.sqls
}

// finally resp fn
func (app *App) Response(code int, json *fastcurd.RetJSON) {
	app.SetCtxRespVal(json)
	procBeginTime := app.GetProcBeginTime()
	reqID := app.GetReqID()
	var procTime string
	if procBeginTime != nil {
		procTime = time.Since(*procBeginTime).String()
	}
	json.Extra = &fastcurd.RespJsonExtra{
		ProcTime: procTime,
		ReqID:    reqID,
	}
	if gin.IsDebugging() {
		<-time.After(2 * time.Millisecond)
		json.Extra.SQLs = &app.sqls
	}
	// todo remove to project mid
	apiCacheKey := app.GetApiCacheKey()
	apiCacheExpire := app.GetApiCacheExpire()
	if apiCacheKey != nil && apiCacheExpire != nil {
		_ = cache.SetAPICache(*apiCacheKey, json, *apiCacheExpire)
	}
	app.SetStatusCode(code)
	app.SetProcTime(procTime)
	app.C.JSON(code, json)
	app.C.Abort()
}

// resp helper
func (app *App) Ok(msg string, data ...interface{}) {
	var actData interface{} = nil
	if len(data) == 1 {
		actData = data[0]
	}
	app.JSON(&fastcurd.RetJSON{Code: retcode.Ok, Msg: msg, Data: actData})
}
func (app *App) Data(data interface{}) {
	app.RetData(data)
}
func (app *App) ServerError(err error) {
	json := &fastcurd.RetJSON{Code: retcode.ServerError, Msg: err.Error()}
	app.Response(http.StatusInternalServerError, json)
}
func (app *App) RetData(data interface{}, msgParam ...string) {
	msg := ""
	if len(msgParam) == 1 {
		msg = msgParam[0]
	}
	app.Ok(msg, data)
}
func (app *App) JSON(json *fastcurd.RetJSON) {
	app.Response(http.StatusOK, json)
}
func (app *App) XML(xml interface{}) {
	procBeginTime := app.GetProcBeginTime()
	var procTime string
	if procBeginTime != nil {
		procTime = time.Since(*procBeginTime).String()
	}
	code := http.StatusOK
	app.SetStatusCode(code)
	app.SetProcTime(procTime)
	app.C.XML(code, xml)
	app.C.Abort()
}
func (app *App) SendList(list interface{}, count int) {
	app.Response(http.StatusOK, &fastcurd.RetJSON{
		Code:  retcode.Ok,
		Data:  list,
		Count: &count,
	})
}
func (app *App) BadReq() {
	app.Response(http.StatusBadRequest, respBadReq)
}
func (app *App) String(html string) {
	app.C.String(http.StatusOK, html)
}
func (app *App) ValidError(err error) {
	json := &fastcurd.RetJSON{}
	switch actErr := err.(type) {
	case validator.ValidationErrors:
		json.Code = retcode.ValidError
		json.Msg = actErr[0].Error()
		translator := app.GetTranslator()
		if translator != nil {
			json.Msg = actErr[0].Translate(*translator)
		}
	default:
		if err.Error() == "EOF" {
			json.Code = retcode.ValidError
			json.Msg = "request param is required"
		} else {
			json.Code = retcode.DefaultError
			json.Msg = actErr.Error()
		}
	}
	app.Response(http.StatusBadRequest, json)
}
func (app *App) NoChange() {
	app.JSON(respNoChange)
}
func (app *App) NoAuth() {
	app.Response(http.StatusUnauthorized, respNoAuth)
}
func (app *App) NoLogin() {
	app.Response(http.StatusUnauthorized, respNoLogin)
}
func (app *App) ErrorMsg(msg string) {
	json := &fastcurd.RetJSON{Code: retcode.DefaultError, Msg: msg}
	app.Response(http.StatusOK, json)
}
func (app *App) CommonError(err error) {
	app.CommonError(err)
}
func (app *App) RateLimitError() {
	app.Response(http.StatusBadRequest, respReqFrequency)
}
func (app *App) Success() {
	json := &fastcurd.RetJSON{Code: retcode.Ok}
	app.Response(http.StatusOK, json)
}
func (app *App) SendAffectRows(affectRows int) {
	app.Data(gin.H{
		"affectRows": affectRows,
	})
}

// utils fn
func (app *App) ISWx() bool {
	return bdk.IsWx(app.GetUA())
}
func (app *App) GetFullReqURL() string {
	schema := "http://"
	req := app.C.Request
	if req.TLS != nil {
		schema = "https://"
	}
	return schema + req.Host + req.RequestURI
}

// sql record
func (app *App) RecordSql(record logger.SqlRecord) {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.sqls = append(app.sqls, record)
}

// head field helper
func (app *App) GetLocale() string {
	return app.C.Request.Header.Get(HeadLocale)
}
func (app *App) GetUserAgent() string {
	return app.GetUA()
}
func (app *App) GetToken() string {
	return app.C.Request.Header.Get(HeadToken)
}
func (app *App) GetUA() string {
	return app.C.GetHeader(HeadUserAgent)
}
func (app *App) GetContentType() string {
	return app.C.GetHeader(HeadContentType)
}

// ctx value helper
func (app *App) setRecordSqlFn() {
	app.C.Set(KeyRecordSqlFn, app.RecordSql)
}
func (app *App) GetNotSaveResp() *bool {
	b, ok := app.C.Get(KeyNotSaveResp)
	if !ok {
		return nil
	}
	if t, ok := b.(bool); ok {
		return &t
	}
	return nil
}
func (app *App) SetNotSaveResp() {
	app.C.Set(KeyNotSaveResp, true)
}
func (app *App) IsShouldSaveResp() bool {
	t := app.GetNotSaveResp()
	return t == nil || *t
}
func (app *App) GetCtxAuthUser() (*AuthUser, bool) {
	u, _ := app.C.Get(KeyAuthUser)
	authUser, ok := u.(*AuthUser)
	return authUser, ok
}
func (app *App) SetCtxAuthUser(u *AuthUser) {
	app.C.Set(KeyAuthUser, u)
}
func (app *App) GetCtxRespVal() *fastcurd.RetJSON {
	if json, ok := app.C.Get(KeyResp); ok {
		return json.(*fastcurd.RetJSON)
	}
	return nil
}
func (app *App) SetCtxRespVal(json *fastcurd.RetJSON) {
	app.C.Set(KeyResp, json)
}
func (app *App) GetProcBeginTime() *time.Time {
	if procBeginTime, ok := app.C.Get(KeyProcBeginTime); ok {
		return procBeginTime.(*time.Time)
	}
	return nil
}
func (app *App) SetProcBeginTime(beginTime *time.Time) {
	app.C.Set(KeyProcBeginTime, beginTime)
}
func (app *App) GetReqID() string {
	if reqID, ok := app.C.Get(KeyReqID); ok {
		return reqID.(string)
	}
	return ""
}
func (app *App) SetReqID(reqID string) {
	app.C.Set(KeyReqID, reqID)
}
func (app *App) GetStatusCode() int {
	if code, ok := app.C.Get(KeyStatusCode); ok {
		return code.(int)
	}
	return 0
}
func (app *App) SetStatusCode(code int) {
	app.C.Set(KeyStatusCode, code)
}
func (app *App) GetProcTime() string {
	if procTime, ok := app.C.Get(KeyProcTime); ok {
		return procTime.(string)
	}
	return ""
}
func (app *App) SetProcTime(procTime string) {
	app.C.Set(KeyProcTime, procTime)
}
func (app *App) GetApiCacheKey() *string {
	if key, ok := app.C.Get(KeyApiCacheKey); ok {
		fmtKey := key.(string)
		return &fmtKey
	}
	return nil
}
func (app *App) SetApiCacheKey(key string) {
	app.C.Set(KeyApiCacheKey, key)
}
func (app *App) GetApiCacheExpire() *time.Duration {
	if ex, ok := app.C.Get(KeyApiCacheExpire); ok {
		return ex.(*time.Duration)
	}
	return nil
}
func (app *App) SetApiCacheExpire(d *time.Duration) {
	app.C.Set(KeyApiCacheExpire, d)
}
func (app *App) GetTranslator() *ut.Translator {
	if translator, ok := app.C.Get(KeyTrans); ok {
		return translator.(*ut.Translator)
	}
	return nil
}
func (app *App) SetTranslator(translator *ut.Translator) {
	app.C.Set(KeyTrans, translator)
}

// middleware
func PrepareProc(c *gin.Context) {
	now := time.Now()
	c.Set(KeyInitOnce, &sync.Once{})
	c.Set(KeyProcBeginTime, &now)
	c.Set(KeyReqID, uuid.NewV4().String())
	c.Next()
}
