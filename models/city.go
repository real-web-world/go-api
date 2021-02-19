package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/real-web-world/go-web-api/global"
)

const (
	LevelProvince = "province"
	LevelCity     = "city"
	LevelCounty   = "county"
)

var (
	CityFilterNameMapDBField = map[string]string{
		"search": "id|year|description",
		"id":     "id",
		"level":  "level",
		"cid":    "code",
		"code":   "code",
		"ctime":  "ctime",
	}
	CityOrderKeyMap = map[string]string{
		"id":    "id",
		"ctime": "ctime",
	}
)

// 城市
type City struct {
	Base
	Level        string `json:"level" gorm:"type:varchar(8);not null;comment:层级"`
	Code         string `json:"code" gorm:"type:varchar(16);not null;index:code;comment:区域代码"`
	Name         string `json:"name" gorm:"type:varchar(24);not null;comment:区域名称"`
	Pinyin       string `json:"pinyin" gorm:"type:varchar(128);not null;comment:区域拼音"`
	PinyinAbbr   string `json:"pinyinAbbr" gorm:"type:varchar(32);not null;comment:拼音缩写"`
	ProvinceCode string `json:"provinceCode" gorm:"type:varchar(32);not null;index:province_code;comment:省代码"`
	CityCode     string `json:"cityCode" gorm:"type:varchar(32);not null;index:city_code;comment:市代码"`
	CountyCode   string `json:"countyCode" gorm:"type:varchar(32);not null;index:county_code;comment:区县代码"`
}
type DefaultSceneCity struct {
	CID      string              `json:"cid"`
	AreaName string              `json:"areaName"`
	Level    string              `json:"level"`
	Children []*DefaultSceneCity `json:"children"`
}

type SingleSceneCity struct {
	ID           int        `json:"id"`
	Level        string     `json:"level"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	Pinyin       string     `json:"pinyin"`
	PinyinAbbr   string     `json:"pinyinAbbr"`
	ProvinceCode string     `json:"provinceCode"`
	CityCode     string     `json:"cityCode"`
	CountyCode   string     `json:"countyCode"`
	Ctime        *time.Time `json:"ctime,omitempty" gorm:"default:datetime"`
	Utime        *time.Time `json:"utime,omitempty" gorm:"default:datetime"`
}

func NewDefaultSceneCity(m *City) *DefaultSceneCity {
	model := &DefaultSceneCity{
		CID:      m.Code,
		Level:    m.Level,
		AreaName: m.Name,
	}
	model.Children = m.GetFmtChildRegion()
	return model
}
func (m *City) GetFmtChildRegion() []*DefaultSceneCity {
	fmtChildRegionList := make([]*DefaultSceneCity, 0, 1)
	cityList := make([]*City, 0, 1)
	cityModel := &City{
		Base: Base{Ctx: m.Ctx},
	}
	q := cityModel.GetGormQuery()
	switch m.Level {
	case LevelProvince:
		q = q.Where("level = ? and province_code = ?", LevelCity, m.Code)
	case LevelCity:
		q = q.Where("level = ? and city_code = ?", LevelCounty, m.Code)
	default:
		return nil
	}
	q.Find(&cityList)
	if len(cityList) == 0 {
		return nil
	}
	for _, city := range cityList {
		fmtChildRegionList = append(fmtChildRegionList, NewDefaultSceneCity(city))
	}
	return fmtChildRegionList
}
func (m *City) NewModel() BaseModel {
	return &City{
		Base: Base{Ctx: m.Ctx},
	}
}
func (m *City) GetFilterMap() map[string]string {
	return CityFilterNameMapDBField
}
func (m *City) GetOrderMap() map[string]string {
	return CityOrderKeyMap
}

func (m *City) GetDetail(id int) error {
	return Detail(m, id)
}
func (m *City) GetGormQuery() *gorm.DB {
	db := global.DB
	if m.Ctx != nil {
		db = db.WithContext(m.Ctx)
	}
	return db.Model(&City{}).Scopes(NotDelScope)
}
func (m *City) GormFindList(q *gorm.DB) Any {
	actModels := make([]*City, 0, 1)
	// nolint
	q = q.Find(&actModels)
	return &actModels
}

func (m *City) GetFmtList(list Any, scene string) Any {
	fmtList := make([]Any, 0, 1)
	actList := list.(*[]*City)
	for _, item := range *actList {
		fmtList = append(fmtList, item.GetFmtDetail(scene))
	}
	return fmtList
}
func (m *City) GetFmtDetail(scenes ...string) Any {
	var scene string
	if len(scenes) == 1 {
		scene = scenes[0]
	}
	var model Any
	// nolint
	switch scene {
	default:
		model = NewDefaultSceneCity(m)
	}
	return model
}
