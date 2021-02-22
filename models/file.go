package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/real-web-world/go-api/global"
	"github.com/real-web-world/go-api/pkg/fastcurd"
	"github.com/real-web-world/go-api/pkg/gin"
)

type File struct {
	Base
	Bucket    string `json:"bucket" gorm:"type:varchar(128);not null;comment:oss bucket"`
	Mime      string `json:"mime" gorm:"type:varchar(32);not null;comment:文件类型"`
	Size      *int   `json:"size" gorm:"type:int unsigned;not null;comment:文件大小(字节)"`
	OssPath   string `json:"ossPath" gorm:"type:varchar(256);not null;comment:oss的路径"`
	Name      string `json:"name" gorm:"type:varchar(256);not null;comment:文件名"`
	Ext       string `json:"ext" gorm:"type:varchar(32);not null;comment:后缀名"`
	UploadUID int    `json:"uploadUid" gorm:"type:int unsigned;not null;default:0;comment:上传用户id"`
}

var (
	FileFilterNameMapDBField = map[string]string{
		"search": "id|name|excerpt|description",
		"id":     "id",
		"name":   "name",
		"ctime":  "ctime",
	}
	FileOrderKeyMap = map[string]string{
		"id":    "id",
		"ctime": "ctime",
	}
)

type DefaultSceneFile struct {
	ID      int        `json:"id"`
	Ctime   *time.Time `json:"ctime,omitempty" gorm:"default:datetime"`
	OssPath string     `json:"url"`
	Name    string     `json:"name"`
}

func NewDefaultSceneFile(m *File) *DefaultSceneFile {
	model := &DefaultSceneFile{
		ID:      m.ID,
		Ctime:   m.Ctime,
		OssPath: m.OssPath,
		Name:    m.Name,
	}
	return model
}
func (m *File) NewModel() BaseModel {
	return &File{
		Base: Base{Ctx: m.Ctx},
	}
}
func (m *File) GetFilterMap() map[string]string {
	return FileFilterNameMapDBField
}
func (m *File) GetOrderMap() map[string]string {
	return FileOrderKeyMap
}
func (m *File) DelWithPostData(data Any, _ *ginApp.App) (int, error) {
	d := data.(*fastcurd.DelData)
	return Delete(m, d.IDs)
}

func (m *File) GetDetail(id int) error {
	return Detail(m, id)
}
func (m *File) GetGormQuery() *gorm.DB {
	db := global.DB
	if m.Ctx != nil {
		db = db.WithContext(m.Ctx)
	}
	return db.Model(&File{}).Scopes(NotDelScope)
}
func (m *File) GormFindList(q *gorm.DB) Any {
	actModels := make([]*File, 0, 1)
	// nolint
	q = q.Find(&actModels)
	return &actModels
}

func (m *File) GetFmtList(list Any, scene string) Any {
	fmtList := make([]Any, 0, 1)
	actList := list.(*[]*File)
	for _, item := range *actList {
		fmtList = append(fmtList, item.GetFmtDetail(scene))
	}
	return fmtList
}
func (m *File) GetFmtDetail(scenes ...string) Any {
	var scene string
	if len(scenes) == 1 {
		scene = scenes[0]
	}
	var model Any
	// nolint
	switch scene {
	default:
		model = NewDefaultSceneFile(m)
	}
	return model
}
