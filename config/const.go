package config

const (

	// CityListUrl 城市列表
	CityListUrl = "https://wx.scmttec.com/base/region/childRegions.do"

	// DepartListUrl 社区列表
	DepartListUrl = "https://wx.scmttec.com/base/department/getDepartments.do"

	// DepartDetailUrl 社区详情
	DepartDetailUrl = "https://wx.scmttec.com/base/departmentVaccine/item.do"

	// CountSubscribeUrl 订阅人数
	CountSubscribeUrl = "https://wx.scmttec.com/passport/register/countSubscribe.do"

	// DepartPageUrl 社区页面
	DepartPageUrl = "https://wx.scmttec.com/index.html#/vaccines"

	// QQBotServ qq机器人服务器地址
	QQBotServ = "http://127.0.0.1:5000/"

	// RefreshWxToken 刷新tk地址
	RefreshWxToken = "https://wx.scmttec.com/passport/wx/login.do"

	// QQGroupID qq群号
	QQGroupID = 981907686

	// SubscribeAbleNum 订阅人数阈值
	SubscribeAbleNum int64 = 2000

	// NoticeMaxNum 通知最大次数
	NoticeMaxNum int8 = 3
)
