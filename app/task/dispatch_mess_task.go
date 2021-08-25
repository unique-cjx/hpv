package task

import (
	"errors"
	"go.uber.org/zap"
	"hpv/app/util"
	"hpv/bootstrap/context"
	"hpv/config"
	"log"
	"sync"
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
	wg := values[1].(*sync.WaitGroup)

	ymConf := ctx.GetAppConfig().YueMiao
	TaskStorage.Tk = ymConf.Tk

	resp, err := TaskStorage.GetResource(config.CityListUrl, map[string]string{"parentCode": ymConf.Province.Code})
	if err != nil {
		log.Panic("get city list fail")
	}
	zap.L().Debug("city list", zap.Any("data", resp))

	var cityList []City
	respBytes, _ := json.Marshal(resp.Data)
	json.Unmarshal(respBytes, &cityList)

	tick := time.NewTicker(time.Second * 8)
	for {
		<-tick.C

		var departList []*DepartRows
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
			depart.SubScribeNum, err = GetSubscribeNum(depart.DepaVaccId)
			if err != nil {
				zap.L().Error("get subscribe num fail", zap.Error(err))
				continue
			}
			zap.L().Debug("depart detail", zap.Any("data", depart))

			TaskStorage.Lock.RLock()
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
			TaskStorage.Lock.RUnlock()
		}
	}
	wg.Done()
}

// GetSubscribeNum 获取指定社区的订阅人数
func GetSubscribeNum(id int64) (data int64, err error) {
	params := map[string]string{"depaVaccId": util.ToString(id)}
	resp, err := TaskStorage.GetResource(config.SubscribeUrl, params)
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
