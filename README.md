# go api project

基于gin+gorm的快速curd框架

[网站](https://blog.buffge.com)

## 功能

- fastcurd
- logger
- 数据库自动迁移
- docker-compose 部署
- sentry 错误收集
- prometheus 指标收集
- grafana 系统指标可视化
- 验证器多语言
- 权限管理
- json api
- 调试模式 返回sql详请
- go linter
- swagger 文档
- 环境变量载入配置
- 记录请求信息,根据reqID记录所有信息
- Jenkins 自动构建部署
- 记录每个请求的sql信息,请求完毕自动写入文件
- 验证码
- oss上传文件
- redis 缓存
- 封装gin 简单易用

## 中间件

- Token Auth
- Cache
- 跨域
- http trace 保存 req resp
- recover

## 部署流程

```shell

```

[comment]: <> (## Todo)

## 计划
- 调整部分slice默认容量 把之前的make(xx,0,1)改掉
- 更新关联项 使用批量insert update 见article update picIDList tagIDList
- app,string,byte 等常用碎片 对象池，
- react seo
- 根据model 自动生成前端ts文件
- firebase第三方登录
- ws聊天室
- 谷歌人机验证
- 发文章自动提交 百度 谷歌

### 更新优先级

- 写出安全的可用代码
- 优化速度
- 优化sql  
- 优化内存
- 优化代码数量

## fastcurd简介

模型：

每个模型都有 id，ctime（创建时间），utime（更新时间），dtime（删除时间） 4个字段

请求均使用post方法，每个接口模型均有5个常用接口
请求路径： /model/operation

- add 添加一条记录
- edit 编辑一条记录
- detail 获取一条记录
- list 获取一组记录
- del 软删除一组记录

其他定制接口均需要自己编写，若常用接口内有非登录权限问题（如管理员与其他用户操作不同）
则需要重新编写
安全性：
在路由定义中加入授权中间件 比如user模型 未登录5个接口都不能访问，普通用户只能访问detail接口
管理员可以访问所有。
返回值：

```
// 通用返回json
// 所有的接口均返回此对象
type RetJSON struct {
    Code  int            `json:"code" example:"0"`
    Data  interface{}    `json:"data,omitempty"`
    Msg   string         `json:"msg,omitempty" example:"提示信息"`
    Count *int           `json:"count,omitempty"`
    Page  int            `json:"page,omitempty"`
    Limit int            `json:"limit,omitempty"`
    Extra *RespJsonExtra `json:"extra,omitempty"`
}

type RespJsonExtra struct {
    ReqID    string      `json:"requestID"`
    SQLs     interface{} `json:"sqls,omitempty"`
    ProcTime string      `json:"procTime" example:"0.2s"`
    TempData interface{} `json:"tempData,omitempty"`
}
```

### add

请求参数 AddData 一般为模型中的字段，如有关联 如预览图列表 则有 profilePicIDList字段

### edit

请求参数 EditData 一般为 id + AddData 并将AddData必填字段设置为非必填

### detail

请求参数 id+scene 根据不同scene返回不同json 不传则返回 defaultScene，

admin scene 需要admin权限 如 请求用户详情 default scene返回信息较多，

profile scene 只返回id 昵称 头像 介绍4个字段

### list

请求参数 
- page: 页码
- limit： 每页数量 
- filter：查询条件（db字段）
- order：排序条件（db字段）
- extra： 其他参数 定制化
```
ListData struct {
    Page   int                    `json:"page" binding:"omitempty,required,min=0"`
    Limit  int                    `json:"limit" binding:"omitempty,required,min=0,max=50"`
    Filter Filter                 `json:"filter" binding:""`
    Order  map[string]string      `json:"order" binding:""`
    Extra  map[string]interface{} `json:"extra" binding:""`
}

type (
	Filter     = map[string]FilterItem
	FilterCond string
	FilterItem struct {
		Condition FilterCond
		Val       interface{}
	}
)

const (
	// 筛选条件
	CondUndefined FilterCond = "undefined"
	// 数值
	CondEq           FilterCond = "eq"
	CondLt           FilterCond = "lt"
	CondElt          FilterCond = "elt"
	CondGt           FilterCond = "gt"
	CondEgt          FilterCond = "egt"
	CondNeq          FilterCond = "neq"
	CondBetweenValue FilterCond = "betweenValue"
	// 字符串
	CondEqString  FilterCond = "eqString"
	CondLike      FilterCond = "like"
	// ... other condition
	)
```
### del
请求参数： ids id列表。 根据id删除


## 如何快速创建一个模型的curd
// todo b站做个视频

- 编辑并运行 cmd/generate_model/main.go 编辑需要添加的模型
会在models目录下 model.go文件
  编辑模型
  编辑文件中的 FilterNameMapDBField 此字段为查询条件对db字段的映射
  编辑文件中的 OrderKeyMap  此字段为排序字段对db字段的映射
  修改AddData EditData
  编辑各种scene模型和构造函数 如DefaultSceneModel AdminSceneModel
- 编辑并运行 数据库迁移
- 添加路由

