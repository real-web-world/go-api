package api

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/real-web-world/go-web-api/models"
	"github.com/real-web-world/go-web-api/pkg/bdk"
	"github.com/real-web-world/go-web-api/pkg/dto/retcode"
	"github.com/real-web-world/go-web-api/pkg/fastcurd"
	"github.com/real-web-world/go-web-api/pkg/gin"
	"github.com/real-web-world/go-web-api/services/cache"
	auth "github.com/real-web-world/go-web-api/services/user"
)

type LoginData struct {
	Account       string         `json:"account" binding:"required"`
	Pwd           string         `json:"pwd" binding:"required"`
	VerifyCode    string         `json:"verifyCode" binding:"required"`
	CapID         string         `json:"capID" binding:""`
	Source        *models.Source `json:"source" binding:"omitempty,oneof=android ios web"`
	AppUniqueID   *string        `json:"appUniqueID"`
	AppDeviceInfo *string        `json:"appDeviceInfo"`
}
type LoginResData struct {
	Token  string      `json:"token"`
	Expire int64       `json:"expire"`
	Detail interface{} `json:"detail"`
}

// for swag
type (
	_ fastcurd.RetJSON
	_ models.DefaultSceneUser
)

func validError(app *ginApp.App, args ...string) {
	msg := "账号或密码错误"
	if len(args) == 1 {
		msg = args[0]
	}
	app.ErrorMsg(msg)
}

// @summary 用户登录
// @tags user
// @Accept  json
// @Produce  json
// @param data body LoginData true "123"
// @Success 200 {object} fastcurd.RetJSON{data=LoginResData{detail=models.DefaultSceneUser}}	"token token过期时间 用户详情"
// @router /user/login [post]
func Login(c *gin.Context) {
	app := ginApp.GetApp(c)
	d := &LoginData{}
	webUserAgent := app.GetUserAgent()
	clientIP := c.ClientIP()
	loginHistory := &models.AddLoginHistoryData{
		OK:           ginApp.BoolStrFalse,
		ClientIP:     clientIP,
		Account:      &d.Account,
		WebUserAgent: &webUserAgent,
	}
	if err := c.ShouldBindJSON(d); err != nil {
		app.ValidError(err)
		go func() { _ = models.AddLoginHistory(loginHistory) }()
		return
	}
	if d.Source != nil {
		loginHistory.Source = *d.Source
	}
	if d.AppUniqueID != nil {
		loginHistory.AppUniqueID = d.AppUniqueID
	}
	if d.AppDeviceInfo != nil {
		loginHistory.AppDeviceInfo = d.AppDeviceInfo
	}
	if !cache.VerifyCaptcha(clientIP, d.CapID, d.VerifyCode) {
		validError(app, "验证码不正确")
		go func() { _ = models.AddLoginHistory(loginHistory) }()
		return
	}
	user := models.NewCtxUser(c)
	if err := models.Get(user, "account", d.Account); err != nil {
		validError(app)
		go func() { _ = models.AddLoginHistory(loginHistory) }()
		return
	}
	// 验证argon2id
	if !auth.ValidUserPwd(d.Pwd, user.Pwd) {
		validError(app)
		go func() { _ = models.AddLoginHistory(loginHistory) }()
		return
	}
	token, expire, _ := auth.Login(user.ID)
	json := &fastcurd.RetJSON{
		Code: retcode.Ok,
		Data: LoginResData{
			Token:  token,
			Expire: time.Now().Add(time.Duration(expire) * time.Second).Unix(),
			Detail: user.GetFmtDetail(),
		},
	}
	app.JSON(json)
	loginHistory.OK = ginApp.BoolStrTrue
	go func() { _ = models.AddLoginHistory(loginHistory) }()
}
func Logout(c *gin.Context) {
	app := ginApp.GetApp(c)
	if app.IsLogin {
		go func() {
			token := app.GetToken()
			_ = auth.Logout(token)
			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()
		}()
	}
	app.Success()
}
func ModifyPwd(c *gin.Context) {
	type data struct {
		UID    *int   `json:"uid" binding:"omitempty"`
		OriPwd string `json:"oriPwd" binding:"required"`
		// 新密码不能与原密码相同
		CurrPwd string `json:"currPwd" binding:"required,min=4,max=32,necsfield=OriPwd"`
		RePwd   string `json:"rePwd" binding:"required,eqfield=CurrPwd"`
	}
	app := ginApp.GetApp(c)
	u := app.AuthUser.User.(*models.User)
	d := &data{}
	if err := c.ShouldBindJSON(d); err != nil {
		app.ValidError(err)
		return
	}
	actUser := models.NewCtxUser(c)
	if d.UID != nil && *d.UID != actUser.ID {
		if !app.IsSuper {
			app.ErrorMsg("只有超级管理员才可以修改其他用户的密码")
			return
		}
		if !models.IsExist(actUser, map[string]interface{}{
			fastcurd.PrimaryField: d.UID,
		}) {
			app.ErrorMsg("用户不存在")
			return
		}
	} else {
		d.UID = new(int)
		*d.UID = u.ID
		if !auth.ValidUserPwd(d.OriPwd, u.Pwd) {
			app.ErrorMsg("原密码不正确")
			return
		}
	}
	hashPwd, err := auth.CreateHashPwd(d.CurrPwd)
	if err != nil {
		app.ErrorMsg("生成密码失败 " + err.Error())
		return
	}
	actUser.Pwd = hashPwd
	if _, err := models.Update(actUser, d.UID); err != nil {
		app.CommonError(err)
		return
	}
	app.Success()
}

func AddUser(c *gin.Context) {
	app := ginApp.GetApp(c)
	d := &models.AddUserData{}
	if err := c.ShouldBindJSON(d); err != nil {
		app.ValidError(err)
		return
	}
	if d.Level == ginApp.LevelAdmin && !app.IsSuper {
		app.ErrorMsg("只有超级管理员才可以添加管理员")
		return
	}
	user := models.NewCtxUser(c, &models.User{
		Account:  d.Account,
		Level:    d.Level,
		Name:     d.Name,
		Gender:   d.Gender,
		Age:      d.Age,
		Phone:    d.Phone,
		AvatarID: d.AvatarID,
		Profile:  d.Profile,
	})
	if models.IsExist(user, map[string]interface{}{
		"account": user.Account,
	}) {
		app.ErrorMsg("账号" + user.Account + "已存在")
		return
	}
	hashPwd, err := auth.CreateHashPwd(d.Pwd)
	if err != nil {
		app.ErrorMsg("生成密码失败")
		return
	}
	user.Pwd = hashPwd
	if err := models.Add(user); err != nil {
		app.CommonError(err)
		return
	}
	app.Success()
}
func EditUser(c *gin.Context) {
	app := ginApp.GetApp(c)
	d := &models.EditUserData{}
	if err := c.ShouldBindJSON(d); err != nil {
		app.ValidError(err)
		return
	}
	user := models.NewCtxUser(c)
	user.Gender = d.Gender
	user.Name = d.Name
	user.Age = d.Age
	if !app.IsLogin {
		app.NoAuth()
		return
	}
	if app.IsAdmin && d.ID != nil && *d.ID != 0 {
		user.ID = *d.ID
	} else {
		user.ID = app.AuthUser.UID
		d.ID = new(int)
		*d.ID = user.ID
	}
	if user.ID == ginApp.RootID && !app.IsSuper {
		app.ErrorMsg("can't modify super master")
		return
	}
	if app.IsSuper {
		user.Level = d.Level
	}
	if app.IsAdmin {
		user.Phone = d.Phone
	}
	// todo 应该判断图片id必须是此用户上传的
	user.AvatarID = d.AvatarID
	user.Profile = d.Profile
	affectRows, err := models.Update(user, user.ID)
	if err != nil {
		app.CommonError(err)
		return
	}
	if affectRows+user.RelationAffectRows == 0 {
		app.NoChange()
		return
	}
	app.Success()
}
func DelUser(c *gin.Context) {
	app := ginApp.GetApp(c)
	d := &fastcurd.DelData{}
	if err := c.ShouldBindJSON(d); err != nil {
		app.ValidError(err)
		return
	}
	if bdk.InArrayInt(ginApp.RootID, d.IDs) {
		app.ErrorMsg("can't delete root")
		return
	}
	affectRows, err := models.Delete(models.NewCtxUser(c), d.IDs)
	if err != nil {
		app.CommonError(err)
		return
	}
	app.SendAffectRows(affectRows)
}
func UserDetail(c *gin.Context) {
	app := ginApp.GetApp(c)
	d := &fastcurd.NullableDetailData{}
	if c.Request.ContentLength != 0 {
		if err := c.ShouldBindJSON(d); err != nil {
			app.ValidError(err)
			return
		}
	}
	var err error
	var user *models.User
	if d.ID == nil || *d.ID == 0 || app.IsAuthor {
		var ok bool
		if user, ok = app.AuthUser.User.(*models.User); !ok {
			log.Printf("%+v", user)
			app.NoLogin()
			return
		}
	} else if app.IsAdmin {
		user = models.NewCtxUser(c)
		if err = models.Detail(user, d.ID); err != nil {
			app.CommonError(err)
			return
		}
	} else {
		log.Println(316)
		app.NoLogin()
		return
	}
	scene := d.GetScene()
	if scene == models.SceneAdmin && !app.IsAdmin && !app.IsAuthor {
		app.NoAuth()
		return
	}
	app.Data(user.GetFmtDetail(scene))
}
func ListUser(c *gin.Context) {
	app := ginApp.GetApp(c)
	m := models.NewCtxUser(c)
	d := &fastcurd.ListData{}
	if err := c.ShouldBindJSON(d); err != nil {
		app.ValidError(err)
		return
	}
	if d.Order == nil {
		d.Order = make(map[string]string)
	}
	_, hasSortOrderField := d.Order["sort"]
	if m.HasSortField() && !hasSortOrderField {
		d.Order["sort"] = "desc"
	}
	if d.Filter == nil {
		d.Filter = fastcurd.Filter{}
	}
	if searchFilter, ok := d.Filter["search"]; ok {
		if phone, ok := searchFilter.Val.(string); ok {
			if bdk.IsPhone(phone) {
				delete(d.Filter, "search")
				user := models.NewCtxUser(c)
				_ = models.Get(user, "phone", phone)
				if user.ID != 0 {
					d.Filter["id"] = struct {
						Condition fastcurd.FilterCond
						Val       interface{}
					}{Condition: fastcurd.CondEq, Val: user.ID}
				}
			}
		}
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
