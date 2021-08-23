package config

const (

	// CityListUrl 城市列表
	CityListUrl = "https://wx.scmttec.com/base/region/childRegions.do"

	// DepartmentsUrl 门诊列表
	DepartmentsUrl = "https://wx.scmttec.com/base/department/getDepartments.do"

	// SubscribeUrl 订阅
	SubscribeUrl = "https://wx.scmttec.com/passport/register/countSubscribe.do"

	// DetailVoUrl 门诊详情
	DetailVoUrl = "https://wx.scmttec.com/index.html#/vaccines"

	// QQBotServ qq机器人服务器
	QQBotServ = "http://127.0.0.1:5000/"

	// RefreshWxToken 刷新
	RefreshWxToken = "https://wx.scmttec.com/passport/wx/login.do"

	// QQGroupID qq群号
	QQGroupID = 981907686

	// SubscribeAbleNum 订阅人数阈值
	SubscribeAbleNum int64 = 2500
)
