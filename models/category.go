package models

import (
	"gorm.io/gorm"

	"github.com/real-web-world/go-api/global"
	"github.com/real-web-world/go-api/pkg/fastcurd"
	"github.com/real-web-world/go-api/pkg/gin"
)

type Category struct {
	Base
	Name string `json:"name" gorm:"type:varchar(32);not null"`
	Pid  int    `json:"pid"  gorm:"type:int unsigned;not null;default:0"`
}

const (
	sceneWithArticleCount = "withArticleCount"
)

var (
	CategoryFilterNameMapDBField = map[string]string{
		"search": "id",
		"id":     "id",
		"name":   "name",
		"pid":    "pid",
		"ctime":  "ctime",
	}
	CategoryOrderKeyMap = map[string]string{
		"id":    "id",
		"ctime": "ctime",
	}
)

type DefaultSceneCategory struct {
	*Category
}
type AddCategoryData struct {
	Name string `json:"name" binding:"required,min=1,max=32"`
	Pid  int    `json:"pid" binding:"required,min=0"`
}
type EditCategoryData struct {
	AddCategoryData
	ID int `json:"id" binding:"required"`
}
type WithArticleCountSceneCategory struct {
	*Category
	ArticleCount int `json:"articleCount"`
}

func NewDefaultSceneCategory(m *Category) *DefaultSceneCategory {
	model := &DefaultSceneCategory{
		m,
	}
	return model
}
func NewWithArticleCountSceneCategory(m *Category) *WithArticleCountSceneCategory {
	model := &WithArticleCountSceneCategory{
		Category: m,
	}
	model.ArticleCount = m.GetArticleCount()
	return model
}
func (m *Category) NewModel() BaseModel {
	return &Category{
		Base: Base{Ctx: m.Ctx},
	}
}
func (m *Category) GetFilterMap() map[string]string {
	return CategoryFilterNameMapDBField
}
func (m *Category) GetOrderMap() map[string]string {
	return CategoryOrderKeyMap
}

func (m *Category) NewAddData() Any {
	return &AddCategoryData{}
}
func (m *Category) NewEditData() Any {
	return &EditCategoryData{}
}
func (m *Category) AddWithPostData(data Any, app *ginApp.App) error {
	d := data.(*AddCategoryData)
	_ = d
	m.Name = d.Name
	m.Pid = d.Pid
	return Add(m)
}
func (m *Category) EditWithPostData(data Any, app *ginApp.App) (int, error) {
	d := data.(*EditCategoryData)
	m.ID = d.ID
	m.Name = d.Name
	return Update(m)
}
func (m *Category) DelWithPostData(data Any, app *ginApp.App) (int, error) {
	d := data.(*fastcurd.DelData)
	return Delete(m, d.IDs)
}

func (m *Category) GetDetail(id int) error {
	return Detail(m, id)
}
func (m *Category) GetGormQuery() *gorm.DB {
	db := global.DB
	if m.Ctx != nil {
		db = db.WithContext(m.Ctx)
	}
	return db.Model(&Category{}).Scopes(NotDelScope)
}
func (m *Category) GormFindList(q *gorm.DB) Any {
	actModels := make([]*Category, 0, 1)
	// nolint
	q = q.Find(&actModels)
	return &actModels
}
func (m *Category) GetFmtList(list Any, scene string) Any {
	fmtList := make([]Any, 0, 1)
	actList := list.(*[]*Category)
	for _, item := range *actList {
		fmtList = append(fmtList, item.GetFmtDetail(scene))
	}
	return fmtList
}
func (m *Category) GetFmtDetail(scenes ...string) Any {
	var scene string
	if len(scenes) == 1 {
		scene = scenes[0]
	}
	var model Any
	// nolint
	switch scene {
	case sceneWithArticleCount:
		model = NewWithArticleCountSceneCategory(m)
	default:
		model = NewDefaultSceneCategory(m)
	}
	return model
}

func (m *Category) GetArticleCount() int {
	article := &Article{Base: Base{Ctx: m.Ctx}}
	var count int64 = 0
	article.GetGormQuery().Where("category_id", m.ID).Count(&count)
	return int(count)
}
