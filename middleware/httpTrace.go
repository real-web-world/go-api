package middleware

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"

	"github.com/real-web-world/go-web-api/global"
	"github.com/real-web-world/go-web-api/pkg/bdk"
	"github.com/real-web-world/go-web-api/pkg/fastcurd"
	"github.com/real-web-world/go-web-api/pkg/gin"
)

const (
	MaxTraceSize = 1 << 10
)

type FileInfo struct {
	Size   int64                `json:"size"`
	Name   string               `json:"name"`
	Header textproto.MIMEHeader `json:"header"`
}
type Trace struct {
	ReqID      string                 `json:"reqID"`
	URL        string                 `json:"url"`
	Method     string                 `json:"method"`
	ClientIP   string                 `json:"clientIP"`
	Host       string                 `json:"host"`
	Header     http.Header            `json:"header"`
	FormData   map[string][]string    `json:"formData"`
	FileData   map[string][]*FileInfo `json:"fileData"`
	PostData   string                 `json:"postData"`
	Resp       *fastcurd.RetJSON      `json:"resp"`
	ProcTime   string                 `json:"procTime"`
	ReqTime    time.Time              `json:"reqTime"`
	StatusCode int                    `json:"statusCode"`
}

func HTTPTrace(c *gin.Context) {
	app := ginApp.GetApp(c)
	reqID := app.GetReqID()
	begin := app.GetProcBeginTime()
	trace := &Trace{
		ReqTime:  *begin,
		Method:   c.Request.Method,
		ClientIP: c.ClientIP(),
		ReqID:    reqID,
		URL:      c.Request.RequestURI,
		Host:     c.Request.Host,
		Header:   c.Request.Header,
	}
	contentType := strings.Split(app.GetContentType(), ";")[0]
	switch contentType {
	case "multipart/form-data":
		var err error
		var multiFormData *multipart.Form
		if multiFormData, err = c.MultipartForm(); err == nil && multiFormData != nil {
			trace.FormData = multiFormData.Value
			trace.FileData = make(map[string][]*FileInfo)
			for fileName, filesInfo := range multiFormData.File {
				for _, fileInfo := range filesInfo {
					trace.FileData[fileName] = append(trace.FileData[fileName], &FileInfo{
						Size:   fileInfo.Size,
						Name:   fileInfo.Filename,
						Header: fileInfo.Header,
					})
				}
			}
		}
	case "application/json":
		// todo bts sync.pool
		bodyBts, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBts))
		trace.PostData = bdk.Bytes2Str(bodyBts)
	}
	c.Next()
	procTime := app.GetProcTime()
	trace.ProcTime = procTime
	if app.IsShouldSaveResp() {
		trace.Resp = app.GetCtxRespVal()
	}
	statusCode := app.GetStatusCode()
	trace.StatusCode = statusCode
	switch statusCode {
	case http.StatusForbidden, http.StatusNotFound:
		return
	default:
	}
	// todo 判断 header formData postData
	header := trace.Header
	formData := trace.FormData
	postData := trace.PostData
	response := trace.Resp
	if unsafe.Sizeof(header) > MaxTraceSize {
		header = nil
	}
	if unsafe.Sizeof(formData) > MaxTraceSize {
		formData = nil
	}
	if len(postData) > MaxTraceSize {
		header = nil
	}
	if unsafe.Sizeof(*response) > MaxTraceSize {
		response = nil
	}
	global.Logger.Infow("httpTrace",
		"reqID", trace.ReqID,
		"url", trace.URL,
		"clientIP", trace.ClientIP,
		"statusCode", trace.StatusCode,
		"fileData", trace.FileData,
		"header", header,
		"formData", formData,
		"host", trace.Host,
		"method", trace.Method,
		"procTime", trace.ProcTime,
		"reqTime", trace.ReqTime,
		"postData", postData,
		"resp", response,
	)
}
