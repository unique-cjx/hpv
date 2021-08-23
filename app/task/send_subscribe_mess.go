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

// GetActiveRegions _
func GetActiveRegions(values ...interface{}) {
	zap.L().Info("start subscribe task...")
	ctx := values[0].(*context.Context)
	wg := values[1].(*sync.WaitGroup)

	ymConf := ctx.GetAppConfig().YueMiao

	resp, err := util.GetResp(config.CityListUrl, map[string]string{"parentCode": ymConf.Province.Code}, "")
	if err != nil {
		log.Panic("get city list fail")
	}

	var cityList []City
	respBytes, _ := json.Marshal(resp.Data)
	json.Unmarshal(respBytes, &cityList)

	tick := time.NewTicker(time.Second * 3)
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
			depart.SubScribeNum, err = GetSubscribeDepartNum(depart.DepaVaccId, ymConf.Tk)
			if err != nil {
				zap.L().Error("get subscribe num fail", zap.Error(err))
				continue
			}
			DepartStorage.Lock.RLock()
			did := depart.DepaVaccId
			if depart.SubScribeNum <= config.SubscribeAbleNum {
				for _, v := range DepartStorage.Dids {
					if v == did {
						goto Loop
					}
				}
				DepartChan <- depart
			}
		Loop:
			DepartStorage.Lock.RUnlock()
		}

	}
	wg.Done()
}

// GetSubscribeDepartNum _
func GetSubscribeDepartNum(id int64, tk string) (data int64, err error) {
	params := map[string]string{"depaVaccId": util.ToString(id)}
	resp, err := util.GetResp(config.SubscribeUrl, params, tk)
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
