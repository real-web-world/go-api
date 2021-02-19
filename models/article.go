package models

import (
	"gorm.io/gorm"

	"github.com/real-web-world/go-web-api/global"
	"github.com/real-web-world/go-web-api/pkg/fastcurd"
	"github.com/real-web-world/go-web-api/pkg/gin"
)

type Article struct {
	Base
	Title      string `json:"title" gorm:"type:varchar(128);not null"`
	CategoryID int    `json:"categoryID" gorm:"type:int unsigned;not null;default:0"`
	ViewCount  int    `json:"viewCount" gorm:"type:int unsigned;not null;default:0"`
	StarCount  int    `json:"starCount" gorm:"type:int unsigned;not null;default:0"`
	Excerpt    string `json:"excerpt" gorm:"type:varchar(200) unsigned;not null;default:''"`
	Content    string `json:"content" gorm:"type:mediumtext;not null"`
	AuthorUid  int    `json:"author_uid" gorm:"type:int unsigned;not null;default:0"`
}

var (
	ArticleFilterNameMapDBField = map[string]string{
		"search":     "id",
		"id":         "id",
		"ctime":      "ctime",
		"categoryID": "category_id",
	}

	ArticleOrderKeyMap = map[string]string{
		"id":    "id",
		"ctime": "ctime",
	}
)

type DefaultSceneArticle struct {
	*Article
}
type AddArticleData struct {
}
type EditArticleData struct {
	ID int `json:"id" binding:"required"`
}

func NewDefaultSceneArticle(m *Article) *DefaultSceneArticle {
	model := &DefaultSceneArticle{
		m,
	}
	return model
}
func (m *Article) NewModel() BaseModel {
	return &Article{
		Base: Base{Ctx: m.Ctx},
	}
}
func (m *Article) GetFilterMap() map[string]string {
	return ArticleFilterNameMapDBField
}
func (m *Article) GetOrderMap() map[string]string {
	return ArticleOrderKeyMap
}

func (m *Article) NewAddData() Any {
	return &AddArticleData{}
}
func (m *Article) NewEditData() Any {
	return &EditArticleData{}
}
func (m *Article) AddWithPostData(data Any, app *ginApp.App) error {
	d := data.(*AddArticleData)
	_ = d
	// m.xx= d.xx
	return Add(m)
}
func (m *Article) EditWithPostData(data Any, app *ginApp.App) (int, error) {
	d := data.(*EditArticleData)
	m.ID = d.ID
	return Update(m)
}
func (m *Article) DelWithPostData(data Any, app *ginApp.App) (int, error) {
	d := data.(*fastcurd.DelData)
	return Delete(m, d.IDs)
}

func (m *Article) GetDetail(id int) error {
	return Detail(m, id)
}
func (m *Article) GetGormQuery() *gorm.DB {
	db := global.DB
	if m.Ctx != nil {
		db = db.WithContext(m.Ctx)
	}
	return db.Model(&Article{}).Scopes(NotDelScope)
}
func (m *Article) GormFindList(q *gorm.DB) Any {
	actModels := make([]*Article, 0, 1)
	// nolint
	q = q.Find(&actModels)
	return &actModels
}
func (m *Article) GetFmtList(list Any, scene string) Any {
	fmtList := make([]Any, 0, 1)
	actList := list.(*[]*Article)
	for _, item := range *actList {
		fmtList = append(fmtList, item.GetFmtDetail(scene))
	}
	return fmtList
}
func (m *Article) GetFmtDetail(scenes ...string) Any {
	var scene string
	if len(scenes) == 1 {
		scene = scenes[0]
	}
	var model Any
	// nolint
	switch scene {
	default:
		model = NewDefaultSceneArticle(m)
	}
	return model
}
