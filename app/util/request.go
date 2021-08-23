package util

import (
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Resource struct {
	Code  string      `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg,omitempty"`
	OK    bool        `json:"ok"`
	NotOK bool        `json:"notOk"`
}

// GetResp _
func GetResp(urlStr string, params map[string]string, tk string) (res *Resource, err error) {
	request := url.Values{}
	Url, err := url.Parse(urlStr)
	if err != nil {
		return
	}
	for k, param := range params {
		request.Set(k, param)
	}
	Url.RawQuery = request.Encode()
	path := Url.String()
	zap.S().Debug("get api url: ", path)
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("GET", path, nil)
	if tk != "" {
		req.Header.Set("tk", tk)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_16) AppleWebKit/605.1.15 (KHTML, like Gecko) MicroMessenger/6.8.0(0x16080000) MacWechat/3.0.1(0x13000110) NetType/WIFI WindowsWechat")
	}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	res = new(Resource)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, res)
	zap.S().Debug("resp content: ", string(body))

	if res.NotOK {
		err = fmt.Errorf("get response fail err: %v", res.Msg)
	}

	return
}
