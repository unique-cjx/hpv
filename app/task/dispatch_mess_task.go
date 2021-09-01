package task

import (
	"go.uber.org/zap"
	"hpv/bootstrap/context"
	"hpv/config"
	"runtime/debug"
	"time"
)

type City struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// DispatchMess _
func DispatchMess(values ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			zap.S().Error("panic:", err)
		}
	}()

	zap.L().Info("start dispatch mess task...")
	ctx := values[0].(*context.Context)

	ymConf := ctx.GetAppConfig().YueMiao
	TaskStorage.Tk = ymConf.Tk

	var cityList []*City

	for {
		time.Sleep(time.Second * 5)

		if len(cityList) < 1 {
			resp, err := TaskStorage.GetResource(config.CityListUrl, map[string]string{"parentCode": ymConf.Province.Code})
			if err != nil {
				zap.L().Error("get city list failed", zap.Error(err))
				continue
			}
			zap.L().Debug("get city list", zap.Any("data", resp))

			respBytes, _ := json.Marshal(resp.Data)
			json.Unmarshal(respBytes, &cityList)
			continue
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
			time.Sleep(time.Second * 1)
			var err error
			depart.SubScribeNum, err = GetSubscribeNum(depart.DepaVaccId)
			if err != nil {
				zap.L().Error("get subscribe num fail", zap.Error(err))
				continue
			}
			zap.L().Debug("depart detail", zap.Any("data", depart))

			TaskStorage.DidLock.RLock()
			did := depart.DepaVaccId
			if depart.IsSeckill == 1 {
				num := TaskStorage.SeckillMp[did]
				if num < config.NoticeMaxNum {
					TaskStorage.SeckillMp[did] += 1
					DepartChan <- depart
					zap.L().Info("now depart id that can be subscribed", zap.Any("data", TaskStorage.SeckillMp))
				}
			} else if depart.SubScribeNum <= config.SubscribeAbleNum {
				for _, v := range TaskStorage.DepartIds {
					if v == did {
						goto Loop
					}
				}
				DepartChan <- depart
			}
		Loop:
			TaskStorage.DidLock.RUnlock()
		}
	}
}
