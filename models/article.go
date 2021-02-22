package models

import (
	"github.com/real-web-world/go-api/global"
	"github.com/real-web-world/go-api/pkg/fastcurd"
	ginApp "github.com/real-web-world/go-api/pkg/gin"
	"gorm.io/gorm"
)

type Article struct {
	Base
	Title      *string `json:"title" gorm:"type:varchar(128);not null"`
	CategoryID *int    `json:"categoryID" gorm:"type:int unsigned;not null;default:0"`
	ViewCount  *int    `json:"viewCount" gorm:"type:int unsigned;not null;default:0"`
	StarCount  *int    `json:"starCount" gorm:"type:int unsigned;not null;default:0"`
	Excerpt    *string `json:"excerpt" gorm:"type:varchar(200);not null;default:''"`
	Content    *string `json:"content" gorm:"type:mediumtext;not null"`
	AuthorUid  *int    `json:"author_uid" gorm:"type:int unsigned;not null;default:0"`
}

const (
	ArticleSceneView = "view"
)

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
	// lazy skill
	CategoryID int                   `json:"-"`
	AuthorUid  int                   `json:"-"`
	Category   *DefaultSceneCategory `json:"category,omitempty"`
	Author     *ProfileSceneUser     `json:"author,omitempty"`
	Tags       []*DefaultSceneTag    `json:"tags"`
	Pictures   []*DefaultSceneFile   `json:"pictures"`
}
type AddArticleData struct {
	Title      *string `json:"title" binding:"required,min=1,max=100"`
	CategoryID *int    `json:"categoryID" binding:"required"`
	Excerpt    *string `json:"excerpt" binding:"required,max=300"`
	Content    *string `json:"content" binding:"required,min=10"`
	TagIDList  []int   `json:"tagIDList" binding:""`
	PicIDList  []int   `json:"PicIDList" binding:"max=3"`
}
type EditArticleData struct {
	ID         int     `json:"id" binding:"required"`
	Title      *string `json:"title" binding:"omitempty,min=1,max=100"`
	CategoryID *int    `json:"categoryID" binding:""`
	Excerpt    *string `json:"excerpt" binding:"omitempty,max=300"`
	Content    *string `json:"content" binding:"omitempty,min=10"`
	TagIDList  []int   `json:"tagIDList" binding:""`
	PicIDList  []int   `json:"picIDList" binding:"omitempty,max=3"`
}

func NewDefaultSceneArticle(m *Article) *DefaultSceneArticle {
	model := &DefaultSceneArticle{
		Article: m,
	}
	author := m.GetAuthor()
	if author != nil {
		model.Author = NewProfileSceneUser(author)
	}
	category := m.GetCategory()
	if category != nil {
		model.Category = NewDefaultSceneCategory(category)
	}
	model.Tags = m.GetDefaultSceneTagList()
	model.Pictures = m.GetDefaultScenePicList()
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
	m.CategoryID = d.CategoryID
	m.Content = d.Content
	m.Excerpt = d.Excerpt
	m.Title = d.Title
	m.AuthorUid = new(int)
	*m.AuthorUid = app.AuthUser.UID
	err := Add(m)
	if err == nil {
		go func() {
			m.UpdateCoverPicIDList(d.PicIDList)
			m.UpdateTagIDList(d.TagIDList)
		}()
	}
	return err
}
func (m *Article) EditWithPostData(data Any, _ *ginApp.App) (int, error) {
	d := data.(*EditArticleData)
	m.ID = d.ID
	m.CategoryID = d.CategoryID
	m.Content = d.Content
	m.Excerpt = d.Excerpt
	m.Title = d.Title
	queryChan := make(chan struct{}, 1)
	if d.PicIDList != nil {
		m.RelationAffectRows += m.UpdateCoverPicIDList(d.PicIDList)
	}
	func() {
		if d.TagIDList != nil {
			m.RelationAffectRows += m.UpdateTagIDList(d.TagIDList)
		}
		queryChan <- struct{}{}
	}()
	<-queryChan
	return Update(m)
}
func (m *Article) DelWithPostData(data Any, _ *ginApp.App) (int, error) {
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
	case ArticleSceneView:
		go func() {
			*m.ViewCount = *m.ViewCount + 1
			_, _ = Update(&Article{Base: Base{Ctx: m.Ctx, ID: m.ID}, ViewCount: m.ViewCount})
		}()
		model = NewDefaultSceneArticle(m)
	default:
		model = NewDefaultSceneArticle(m)
	}
	return model
}
func (m *Article) GetAuthor() *User {
	relation := &User{Base: Base{
		Ctx: m.Ctx,
	}}
	if m.AuthorUid == nil || *m.AuthorUid == 0 {
		return nil
	}
	relation.GetGormQuery().Where("id = ?", m.AuthorUid).First(relation)
	if relation.ID == 0 {
		return nil
	}
	return relation
}
func (m *Article) GetCategory() *Category {
	relation := &Category{Base: Base{
		Ctx: m.Ctx,
	}}
	if m.CategoryID == nil || *m.CategoryID == 0 {
		return nil
	}
	relation.GetGormQuery().Where("id = ?", m.CategoryID).First(relation)
	if relation.ID == 0 {
		return nil
	}
	return relation
}
func (m *Article) GetTagList() []*Tag {
	relationList := make([]*ArticleTag, 0, 4)
	tagList := make([]*Tag, 0, 4)
	model := ArticleTag{Base: Base{
		Ctx: m.Ctx,
	}}
	model.GetGormQuery().Where("article_id = ?",
		m.ID).Order("sort desc").Find(&relationList)
	if len(relationList) == 0 {
		return nil
	}
	for _, relation := range relationList {
		if tag := relation.GetTag(); tag != nil {
			tagList = append(tagList, tag)
		}
	}
	return tagList
}
func (m *Article) GetDefaultSceneTagList() []*DefaultSceneTag {
	if tags := m.GetTagList(); tags != nil {
		fmtTags := make([]*DefaultSceneTag, 0, 1)
		for _, tag := range tags {
			fmtTags = append(fmtTags, NewDefaultSceneTag(tag))
		}
		return fmtTags
	}
	return nil
}
func (m *Article) GetPicList() *[]*File {
	relationList := make([]*ArticleProfilePicture, 0, 1)
	picList := make([]*File, 0, 1)
	model := ArticleProfilePicture{Base: Base{
		Ctx: m.Ctx,
	}}
	model.GetGormQuery().Where("article_id = ?",
		m.ID).Order("sort desc").Find(&relationList)
	if len(relationList) == 0 {
		return nil
	}
	for _, relation := range relationList {
		if pic := relation.GetFile(); pic != nil {
			picList = append(picList, pic)
		}
	}
	return &picList
}
func (m *Article) GetDefaultScenePicList() []*DefaultSceneFile {
	if picList := m.GetPicList(); picList != nil {
		fmtPicList := make([]*DefaultSceneFile, 0, 1)
		for _, pic := range *picList {
			fmtPicList = append(fmtPicList, NewDefaultSceneFile(pic))
		}
		return fmtPicList
	}
	return nil
}
func (m *Article) UpdateCoverPicIDList(picIDList []int) int {
	if picIDList == nil {
		return 0
	}
	affectRows := 0
	for _, id := range picIDList {
		if !IsExist(&ArticleProfilePicture{Base: Base{Ctx: m.Ctx}}, map[string]interface{}{
			"pic_id":     id,
			"article_id": m.ID,
		}) {
			if Add(&ArticleProfilePicture{
				Base:      Base{Ctx: m.Ctx},
				PicID:     id,
				ArticleID: m.ID,
			}) == nil {
				affectRows++
			}
		}
	}
	if len(picIDList) == 0 {
		picIDList = append(picIDList, 0)
	}
	filter := fastcurd.Filter{
		"articleID": {
			Condition: fastcurd.CondEq,
			Val:       m.ID,
		},
		"picID": {
			Condition: fastcurd.CondNotIn,
			Val:       picIDList,
		},
	}
	delCount, _ := Delete(&ArticleProfilePicture{Base: Base{Ctx: m.Ctx}}, filter)
	return affectRows + delCount
}
func (m *Article) UpdateTagIDList(tagIDList []int) int {
	if tagIDList == nil {
		return 0
	}
	affectRows := 0
	for _, id := range tagIDList {
		if !IsExist(&ArticleTag{Base: Base{Ctx: m.Ctx}}, map[string]interface{}{
			"tag_id":     id,
			"article_id": m.ID,
		}) {
			if Add(&ArticleTag{
				Base:      Base{Ctx: m.Ctx},
				TagID:     id,
				ArticleID: m.ID,
			}) == nil {
				affectRows++
			}
		}
	}
	if len(tagIDList) == 0 {
		tagIDList = append(tagIDList, 0)
	}
	filter := fastcurd.Filter{
		"articleID": {
			Condition: fastcurd.CondEq,
			Val:       m.ID,
		},
		"tagID": {
			Condition: fastcurd.CondNotIn,
			Val:       tagIDList,
		},
	}
	delCount, _ := Delete(&ArticleTag{Base: Base{Ctx: m.Ctx}}, filter)
	return affectRows + delCount
}
