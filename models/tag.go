package models

import (
	"gorm.io/gorm"

	"github.com/real-web-world/go-web-api/global"
	"github.com/real-web-world/go-web-api/pkg/fastcurd"
	"github.com/real-web-world/go-web-api/pkg/gin"
)

type Tag struct {
	Base
	Name string `json:"name" gorm:"type:varchar(32);not null"`
}

var (
	TagFilterNameMapDBField = map[string]string{
		"search": "id",
		"id":     "id",
		"ctime":  "ctime",
	}
	TagOrderKeyMap = map[string]string{
		"id":    "id",
		"ctime": "ctime",
	}
)

type DefaultSceneTag struct {
	*Tag
}
type AddTagData struct {
	Name string `json:"name" binding:"required,min=1,max=32"`
}
type EditTagData struct {
	ID   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required,min=1,max=32"`
}

func NewDefaultSceneTag(m *Tag) *DefaultSceneTag {
	model := &DefaultSceneTag{
		m,
	}
	return model
}
func (m *Tag) NewModel() BaseModel {
	return &Tag{
		Base: Base{Ctx: m.Ctx},
	}
}
func (m *Tag) GetFilterMap() map[string]string {
	return TagFilterNameMapDBField
}
func (m *Tag) GetOrderMap() map[string]string {
	return TagOrderKeyMap
}

func (m *Tag) NewAddData() Any {
	return &AddTagData{}
}
func (m *Tag) NewEditData() Any {
	return &EditTagData{}
}
func (m *Tag) AddWithPostData(data Any, app *ginApp.App) error {
	d := data.(*AddTagData)
	_ = d
	m.Name = d.Name
	return Add(m)
}
func (m *Tag) EditWithPostData(data Any, app *ginApp.App) (int, error) {
	d := data.(*EditTagData)
	m.ID = d.ID
	m.Name = d.Name
	return Update(m)
}
func (m *Tag) DelWithPostData(data Any, app *ginApp.App) (int, error) {
	d := data.(*fastcurd.DelData)
	return Delete(m, d.IDs)
}

func (m *Tag) GetDetail(id int) error {
	return Detail(m, id)
}
func (m *Tag) GetGormQuery() *gorm.DB {
	db := global.DB
	if m.Ctx != nil {
		db = db.WithContext(m.Ctx)
	}
	return db.Model(&Tag{}).Scopes(NotDelScope)
}
func (m *Tag) GormFindList(q *gorm.DB) Any {
	actModels := make([]*Tag, 0, 1)
	// nolint
	q = q.Find(&actModels)
	return &actModels
}
func (m *Tag) GetFmtList(list Any, scene string) Any {
	fmtList := make([]Any, 0, 1)
	actList := list.(*[]*Tag)
	for _, item := range *actList {
		fmtList = append(fmtList, item.GetFmtDetail(scene))
	}
	return fmtList
}
func (m *Tag) GetFmtDetail(scenes ...string) Any {
	var scene string
	if len(scenes) == 1 {
		scene = scenes[0]
	}
	var model Any
	// nolint
	switch scene {
	default:
		model = NewDefaultSceneTag(m)
	}
	return model
}
