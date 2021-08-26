package task

import (
	"go.uber.org/zap"
	"hpv/bootstrap/context"
	"hpv/config"
	"log"
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

	ymConf := ctx.GetAppConfig().YueMiao
	TaskStorage.Tk = ymConf.Tk

	var cityList []*City

	for {
		time.Sleep(time.Second * 5)

		if len(cityList) < 1 {
			resp, err := TaskStorage.GetResource(config.CityListUrl, map[string]string{"parentCode": ymConf.Province.Code})
			if err != nil {
				log.Panic("get city list fail")
			}
			zap.L().Debug("city list", zap.Any("data", resp))

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
			if depart.SubScribeNum <= config.SubscribeAbleNum {
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
