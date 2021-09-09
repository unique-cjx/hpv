package task

import (
	"go.uber.org/zap"
	"hpv/app/util"
	"hpv/bootstrap/context"
	"hpv/config"
	"time"
)

type City struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// DispatchMess _
func DispatchMess(values ...interface{}) {
	zap.L().Info("start dispatch mess task...")
	ctx := values[0].(*context.Context)

	cfg := ctx.GetAppConfig()
	TaskStorage.Tk = cfg.YM.Tk

	var cityList []*City
	for {
		time.Sleep(time.Second * 5)

		if len(cityList) < 1 {
			for _, region := range cfg.Region {
				code := util.ToString(region.Code)
				resp, err := TaskStorage.GetResource(config.CityListUrl, map[string]string{"parentCode": code})
				if err != nil {
					zap.L().Error("get city list failed", zap.Error(err))
					continue
				}
				zap.L().Debug("get city list", zap.Any("data", resp))

				respBytes, _ := json.Marshal(resp.Data)
				json.Unmarshal(respBytes, &cityList)
				continue
			}
		}

		var departList []*DepartRow
		for _, city := range cityList {
			rows, err := GetActiveDepartList(city.Value)
			if err != nil {
				zap.L().Error("get depart list fail", zap.Error(err))
				continue
			}
			if len(rows) < 1 {
				continue
			}
			for _, row := range rows {
				departList = append(departList, row)
			}
		}
		for _, depart := range departList {
			var err error
			depart.SubScribeNum, err = GetSubscribeNum(depart.DepaVaccId)
			if err != nil {
				zap.L().Error("get subscribe num fail", zap.Error(err))
				continue
			}
			zap.L().Debug("depart detail", zap.Any("data", depart))

			if depart.SubScribeNum <= config.SubscribeAbleMaxNum {
				if depart.IsNotice == 0 && depart.Total > 0 {
					depart.IsNowSubscribe = true
				}
				DepartChan <- depart
			}
			time.Sleep(time.Second * 1)
		}
	}
}
