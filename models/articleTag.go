package models

import (
	"gorm.io/gorm"

	"github.com/real-web-world/go-api/global"
	"github.com/real-web-world/go-api/pkg/fastcurd"
	"github.com/real-web-world/go-api/pkg/gin"
)

type ArticleTag struct {
	Base
	ArticleID int `json:"articleID" gorm:"type:int unsigned;not null;default:0;index:articleID"`
	TagID     int `json:"tagID" gorm:"type:int unsigned;not null;default:0;index:tagID"`
	Sort      int `json:"sort" gorm:"type:int unsigned;not null;default:0"`
}

var (
	ArticleTagFilterNameMapDBField = map[string]string{
		"search":    "id",
		"id":        "id",
		"articleID": "article_id",
		"tagID":     "tag_id",
		"sort":      "sort",
		"ctime":     "ctime",
	}
	ArticleTagOrderKeyMap = map[string]string{
		"id":    "id",
		"ctime": "ctime",
	}
)

type DefaultSceneArticleTag struct {
	*ArticleTag
}
type AddArticleTagData struct {
}
type EditArticleTagData struct {
	ID int `json:"id" binding:"required"`
}

func NewDefaultSceneArticleTag(m *ArticleTag) *DefaultSceneArticleTag {
	model := &DefaultSceneArticleTag{
		m,
	}
	return model
}
func (m *ArticleTag) NewModel() BaseModel {
	return &ArticleTag{
		Base: Base{Ctx: m.Ctx},
	}
}
func (m *ArticleTag) GetFilterMap() map[string]string {
	return ArticleTagFilterNameMapDBField
}
func (m *ArticleTag) GetOrderMap() map[string]string {
	return ArticleTagOrderKeyMap
}

func (m *ArticleTag) NewAddData() Any {
	return &AddArticleTagData{}
}
func (m *ArticleTag) NewEditData() Any {
	return &EditArticleTagData{}
}
func (m *ArticleTag) AddWithPostData(data Any, app *ginApp.App) error {
	d := data.(*AddArticleTagData)
	_ = d
	// m.xx= d.xx
	return Add(m)
}
func (m *ArticleTag) EditWithPostData(data Any, app *ginApp.App) (int, error) {
	d := data.(*EditArticleTagData)
	m.ID = d.ID
	return Update(m)
}
func (m *ArticleTag) DelWithPostData(data Any, app *ginApp.App) (int, error) {
	d := data.(*fastcurd.DelData)
	return Delete(m, d.IDs)
}

func (m *ArticleTag) GetDetail(id int) error {
	return Detail(m, id)
}
func (m *ArticleTag) GetGormQuery() *gorm.DB {
	db := global.DB
	if m.Ctx != nil {
		db = db.WithContext(m.Ctx)
	}
	return db.Model(&ArticleTag{}).Scopes(NotDelScope)
}
func (m *ArticleTag) GormFindList(q *gorm.DB) Any {
	actModels := make([]*ArticleTag, 0, 1)
	// nolint
	q = q.Find(&actModels)
	return &actModels
}
func (m *ArticleTag) GetFmtList(list Any, scene string) Any {
	fmtList := make([]Any, 0, 1)
	actList := list.(*[]*ArticleTag)
	for _, item := range *actList {
		fmtList = append(fmtList, item.GetFmtDetail(scene))
	}
	return fmtList
}
func (m *ArticleTag) GetFmtDetail(scenes ...string) Any {
	var scene string
	if len(scenes) == 1 {
		scene = scenes[0]
	}
	var model Any
	// nolint
	switch scene {
	default:
		model = NewDefaultSceneArticleTag(m)
	}
	return model
}
func (m *ArticleTag) GetTag() *Tag {
	relation := &Tag{Base: Base{
		Ctx: m.Ctx,
	}}
	if m.TagID == 0 {
		return nil
	}
	relation.GetGormQuery().Where("id = ?", m.TagID).First(relation)
	if relation.ID == 0 {
		return nil
	}
	return relation
}
