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
func DispatchMess(ctx *context.Context) {
	zap.L().Info("start dispatch mess task...")

	cfg := ctx.GetAppConfig()
	TaskStorage.Tk = cfg.YM.Tk

	var cityList []*City
	for {
		time.Sleep(time.Second * 10)

		if len(cityList) < 1 {
			for _, region := range cfg.Regions {
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
			rows, err := GetAllDepartList(city.Value)
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
			time.Sleep(time.Millisecond * 500)
			zap.L().Debug("depart detail", zap.Any("data", depart))

			did := depart.DepaVaccId
			// 立即订阅
			if TaskStorage.IsSendDepart(did) {
				if depart.IsNotice == 0 && depart.Total > 0 {
					depart.IsNowSubscribe = true
					DepartChan <- depart
					continue
				}

				var err error
				if depart.SubScribeNum, err = GetSubscribeNum(did); err != nil {
					zap.L().Error("get subscribe num fail", zap.Error(err))
					continue
				}
				zap.L().Debug("get depart subscribe number", zap.Int64("depart_id", did), zap.Int64("num", depart.SubScribeNum))

				// 可以订阅
				if depart.SubScribeNum <= config.SubscribeAbleMaxNum && depart.StopSubscribe == 0 {
					DepartChan <- depart
				}
			}
		}
	}

}
