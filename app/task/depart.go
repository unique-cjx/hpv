package task

import (
	"errors"
	"go.uber.org/zap"
	"hpv/app/util"
	"hpv/config"
)

type DepartRow struct {
	DepaVaccId     int64  `json:"depaVaccId"`
	VaccineCode    string `json:"vaccineCode"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	RegionCode     string `json:"regionCode"`
	Tel            string `json:"tel"`
	IsOpen         int8   `json:"isOpen"`
	Total          int    `json:"total"`
	SubScribeNum   int64  `json:"subscribeNum,omitempty"` // 订阅人数
	StopSubscribe  int8   `json:"stopSubscribe"`
	Prompt         string `json:"prompt"` // 注意事项
	IsNowSubscribe bool   `json:"-"`      // 是否可以立即订阅
	IsNotice       int8   `json:"isNoticedUserAllowed"`
}

type DepartmentsResp struct {
	Offset       int         `json:"offset"`
	End          int         `json:"end"`
	Total        int         `json:"total"`
	Limit        int         `json:"limit"`
	PageNumber   int         `json:"pageNumber"`
	PageListSize int         `json:"pageListSize"`
	PageNumList  []int       `json:"pageNumList"`
	DepartRow    []DepartRow `json:"rows"`
	Pages        int         `json:"pages"`
}

// GetActiveDepartList 获取可订阅的社区列表
func GetActiveDepartList(regionCode string) (rows []*DepartRow, err error) {
	param := map[string]string{
		"offset":     "0",
		"limit":      "80",
		"regionCode": regionCode,
		"sortType":   "1",
		"isOpen":     "1",
		"customId":   "3", // 九价疫苗编号
	}
	resp, err := TaskStorage.GetResource(config.DepartListUrl, param)
	if err != nil {
		zap.L().Error("get depart list error", zap.Error(err))
		return
	}
	departResp := new(DepartmentsResp)
	vByte, _ := json.Marshal(resp.Data)
	json.Unmarshal(vByte, departResp)

	for _, row := range departResp.DepartRow {
		if row.StopSubscribe == 0 && row.DepaVaccId != 0 {
			row := row
			rows = append(rows, &row)
		}
	}
	return
}

// GetSubscribeNum 获取指定社区的订阅人数
func GetSubscribeNum(id int64) (data int64, err error) {
	params := map[string]string{"depaVaccId": util.ToString(id)}
	resp, err := TaskStorage.GetResource(config.CountSubscribeUrl, params)
	if err != nil {
		return
	}
	if resp.Data == nil {
		err = errors.New("temporarily unable to obtain data")
		return
	}
	tmpData, _ := util.ToInt(resp.Data)
	data = int64(tmpData)
	return
}

// SetDepartPrompt _
func (depar *DepartRow) SetDepartPrompt() (err error) {
	s := util.ToString(depar.DepaVaccId)
	resp, err := TaskStorage.GetResource(config.DepartDetailUrl, map[string]string{"id": s})
	if err != nil {
		zap.L().Error("get depart detail err", zap.Error(err))
		return
	}
	data := resp.Data.(map[string]interface{})
	prompt, ok := data["prompt"].(string)
	if !ok {
		depar.Prompt = ""
		return
	}
	depar.Prompt = prompt

	return
}
