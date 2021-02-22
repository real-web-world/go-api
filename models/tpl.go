// +build ignore

package models

import (
	"gorm.io/gorm"

	"github.com/real-web-world/go-api/global"
	"github.com/real-web-world/go-api/pkg/fastcurd"
	"github.com/real-web-world/go-api/pkg/gin"
)

type Tpl struct {
	Base
}

var (
	TplFilterNameMapDBField = map[string]string{
		"search": "id",
		"id":     "id",
		"ctime":  "ctime",
	}
	TplOrderKeyMap = map[string]string{
		"id":    "id",
		"ctime": "ctime",
	}
)

type DefaultSceneTpl struct {
	*Tpl
}
type AddTplData struct {
}
type EditTplData struct {
	ID int `json:"id" binding:"required"`
}

func NewDefaultSceneTpl(m *Tpl) *DefaultSceneTpl {
	model := &DefaultSceneTpl{
		m,
	}
	return model
}
func (m *Tpl) NewModel() BaseModel {
	return &Tpl{
		Base: Base{Ctx: m.Ctx},
	}
}
func (m *Tpl) GetFilterMap() map[string]string {
	return TplFilterNameMapDBField
}
func (m *Tpl) GetOrderMap() map[string]string {
	return TplOrderKeyMap
}

func (m *Tpl) NewAddData() Any {
	return &AddTplData{}
}
func (m *Tpl) NewEditData() Any {
	return &EditTplData{}
}
func (m *Tpl) AddWithPostData(data Any, _ *ginApp.App) error {
	d := data.(*AddTplData)
	_ = d
	// m.xx= d.xx
	return Add(m)
}
func (m *Tpl) EditWithPostData(data Any, _ *ginApp.App) (int, error) {
	d := data.(*EditTplData)
	m.ID = d.ID
	return Update(m)
}
func (m *Tpl) DelWithPostData(data Any, _ *ginApp.App) (int, error) {
	d := data.(*fastcurd.DelData)
	return Delete(m, d.IDs)
}

func (m *Tpl) GetDetail(id int) error {
	return Detail(m, id)
}
func (m *Tpl) GetGormQuery() *gorm.DB {
	db := global.DB
	if m.Ctx != nil {
		db = db.WithContext(m.Ctx)
	}
	return db.Model(&Tpl{}).Scopes(NotDelScope)
}
func (m *Tpl) GormFindList(q *gorm.DB) Any {
	actModels := make([]*Tpl, 0, 1)
	// nolint
	q = q.Find(&actModels)
	return &actModels
}
func (m *Tpl) GetFmtList(list Any, scene string) Any {
	fmtList := make([]Any, 0, 1)
	actList := list.(*[]*Tpl)
	for _, item := range *actList {
		item.SetCtx(m.Ctx)
		fmtList = append(fmtList, item.GetFmtDetail(scene))
	}
	return fmtList
}
func (m *Tpl) GetFmtDetail(scenes ...string) Any {
	var scene string
	if len(scenes) == 1 {
		scene = scenes[0]
	}
	var model Any
	// nolint
	switch scene {
	default:
		model = NewDefaultSceneTpl(m)
	}
	return model
}
