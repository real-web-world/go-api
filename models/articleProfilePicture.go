package models

import (
	"gorm.io/gorm"

	"github.com/real-web-world/go-api/global"
	"github.com/real-web-world/go-api/pkg/fastcurd"
	"github.com/real-web-world/go-api/pkg/gin"
)

type ArticleProfilePicture struct {
	Base
	ArticleID int `json:"title" gorm:"type:varchar(128);not null;index:articleID"`
	PicID     int `json:"picID" gorm:"type:int unsigned;not null;default:0;index:picID"`
	Sort      int `json:"sort" gorm:"type:int unsigned;not null;default:0"`
}

var (
	ArticleProfilePictureFilterNameMapDBField = map[string]string{
		"search":    "id",
		"id":        "id",
		"articleID": "article_id",
		"picID":     "pic_id",
		"sort":      "sort",
		"ctime":     "ctime",
	}
	ArticleProfilePictureOrderKeyMap = map[string]string{
		"id":    "id",
		"ctime": "ctime",
	}
)

type DefaultSceneArticleProfilePicture struct {
	*ArticleProfilePicture
}
type AddArticleProfilePictureData struct {
}
type EditArticleProfilePictureData struct {
	ID int `json:"id" binding:"required"`
}

func NewDefaultSceneArticleProfilePicture(m *ArticleProfilePicture) *DefaultSceneArticleProfilePicture {
	model := &DefaultSceneArticleProfilePicture{
		m,
	}
	return model
}
func (m *ArticleProfilePicture) NewModel() BaseModel {
	return &ArticleProfilePicture{
		Base: Base{Ctx: m.Ctx},
	}
}
func (m *ArticleProfilePicture) GetFilterMap() map[string]string {
	return ArticleProfilePictureFilterNameMapDBField
}
func (m *ArticleProfilePicture) GetOrderMap() map[string]string {
	return ArticleProfilePictureOrderKeyMap
}

func (m *ArticleProfilePicture) NewAddData() Any {
	return &AddArticleProfilePictureData{}
}
func (m *ArticleProfilePicture) NewEditData() Any {
	return &EditArticleProfilePictureData{}
}
func (m *ArticleProfilePicture) AddWithPostData(data Any, app *ginApp.App) error {
	d := data.(*AddArticleProfilePictureData)
	_ = d
	// m.xx= d.xx
	return Add(m)
}
func (m *ArticleProfilePicture) EditWithPostData(data Any, app *ginApp.App) (int, error) {
	d := data.(*EditArticleProfilePictureData)
	m.ID = d.ID
	return Update(m)
}
func (m *ArticleProfilePicture) DelWithPostData(data Any, app *ginApp.App) (int, error) {
	d := data.(*fastcurd.DelData)
	return Delete(m, d.IDs)
}

func (m *ArticleProfilePicture) GetDetail(id int) error {
	return Detail(m, id)
}
func (m *ArticleProfilePicture) GetGormQuery() *gorm.DB {
	db := global.DB
	if m.Ctx != nil {
		db = db.WithContext(m.Ctx)
	}
	return db.Model(&ArticleProfilePicture{}).Scopes(NotDelScope)
}
func (m *ArticleProfilePicture) GormFindList(q *gorm.DB) Any {
	actModels := make([]*ArticleProfilePicture, 0, 1)
	// nolint
	q = q.Find(&actModels)
	return &actModels
}
func (m *ArticleProfilePicture) GetFmtList(list Any, scene string) Any {
	fmtList := make([]Any, 0, 1)
	actList := list.(*[]*ArticleProfilePicture)
	for _, item := range *actList {
		fmtList = append(fmtList, item.GetFmtDetail(scene))
	}
	return fmtList
}
func (m *ArticleProfilePicture) GetFmtDetail(scenes ...string) Any {
	var scene string
	if len(scenes) == 1 {
		scene = scenes[0]
	}
	var model Any
	// nolint
	switch scene {
	default:
		model = NewDefaultSceneArticleProfilePicture(m)
	}
	return model
}
func (m *ArticleProfilePicture) GetFile() *File {
	relation := &File{Base: Base{
		Ctx: m.Ctx,
	}}
	if m.PicID == 0 {
		return nil
	}
	relation.GetGormQuery().Where("id = ?", m.PicID).First(relation)
	if relation.ID == 0 {
		return nil
	}
	return relation
}
