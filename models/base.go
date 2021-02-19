package models

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stoewer/go-strcase"
	"gorm.io/gorm"

	"github.com/real-web-world/go-web-api/pkg/fastcurd"
	"github.com/real-web-world/go-web-api/pkg/gin"
)

const (
	SceneDefault = "default"
	SceneAdmin   = "admin"
	SceneWithNav = "withNav"
	SceneProfile = "profile"
)

var (
	queryCondRequiredErr = errors.New("query condition required")
	implThisErr          = errors.New("must impl this")
)

type Any interface{}
type Base struct {
	ID                 int          `json:"id" gorm:"type:int unsigned auto_increment;primary_key;"`
	Ctime              *time.Time   `json:"ctime,omitempty" gorm:"type:datetime;default:CURRENT_TIMESTAMP;not null"`
	Utime              *time.Time   `json:"utime,omitempty" gorm:"type:datetime ON UPDATE CURRENT_TIMESTAMP;default:CURRENT_TIMESTAMP;not null;"`
	Dtime              *int64       `json:"dtime,omitempty" gorm:"type:int unsigned;default:0;not null"`
	RelationAffectRows int          `json:"-" gorm:"-"` // 更新时用来保存其他关联数据的更新数
	Ctx                *gin.Context `json:"-" gorm:"-"` // gin req ctx
}
type BaseModel interface {
	GetID() int
	GetCtime() *time.Time
	GetUtime() *time.Time
	GetDtime() *int64
	GetFilterMap() map[string]string
	GetOrderMap() map[string]string
	GetDetail(int) error
	GetFmtDetail(scenes ...string) Any
	NewAddData() Any
	NewEditData() Any
	AddWithPostData(d Any, app *ginApp.App) error
	EditWithPostData(d Any, app *ginApp.App) (int, error)
	DelWithPostData(d Any, app *ginApp.App) (int, error)
	GetGormQuery() *gorm.DB
	GormFindList(q *gorm.DB) Any
	GetFmtList(list Any, scene string) Any
	GetRelationAffectRows() int
	NewModel() BaseModel
	HasSortField() bool
	AfterDelAction(*ginApp.App)
	AfterEditAction(*ginApp.App)
	AfterDetailAction(*ginApp.App)
	BeforeFmtDetail(*ginApp.App)
	BeforeFmtList(list interface{}, _ *ginApp.App)
	SetCtx(ctx *gin.Context)
}

func (m *Base) NewModel() BaseModel {
	panic(implThisErr)
}
func (m *Base) SetCtx(ctx *gin.Context) {
	m.Ctx = ctx
}
func (m *Base) GetID() int {
	return m.ID
}
func (m *Base) GetCtime() *time.Time {
	return m.Ctime
}
func (m *Base) GetUtime() *time.Time {
	return m.Utime
}
func (m *Base) GetDtime() *int64 {
	return m.Dtime
}
func (m *Base) HasSortField() bool {
	return false
}
func (m *Base) GetRelationAffectRows() int {
	return m.RelationAffectRows
}
func (m *Base) GetDetail(_ int) error {
	return implThisErr
}
func (m *Base) AfterDelAction(_ *ginApp.App) {
}
func (m *Base) AfterEditAction(_ *ginApp.App) {
}
func (m *Base) AfterDetailAction(_ *ginApp.App) {
}
func (m *Base) BeforeFmtDetail(_ *ginApp.App) {
}
func (m *Base) BeforeFmtList(_ interface{}, _ *ginApp.App) {
}
func (m *Base) GetFilterMap() map[string]string {
	return map[string]string{
		"id":    "id",
		"ctime": "ctime",
	}
}
func (m *Base) GetOrderMap() map[string]string {
	return map[string]string{
		"id":    "id",
		"ctime": "ctime",
	}
}
func (m *Base) GetFmtDetail(_ ...string) Any {
	return m
}
func (m *Base) NewAddData() Any {
	return nil
}
func (m *Base) NewEditData() Any {
	return nil
}
func (m *Base) AddWithPostData(_ Any, _ *ginApp.App) error {
	panic(implThisErr)
}
func (m *Base) EditWithPostData(_ Any, _ *ginApp.App) (int, error) {
	return 0, nil
}
func (m *Base) DelWithPostData(_ Any, _ *ginApp.App) (int, error) {
	return 0, nil
}
func (m *Base) GetGormQuery() *gorm.DB {
	return nil
}
func (m *Base) GormFindList(_ *gorm.DB) Any {
	return nil
}
func (m *Base) GetFmtList(_ Any, _ string) Any {
	return nil
}
func Add(m BaseModel) error {
	return m.GetGormQuery().Create(m).Error
}
func Get(m BaseModel, args ...Any) error {
	q := m.GetGormQuery()
	if len(args) == 1 {
		// nolint
		switch args[0].(type) {
		case int:
			if m.GetID() == 0 {
				q = q.Where("id = ?", args[0])
			}
			return q.First(m).Error
		case fastcurd.Filter:
			filter, _ := args[0].(fastcurd.Filter)
			q = fastcurd.BuildFilterCond(m.GetFilterMap(), q, filter)
			return q.First(m).Error
		case *fastcurd.Filter:
			filter, _ := args[0].(*fastcurd.Filter)
			q = fastcurd.BuildFilterCond(m.GetFilterMap(), q, *filter)
			return q.First(m).Error
		}
	}
	if len(args) == 2 {
		key, ok := args[0].(string)
		if !ok {
			return errors.New("query condtion incorrect ,args[0] " +
				"should be string,act is " + key)
		}
		switch args[1].(type) {
		case int, string:
			key = strcase.SnakeCase(key)
			if !fastcurd.IsValidQueryField(key) {
				return errors.New("query field " + key + " incorrect,only allow a-z and _")
			}
			q = q.Where("`"+key+"` = ?", args[1])
			return q.First(m).Error
		default:
			return errors.New("query condition incorrect ,args[1] type incorrect")
		}
	}
	return nil
}
func NotDelScope(db *gorm.DB) *gorm.DB {
	return db.Where("dtime = 0")
}

// 软删除模型
// all sql should soft delete
// delete(model) ok
// delete(modelType,idList)
//
func Delete(m BaseModel, args ...Any) (int, error) {
	q := m.GetGormQuery()
	if len(args) == 0 {
		return 0, queryCondRequiredErr
	}
	switch actArg := args[0].(type) {
	case int:
		q = q.Where("id = ?", actArg)
	case []int:
		q = q.Where("id in (?)", actArg)
	case fastcurd.Filter:
		filter := actArg
		q = fastcurd.BuildFilterCond(m.GetFilterMap(), q, filter)
	case *fastcurd.Filter:
		filter := *actArg
		q = fastcurd.BuildFilterCond(m.GetFilterMap(), q, filter)
	default:
	}
	res := q.Update(fastcurd.DeleteTimeField, time.Now().Unix())
	if err := res.Error; err != nil {
		return 0, err
	}
	return int(res.RowsAffected), nil
}

func Detail(m BaseModel, args ...Any) error {
	q := m.GetGormQuery()
	if len(args) == 0 {
		return queryCondRequiredErr
	}
	switch acrArg := args[0].(type) {
	case int, *int:
		q = q.Where("id = ?", acrArg)
	case fastcurd.Filter:
		filter, _ := args[0].(fastcurd.Filter)
		q = fastcurd.BuildFilterCond(m.GetFilterMap(), q, filter)
	}
	// 通过主键查询第一条记录
	if err := q.First(m).Error; err != nil {
		return err
	}
	return nil
}

func Update(m BaseModel, args ...Any) (int, error) {
	q := m.GetGormQuery()
	if len(args) != 0 {
		switch actArg := args[0].(type) {
		case int:
			q = q.Where("id = ?", actArg)
		case fastcurd.Filter:
			filter, _ := args[0].(fastcurd.Filter)
			q = fastcurd.BuildFilterCond(m.GetFilterMap(), q, filter)
		}
	}
	q = q.Updates(m)
	if err := q.Error; err != nil {
		return 0, err
	}
	return int(q.RowsAffected), nil
}

func Count(m BaseModel, where map[string]interface{}) (int, error) {
	var count int64
	if err := m.GetGormQuery().Where(where).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func IsExist(m BaseModel, where map[string]interface{}) bool {
	count, _ := Count(m, where)
	return count > 0
}
func List(m BaseModel, d *fastcurd.ListData, field ...[]string) (intCount int, models interface{},
	err error) {
	count := int64(0)
	filter := d.Filter
	page := d.Page
	limit := d.Limit
	order := d.Order
	q := m.GetGormQuery()
	if len(field) == 1 {
		q = q.Select(field)
	}
	q = fastcurd.BuildFilterCond(m.GetFilterMap(), q, filter)
	q = fastcurd.BuildOrderCond(m.GetOrderMap(), q, order)
	offset := 0
	if page > 1 {
		offset = (page - 1) * limit
	}
	if limit != fastcurd.NotLimit {
		q = q.Offset(offset).Limit(limit)
	}
	models = m.GormFindList(q)
	q = q.Offset(-1).Limit(-1).Count(&count)
	err = q.Error
	return int(count), models, err
}
