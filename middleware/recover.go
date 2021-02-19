package middleware

import (
	"fmt"
	"log"
	"runtime/debug"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/real-web-world/go-web-api/pkg/gin"
)

func ErrHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			app := ginApp.GetApp(c)
			reqID := app.GetReqID()
			var actErr error
			log.Println("发生异常: ", err)
			log.Println("reqID: ", reqID)
			debug.PrintStack()
			if app.IsLogin {
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetUser(sentry.User{
						ID:        strconv.Itoa(app.AuthUser.UID),
						IPAddress: c.ClientIP(),
						Username:  app.AuthUser.Name,
					})
					scope.SetExtra("reqID", reqID)
				})
			}
			sentry.CaptureException(errors.New(fmt.Sprintf("%v", err)))
			switch err := err.(type) {
			case error:
				actErr = err
			case string:
				errMsg := err
				actErr = errors.New(errMsg)
			default:
				actErr = errors.New("server exception")
			}
			app.ServerError(actErr)
			return
		}
	}()
	c.Next()
}
