package task

import (
	"fmt"
	"go.uber.org/zap"
	"hpv/config"
	"net/http"
	"sync"
	"time"
)

const tokenPrefix = "_xzkj_"

// RefreshToken _
func RefreshToken(values ...interface{}) {
	wg := values[0].(*sync.WaitGroup)
	tick := time.NewTicker(time.Minute * 30)

	for {
		<-tick.C

		zap.S().Info("start refresh token task...")
		// 禁止重定向
		client := &http.Client{Timeout: time.Second * 5, CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}
		req, _ := http.NewRequest("GET", config.RefreshWxToken, nil)
		cookie := fmt.Sprintf("%s=%s", tokenPrefix, TaskStorage.Tk)
		req.Header.Set("cookie", cookie)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_16) AppleWebKit/605.1.15 (KHTML, like Gecko) MicroMessenger/6.8.0(0x16080000) MacWechat/3.0.1(0x13000110) NetType/WIFI WindowsWechat")

		resp, err := client.Do(req)
		resp.Body.Close()

		if err != nil {
			zap.L().Error("refresh token err", zap.Error(err))
			goto Loop
		}

		var respCookies []*http.Cookie
		respCookies = resp.Cookies()
		for _, ck := range respCookies {
			if ck.Name == tokenPrefix {
				TaskStorage.Tk = ck.Value
			}
		}
		zap.L().Debug("refresh wxtoken", zap.String("data", TaskStorage.Tk))
	}

Loop:
	wg.Done()
}
