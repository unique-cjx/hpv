package task

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	DepartChan  chan *DepartRows
	TaskStorage *taskStorage
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// taskStorage _
type taskStorage struct {
	DidLock   sync.RWMutex
	DepartIds []int64
	Tk        string
}

// Resource _
type Resource struct {
	Code  string      `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg,omitempty"`
	OK    bool        `json:"ok"`
	NotOK bool        `json:"notOk"`
}

func init() {
	DepartChan = make(chan *DepartRows, 2<<4)
	TaskStorage = &taskStorage{DepartIds: make([]int64, 0)}
}

// GetResource _
func (t taskStorage) GetResource(urlStr string, params map[string]string) (res *Resource, err error) {
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
	client := &http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("GET", path, nil)

	req.Header.Set("tk", t.Tk)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_16) AppleWebKit/605.1.15 (KHTML, like Gecko) MicroMessenger/6.8.0(0x16080000) MacWechat/3.0.1(0x13000110) NetType/WIFI WindowsWechat")
	zap.L().Debug("get req", zap.String("url", path), zap.String("tk", t.Tk))

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	res = new(Resource)
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, res)

	if res.NotOK {
		err = fmt.Errorf("get response fail err: %v", res.Msg)
	}

	return
}
