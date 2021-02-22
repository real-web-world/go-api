package models

import (
	"time"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"github.com/real-web-world/go-api/global"
	"github.com/real-web-world/go-api/pkg/bdk"
	"github.com/real-web-world/go-api/pkg/gin"
	"github.com/real-web-world/go-api/pkg/valid"
)

func init() {
	global.ValidFuncList = append(global.ValidFuncList, regUserValid)
}

type UserStatus string

const (
	UserSceneProfile            = "profile"
	UserSceneAuthor             = "author"
	UserSceneExport             = "export"
	UserStatusNormal UserStatus = "正常"
	UserStatusFrozen UserStatus = "冻结"
)

var (
	UserFilterNameMapDBField = map[string]string{
		"search":       "account|name|id",
		"id":           "id",
		"account":      "account",
		"ctime":        "ctime",
		"utime":        "utime",
		"level":        "level",
		"phone":        "phone",
		"gender":       "gender",
		"name":         "name",
		"age":          "age",
		"provinceCode": "province_code",
		"cityCode":     "city_code",
		"countyCode":   "county_code",
	}
	UserOrderKeyMap = map[string]string{
		"id":    "id",
		"ctime": "ctime",
	}
)

func regUserValid() {
	valid.RegRule(map[string]validator.Func{
		"validProvinceCode":    validCityCode(LevelProvince),
		"validCityCode":        validCityCode(LevelCity),
		"validCountyCode":      validCityCode(LevelCounty),
		"validUniqueAccount":   validUniqueAccount,
		"validUniquePhone":     validUniquePhone,
		"validUniqueEditPhone": validUniqueEditPhone,
	})
	valid.RegTrans(map[string][]interface{}{
		"validProvinceCode": {
			func(ut ut.Translator) error {
				return ut.Add("validProvinceCode", "省份code{0}不存在",
					true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("validProvinceCode", fe.Value().(string))
				return t
			},
		},
		"validCityCode": {
			func(ut ut.Translator) error {
				return ut.Add("validCityCode", "城市code{0}不存在",
					true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("validCityCode", fe.Value().(string))
				return t
			},
		},
		"validCountyCode": {
			func(ut ut.Translator) error {
				return ut.Add("validCountyCode", "区县code{0}不存在",
					true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("validCountyCode", fe.Value().(string))
				return t
			},
		},
		"validUniqueAccount": {
			func(ut ut.Translator) error {
				return ut.Add("validUniqueAccount", "账号 {0} 已存在",
					true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("validUniqueAccount", fe.Value().(string))
				return t
			},
		},
		"validUniquePhone": {
			func(ut ut.Translator) error {
				return ut.Add("validUniquePhone", "手机号 {0} 已存在",
					true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("validUniquePhone", fe.Value().(string))
				return t
			},
		},
		"validUniqueEditPhone": {
			func(ut ut.Translator) error {
				return ut.Add("validUniqueEditPhone", "手机号 {0} 已存在",
					true)
			}, func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("validUniqueEditPhone", fe.Value().(string))
				return t
			},
		},
	}, global.GetTrans())
}

type AddUserData struct {
	Account  string           `json:"account" binding:"required,validUniqueAccount"`
	Pwd      string           `json:"pwd" binding:"required,min=4,max=32"`
	RePwd    string           `json:"rePwd" binding:"required,eqfield=Pwd"`
	Level    ginApp.UserLevel `json:"level" binding:"required,validUserLevel"`
	Name     string           `json:"name" binding:"required,min=1,max=12"`
	Gender   ginApp.Gender    `json:"gender" binding:"omitempty,validGender"`
	Age      *int             `json:"age" binding:"omitempty,min=0,max=120"`
	Phone    *string          `json:"phone,omitempty" binding:"omitempty,phone,validUniquePhone"`
	AvatarID *int             `json:"avatarID,omitempty"`
	Profile  *string          `json:"profile" binding:"omitempty,max=500"`
}
type EditUserData struct {
	ID       *int             `json:"id" gorm:"primaryKey" binding:"omitempty,min=0"`
	Name     string           `json:"name" binding:"required,min=1,max=12"`
	Level    ginApp.UserLevel `json:"level" binding:"omitempty,validUserLevel"`
	Gender   ginApp.Gender    `json:"gender" binding:"omitempty,validGender"`
	Age      *int             `json:"age" binding:"omitempty,min=0,max=120"`
	Phone    *string          `json:"phone,omitempty" binding:"omitempty,phone,validUniqueEditPhone"`
	AvatarID *int             `json:"avatarID"`
	Profile  *string          `json:"profile" binding:"omitempty,max=500"`
}
type ExportSceneUser struct {
	ID       int              `json:"id"`
	Account  string           `json:"account"`
	Level    ginApp.UserLevel `json:"level" binding:"required"`
	Gender   ginApp.Gender    `json:"gender" binding:"omitempty,validGender"`
	Name     string           `json:"name" binding:"required,min=1,max=12"`
	Age      *int             `json:"age" binding:"omitempty,min=0,max=120"`
	Phone    *string          `json:"phone" binding:"omitempty"`
	AvatarID *int             `json:"avatarID"`
	Avatar   Any              `json:"avatar,omitempty"`
	Profile  *string          `json:"profile" binding:"omitempty,max=500"`
	Ctime    *time.Time       `json:"ctime,omitempty" gorm:"default:datetime"`
	Utime    *time.Time       `json:"utime,omitempty" gorm:"default:datetime"`
}
type DefaultSceneUser struct {
	EditUserData
	ID           int               `json:"id"`
	Account      string            `json:"account"`
	Status       UserStatus        `json:"status"`
	Avatar       *DefaultSceneFile `json:"avatar,omitempty"`
	Ctime        *time.Time        `json:"ctime,omitempty" gorm:"default:datetime"`
	Utime        *time.Time        `json:"utime,omitempty" gorm:"default:datetime"`
	CollectCount *int              `json:"collectCount,omitempty"` // 收藏文章数
	FollowCount  *int              `json:"followCount,omitempty"`  // 关注数
}
type ProfileSceneUser struct {
	ID      int               `json:"id" gorm:"primaryKey" binding:"required,min=0"`
	Name    string            `json:"name" binding:""`
	Profile *string           `json:"profile"`
	Avatar  *DefaultSceneFile `json:"avatar,omitempty"`
}
type AdminSceneUser struct {
	DefaultSceneUser
	Dtime *int64 `json:"dtime,omitempty" gorm:"default:0"`
}
type User struct {
	Base
	Account      string           `json:"account" gorm:"type:varchar(32);not null;index:account,unique"`
	Pwd          string           `json:"pwd" gorm:"type:varchar(128);not null;"`
	Level        ginApp.UserLevel `json:"level" gorm:"type:varchar(16);not null;"`
	Name         string           `json:"name" gorm:"type:varchar(32);not null;default:''"`
	Gender       ginApp.Gender    `json:"gender" gorm:"type:varchar(16);not null;default:'未知'"`
	Status       UserStatus       `json:"status" gorm:"type:varchar(16);not null;default:'正常'"`
	Age          *int             `json:"age,omitempty" gorm:"type:tinyint unsigned;not null;default:0"`
	Phone        *string          `json:"phone,omitempty" gorm:"varchar(16);not null;default:'';index:phone"`
	Email        *string          `json:"email,omitempty" gorm:"varchar(128);not null;default:'';index:email"`
	AvatarID     *int             `json:"avatarID,omitempty" gorm:"type:int unsigned;not null;default:0"`
	Profile      *string          `json:"profile,omitempty" gorm:"type:varchar(512);not null;default:''"`
	ProvinceCode *string          `json:"provinceCode,omitempty" gorm:"type:varchar(16);not null;default:''"`
	CityCode     *string          `json:"cityCode,omitempty"  gorm:"type:varchar(16);not null;default:''"`
	CountyCode   *string          `json:"countyCode,omitempty" gorm:"type:varchar(16);not null;default:''"`
	HasFollow    *bool            `json:"-" gorm:"-"`
}

func NewCtxUser(ctx *gin.Context, userParam ...*User) *User {
	var m *User
	if len(userParam) > 0 {
		m = userParam[0]
	} else {
		m = &User{}
	}
	m.Ctx = ctx
	return m
}
func NewDefaultSceneUser(u *User) *DefaultSceneUser {
	user := &DefaultSceneUser{
		Ctime: u.Ctime,
		Utime: u.Utime,
	}
	user.ID = u.ID
	user.Account = u.Account
	user.Level = u.Level
	user.Gender = u.Gender
	user.Name = u.Name
	user.Age = u.Age
	user.AvatarID = u.AvatarID
	user.Profile = u.Profile
	user.Status = u.Status
	if u.Phone != nil && *u.Phone != "" {
		user.Phone = new(string)
		*user.Phone = bdk.MaskPhone(*u.Phone)
	}
	if avatar := u.GetAvatar(); avatar != nil {
		user.Avatar = avatar.GetFmtDetail().(*DefaultSceneFile)
	}
	return user
}
func NewExportSceneUser(u *User) *ExportSceneUser {
	user := &ExportSceneUser{
		ID:       u.ID,
		Account:  u.Account,
		Level:    u.Level,
		Gender:   u.Gender,
		Name:     u.Name,
		Age:      u.Age,
		Phone:    u.Phone,
		AvatarID: u.AvatarID,
		Avatar:   nil,
		Profile:  u.Profile,
		Ctime:    u.Ctime,
		Utime:    u.Utime,
	}

	if avatar := u.GetAvatar(); avatar != nil {
		user.Avatar = avatar.GetFmtDetail()
	}
	return user
}
func NewProfileSceneUser(u *User) *ProfileSceneUser {
	user := &ProfileSceneUser{}
	user.ID = u.ID
	user.Name = u.Name
	if u.Profile != nil {
		user.Profile = u.Profile
	}
	if avatar := u.GetAvatar(); avatar != nil {
		user.Avatar = avatar.GetFmtDetail().(*DefaultSceneFile)
	}
	return user
}
func NewAdminSceneUser(u *User) *AdminSceneUser {
	defaultUserInfo := NewDefaultSceneUser(u)
	user := &AdminSceneUser{
		DefaultSceneUser: *defaultUserInfo,
		Dtime:            u.Dtime,
	}
	user.Phone = u.Phone
	return user
}
func (m *User) NewModel() BaseModel {
	return &User{}
}
func (m *User) GetFmtDetail(scenes ...string) Any {
	var scene string
	if len(scenes) == 1 {
		scene = scenes[0]
	}
	var user Any
	switch scene {
	case UserSceneProfile:
		user = NewProfileSceneUser(m)
	case SceneAdmin:
		user = NewAdminSceneUser(m)
	case UserSceneExport:
		user = NewExportSceneUser(m)
	default:
		user = NewDefaultSceneUser(m)
	}
	return user
}
func (m *User) GetFilterMap() map[string]string {
	return UserFilterNameMapDBField
}
func (m *User) GetOrderMap() map[string]string {
	return UserOrderKeyMap
}
func (m *User) GetDetail(id int) error {
	return Detail(m, id)
}
func (m *User) GetGormQuery() *gorm.DB {
	db := global.DB
	if m.Ctx != nil {
		db = db.WithContext(m.Ctx)
	}
	return db.Model(&User{}).Scopes(NotDelScope)
}
func (m *User) GormFindList(q *gorm.DB) Any {
	actModels := make([]*User, 0, 1)
	// nolint
	q = q.Find(&actModels)
	return &actModels
}
func (m *User) GetFmtList(list Any, scene string) Any {
	fmtList := make([]Any, 0, 1)
	actList := list.(*[]*User)
	for _, item := range *actList {
		fmtList = append(fmtList, item.GetFmtDetail(scene))
	}
	return fmtList
}

func (m *User) GetAvatar() *File {
	relation := &File{}
	if m.AvatarID == nil || *m.AvatarID == 0 {
		return nil
	}
	global.DB.Model(relation).Scopes(NotDelScope).
		Where("id = ?", m.AvatarID).First(relation)
	if relation.ID == 0 {
		return nil
	}
	return relation
}
func (m *User) IsAccountExist(account string) bool {
	count := int64(0)
	m.GetGormQuery().Where("account = ?", account).Count(&count)
	return count > 0
}
func (m *User) IsPhoneExist(phone string) bool {
	count := int64(0)
	m.GetGormQuery().Where("phone = ?", phone).Count(&count)
	return count > 0
}
func isCityExist(code, level string) bool {
	m := &City{}
	count := int64(0)
	m.GetGormQuery().Where("level = ? and code = ?", level, code).Count(&count)
	return count > 0
}

// valid
func validCityCode(level string) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		code := fl.Field().String()
		return code == "" || isCityExist(code, level)
	}
}
func validUniqueAccount(fl validator.FieldLevel) bool {
	account := fl.Field().String()
	m := &User{}
	return !m.IsAccountExist(account)
}
func validUniquePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	m := &User{}
	return !m.IsPhoneExist(phone)
}

// 判断编辑的用户手机号是否已存在
func validUniqueEditPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	iRecord := fl.Top().Interface()
	if record, ok := iRecord.(*EditUserData); ok {
		userModel := &User{}
		if err := userModel.GetDetail(*record.ID); err != nil {
			return false
		}
		if userModel.Phone != nil && phone == *userModel.Phone {
			return true
		}
	}
	m := &User{}
	return !m.IsPhoneExist(phone)
}
