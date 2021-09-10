package task

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	DepartChan  chan *DepartRow
	TaskStorage *taskStorage
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// taskStorage _
type taskStorage struct {
	Tk       string
	Lock     sync.RWMutex
	DepartMp map[int64]int64
}

// Resource _
type Resource struct {
	Code  string      `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg,omitempty"`
	OK    bool        `json:"ok"`
	NotOK bool        `json:"notOk"`
}

func InitTask() {
	DepartChan = make(chan *DepartRow, 2<<4)
	TaskStorage = new(taskStorage)
	TaskStorage.DepartMp = make(map[int64]int64)
	TaskStorage.initData()
}

// getDepartDataPath _
func getDepartDataPath() (path string) {
	path, _ = os.Getwd()
	path += "/depart_data.json"
	return
}

// initData _
func (t *taskStorage) initData() {
	path := getDepartDataPath()
	var err error
	defer func() {
		if err != nil {
			zap.L().Error("load stored depart ids err", zap.Error(err))
			log.Panic(err)
		}
	}()

	var departIds []int64
	if _, err = os.Stat(path); os.IsNotExist(err) {
		f, fErr := os.Create(path)
		if fErr != nil {
			err = fErr
		}
		f.WriteString("[")
		err = nil

	} else {
		var (
			departList []*DepartRow
			data       []byte
		)
		if data, err = ioutil.ReadFile(path); err != nil {
			return
		}
		json.Unmarshal(data, &departList)

		for _, depart := range departList {
			departIds = append(departIds, depart.DepaVaccId)
		}
	}
	zap.L().Debug("load stored depart ids", zap.Int64s("data", departIds))
	return
}

// IsSendDepart _
func (t *taskStorage) IsSendDepart(did int64) (check bool) {
	t.Lock.RLock()
	timeStamp, ok := t.DepartMp[did]
	if ok {
		if time.Now().Unix()-timeStamp > 300 {
			check = false
		} else {
			check = true
		}

	} else {
		check = true
	}

	t.Lock.RUnlock()
	return
}

// AddDepartData 写入已发送的社区到json文件
func (t *taskStorage) AddDepartData(depart *DepartRow) (err error) {
	did := depart.DepaVaccId

	t.Lock.Lock()
	defer t.Lock.Unlock()

	tst, exist := t.DepartMp[did]
	zap.S().Debugf("depart_id: %d timestamp: %d", did, tst)

	if !exist {
		t.DepartMp[did] = time.Now().Unix()
	} else {
		return
	}

	path := getDepartDataPath()
	f, err := os.OpenFile(path, os.O_RDWR, 6)
	defer f.Close()
	if err != nil {
		return
	}

	var data []byte
	if data, err = json.Marshal(depart); err != nil {
		return
	}

	contByte, _ := ioutil.ReadFile(path)
	contStr := string(contByte)
	index := int64(strings.Index(contStr, "]"))

	var writeStr string
	if index < 0 {
		index = 1
		writeStr = fmt.Sprintf("%s]", string(data))
	} else {
		writeStr = fmt.Sprintf(",%s]", string(data))
	}
	f.Seek(index, 0)
	f.WriteString(writeStr)

	zap.L().Debug("storage depart", zap.Int64("id", did))
	return
}

// GetResource _
func (t taskStorage) GetResource(urlStr string, params map[string]string) (res *Resource, err error) {
	res = new(Resource)
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

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, res); err != nil {
		return
	}

	if res.NotOK {
		err = fmt.Errorf("get response fail err: %v", res.Msg)
	}
	return
}
