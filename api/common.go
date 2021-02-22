package api

import (
	"github.com/real-web-world/go-api/assets"
	"image/png"

	"github.com/afocus/captcha"
	"github.com/gin-gonic/gin"

	"github.com/real-web-world/go-api"
	"github.com/real-web-world/go-api/models"
	"github.com/real-web-world/go-api/pkg/dto/retcode"
	"github.com/real-web-world/go-api/pkg/fastcurd"
	"github.com/real-web-world/go-api/pkg/gin"
	"github.com/real-web-world/go-api/services/cache"
)

type VersionInfo struct {
	Version string `json:"version" example:"1.0.0"`
	Commit  string `json:"commit" example:"github commit sha256"`
}

// @summary 测试接口 用于返回一些临时数据
// @tags test
// @router /test [post]
func TestHand(c *gin.Context) {
	app := ginApp.GetApp(c)
	// <-time.After(20 * time.Second)
	json := &fastcurd.RetJSON{
		Code: 0,
		Data: gin.H{
			"user": 123,
			"ip":   c.Request.RemoteAddr,
		},
	}
	app.JSON(json)
}

// @summary 获取当前api版本
// @Produce  json
// @success 200 {object} fastcurd.RetJSON{data=VersionInfo} "版本信息"
// @router /version [get]
func ShowVersion(c *gin.Context) {
	app := ginApp.GetApp(c)
	json := &fastcurd.RetJSON{
		Code: retcode.Ok,
		Data: &VersionInfo{
			Version: gowebapi.APIVersion,
			Commit:  gowebapi.Commit,
		},
	}
	app.JSON(json)
}

// @summary 获取验证码
// @description 获取后端登陆验证码
// @Produce  png
// @success 200 {object} object
// @router /getVerifyCode [get]
func GetCaptcha(c *gin.Context) {
	capWidth := 200
	capHeight := 60
	captchaCode := captcha.New()
	// 设置字体
	if err := captchaCode.AddFontFromBytes(assets.DefaultFont); err != nil {
		panic(err.Error())
	}
	captchaCode.SetSize(capWidth, capHeight)
	img, code := captchaCode.Create(4, captcha.CLEAR)
	if err := png.Encode(c.Writer, img); err != nil {
		panic(err.Error())
	}
	reqIP := c.ClientIP()
	capID, _ := c.GetQuery("capID")
	_ = cache.SetCaptcha(reqIP, capID, code)
}

func UserCommonList(m models.BaseModel) gin.HandlerFunc {
	return func(c *gin.Context) {
		app := ginApp.GetApp(c)
		m = m.NewModel()
		m.SetCtx(c)
		d := &fastcurd.ListData{}
		if err := c.ShouldBindJSON(d); err != nil {
			app.ValidError(err)
			return
		}
		scene := d.GetScene()
		if scene == models.SceneAdmin && !app.IsAdmin {
			app.NoAuth()
			return
		}
		if !app.IsAdmin {
			if d.Filter == nil {
				d.Filter = make(fastcurd.Filter)
			}
			d.Filter["uid"] = struct {
				Condition fastcurd.FilterCond
				Val       interface{}
			}{Condition: fastcurd.CondEq, Val: app.AuthUser.UID}
		}
		var count int
		var err error
		var list interface{}
		if count, list, err = models.List(m, d); err != nil {
			app.CommonError(err)
			return
		}
		m.BeforeFmtList(list, app)
		fmtList := m.GetFmtList(list, scene)
		app.SendList(fmtList, count)
	}
}

func UserCommonDetail(m models.BaseModel) gin.HandlerFunc {
	return func(c *gin.Context) {
		app := ginApp.GetApp(c)
		m = m.NewModel()
		m.SetCtx(c)
		type postData struct {
			ID  int `json:"id" binding:"required"`
			UID int `json:"uid"`
		}
		d := &postData{}
		if err := c.ShouldBindJSON(d); err != nil {
			app.ValidError(err)
			return
		}
		filter := fastcurd.Filter{
			"id": {
				Condition: fastcurd.CondEq,
				Val:       d.ID,
			},
		}
		if !app.IsAdmin {
			filter["uid"] = struct {
				Condition fastcurd.FilterCond
				Val       interface{}
			}{Condition: fastcurd.CondEq, Val: app.AuthUser.UID}
		}
		err := models.Detail(m, filter)
		if err != nil {
			app.CommonError(err)
			return
		}
		m.BeforeFmtDetail(app)
		go m.AfterDetailAction(app)
		app.Data(m.GetFmtDetail())
	}
}
func CommonDetail(m models.BaseModel) gin.HandlerFunc {
	return func(c *gin.Context) {
		m = m.NewModel()
		m.SetCtx(c)
		app := ginApp.GetApp(c)
		d := &fastcurd.DetailData{}
		if err := c.ShouldBindJSON(d); err != nil {
			app.ValidError(err)
			return
		}
		scene := d.GetScene()
		if scene == models.SceneAdmin && !app.IsAdmin {
			app.NoAuth()
			return
		}
		if err := m.GetDetail(d.ID); err != nil {
			app.CommonError(err)
			return
		}
		m.BeforeFmtDetail(app)
		go m.AfterDetailAction(app)
		app.Data(m.GetFmtDetail(scene))
	}
}
func CommonList(m models.BaseModel) gin.HandlerFunc {
	return func(c *gin.Context) {
		m = m.NewModel()
		m.SetCtx(c)
		app := ginApp.GetApp(c)
		d := &fastcurd.ListData{}
		if d.Order == nil {
			d.Order = make(map[string]string)
		}
		_, hasSortOrderField := d.Order["sort"]
		if m.HasSortField() && !hasSortOrderField {
			d.Order["sort"] = "desc"
		}
		if err := c.ShouldBindJSON(d); err != nil {
			app.ValidError(err)
			return
		}
		scene := d.GetScene()
		if scene == models.SceneAdmin && !app.IsAdmin {
			app.NoAuth()
			return
		}
		var count int
		var err error
		var list interface{}
		if count, list, err = models.List(m, d); err != nil {
			app.CommonError(err)
			return
		}
		m.BeforeFmtList(list, app)
		fmtList := m.GetFmtList(list, scene)
		app.SendList(fmtList, count)
	}
}
func CommonAdd(m models.BaseModel) gin.HandlerFunc {
	return func(c *gin.Context) {
		m = m.NewModel()
		m.SetCtx(c)
		app := ginApp.GetApp(c)
		d := m.NewAddData()
		if err := c.ShouldBindJSON(d); err != nil {
			app.ValidError(err)
			return
		}
		if err := m.AddWithPostData(d, app); err != nil {
			app.CommonError(err)
			return
		}
		app.Success()
	}
}
func CommonEdit(m models.BaseModel) gin.HandlerFunc {
	return func(c *gin.Context) {
		m = m.NewModel()
		m.SetCtx(c)
		app := ginApp.GetApp(c)
		d := m.NewEditData()
		if err := c.ShouldBindJSON(d); err != nil {
			app.ValidError(err)
			return
		}
		affectRows, err := m.EditWithPostData(d, app)
		if err != nil {
			app.CommonError(err)
			return
		}
		if affectRows != 1 && m.GetRelationAffectRows() == 0 {
			app.NoChange()
			return
		}
		go func() {
			app.C.Set("postData", d)
			m.AfterEditAction(app)
		}()
		app.Success()
	}
}
func CommonDel(m models.BaseModel) gin.HandlerFunc {
	return func(c *gin.Context) {
		m = m.NewModel()
		m.SetCtx(c)
		app := ginApp.GetApp(c)
		d := &fastcurd.DelData{}
		if err := c.ShouldBindJSON(d); err != nil {
			app.ValidError(err)
			return
		}
		affectRows, err := m.DelWithPostData(d, app)
		if err != nil {
			app.CommonError(err)
			return
		}
		json := &fastcurd.RetJSON{
			Code: retcode.Ok,
			Data: gin.H{
				"affectRows": affectRows,
			},
		}
		go func() {
			app.C.Set("postData", d)
			m.AfterDelAction(app)
		}()
		app.JSON(json)
	}
}
