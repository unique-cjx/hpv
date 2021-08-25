package task

import (
	"go.uber.org/zap"
	"hpv/config"
)

type DepartRows struct {
	DepaVaccId    int64  `json:"depaVaccId"`
	VaccineCode   string `json:"vaccineCode"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	RegionCode    string `json:"regionCode"`
	Tel           string `json:"tel"`
	IsOpen        int8   `json:"isOpen"`
	Address       string `json:"address"`
	WorktimeDesc  string `json:"worktimeDesc"`
	Total         int    `json:"total"`
	SubScribeNum  int64  `json:"subscribeNum"` // 订阅人数
	IsSeckill     int8   `json:"isSeckill"`
	StopSubscribe int8   `json:"stopSubscribe"`
}

type DepartmentsResp struct {
	Offset       int          `json:"offset"`
	End          int          `json:"end"`
	Total        int          `json:"total"`
	Limit        int          `json:"limit"`
	PageNumber   int          `json:"pageNumber"`
	PageListSize int          `json:"pageListSize"`
	PageNumList  []int        `json:"pageNumList"`
	DepartRows   []DepartRows `json:"rows"`
	Pages        int          `json:"pages"`
}

// GetActiveDepartList 获取可订阅的社区列表
func GetActiveDepartList(regionCode string) (rows []*DepartRows, err error) {
	param := map[string]string{
		"offset":     "0",
		"limit":      "80",
		"regionCode": regionCode,
		"sortType":   "1",
		"isOpen":     "1",
		"customId":   "3",
	}
	resp, err := TaskStorage.GetResource(config.DepartmentsUrl, param)
	if err != nil {
		zap.L().Error("get departments error", zap.Error(err))
		return
	}
	departResp := new(DepartmentsResp)
	vByte, _ := json.Marshal(resp.Data)
	json.Unmarshal(vByte, departResp)

	for _, row := range departResp.DepartRows {
		if row.StopSubscribe == 0 {
			row := row
			rows = append(rows, &row)
		}
	}
	return
}
