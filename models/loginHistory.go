package models

import (
	"gorm.io/gorm"

	"github.com/real-web-world/go-web-api/global"
	"github.com/real-web-world/go-web-api/pkg/fastcurd"
	"github.com/real-web-world/go-web-api/pkg/gin"
)

type Source string

const (
	SourceAndroid Source = "android"
	SourceIos     Source = "ios"
	SourceWeb     Source = "web"
)

type LoginHistory struct {
	Base
	Account       string `json:"account" gorm:"type:varchar(128);not null;default:''"`
	Phone         string `json:"phone" gorm:"type:varchar(16);not null;default:''"`
	OK            string `json:"ok" gorm:"type:varchar(8);not null;"`
	Source        Source `json:"source" gorm:"type:varchar(128);not null;default:''"`
	ClientIP      string `json:"clientIP" gorm:"type:varchar(32);not null;"`
	AppUniqueID   string `json:"appUniqueID" gorm:"type:varchar(128);not null;default:''"`
	WebUserAgent  string `json:"webUserAgent" gorm:"type:varchar(256);not null;default:''"`
	AppDeviceInfo string `json:"appDeviceInfo" gorm:"type:varchar(1024);not null;default:''"`
}
type AddLoginHistoryData struct {
	OK            string  `json:"ok"`
	Source        Source  `json:"source"`
	ClientIP      string  `json:"client_ip"`
	Account       *string `json:"account"`
	Phone         *string `json:"phone"`
	AppUniqueID   *string `json:"appUniqueID"`
	WebUserAgent  *string `json:"webUserAgent"`
	AppDeviceInfo *string `json:"appDeviceInfo"`
}

var (
	LoginHistoryFilterNameMapDBField = map[string]string{
		"search":  "id|account|phone",
		"id":      "id",
		"account": "account",
		"phone":   "phone",
		"ok":      "ok",
		"source":  "source",
		"ctime":   "ctime",
	}
	LoginHistoryOrderKeyMap = map[string]string{
		"id":    "id",
		"ctime": "ctime",
	}
)

type DefaultSceneLoginHistory struct {
	*LoginHistory
}

func NewDefaultSceneLoginHistory(m *LoginHistory) *DefaultSceneLoginHistory {
	model := &DefaultSceneLoginHistory{
		m,
	}
	return model
}
func (m *LoginHistory) NewModel() BaseModel {
	return &LoginHistory{
		Base: Base{Ctx: m.Ctx},
	}
}
func (m *LoginHistory) GetFilterMap() map[string]string {
	return LoginHistoryFilterNameMapDBField
}
func (m *LoginHistory) GetOrderMap() map[string]string {
	return LoginHistoryOrderKeyMap
}
func (m *LoginHistory) DelWithPostData(data Any, _ *ginApp.App) (int, error) {
	d := data.(*fastcurd.DelData)
	return Delete(m, d.IDs)
}

func (m *LoginHistory) GetDetail(id int) error {
	return Detail(m, id)
}
func (m *LoginHistory) GetGormQuery() *gorm.DB {
	db := global.DB
	if m.Ctx != nil {
		db = db.WithContext(m.Ctx)
	}
	return db.Model(&LoginHistory{})
}
func (m *LoginHistory) GormFindList(q *gorm.DB) Any {
	actModels := make([]*LoginHistory, 0, 1)
	// nolint
	q = q.Find(&actModels)
	return &actModels
}

func (m *LoginHistory) GetFmtList(list Any, scene string) Any {
	fmtList := make([]Any, 0, 1)
	actList := list.(*[]*LoginHistory)
	for _, item := range *actList {
		fmtList = append(fmtList, item.GetFmtDetail(scene))
	}
	return fmtList
}
func (m *LoginHistory) GetFmtDetail(scenes ...string) Any {
	var scene string
	if len(scenes) == 1 {
		scene = scenes[0]
	}
	var model Any
	// nolint
	switch scene {
	default:
		model = NewDefaultSceneLoginHistory(m)
	}
	return model
}
func AddLoginHistory(d *AddLoginHistoryData) error {
	m := &LoginHistory{
		OK:       d.OK,
		Source:   d.Source,
		ClientIP: d.ClientIP,
	}
	if d.Phone != nil {
		m.Phone = *d.Phone
	}
	if d.Account != nil {
		m.Account = *d.Account
	}
	if d.AppDeviceInfo != nil {
		m.AppDeviceInfo = *d.AppDeviceInfo
	}
	if d.AppUniqueID != nil {
		m.AppUniqueID = *d.AppUniqueID
	}
	if d.WebUserAgent != nil {
		m.WebUserAgent = *d.WebUserAgent
	}
	return Add(m)
}
