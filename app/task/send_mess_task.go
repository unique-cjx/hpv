package task

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"hpv/app/util"
	"hpv/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
)

// SendMess _
func SendMess(values ...interface{}) {
	zap.L().Info("start send mess task...")

	pwd, _ := os.Getwd()
	imgPath := pwd + "/img"
	sendGroupMess := "[CQ:at,qq=all] \n %s 打开微信扫描二维码 \n [CQ:image,file=file://%s]"
	body := map[string]interface{}{
		"group_id":    config.QQGroupID,
		"message":     "",
		"auto_escape": false,
	}

	for {
		depart := <-DepartChan

		if err := depart.SetDepartPrompt(); err != nil {
			zap.S().Error(err)
			continue
		}

		var text = [5]interface{}{}
		if depart.IsNowSubscribe {
			text[0] = "检测到可立即预约的社区"
		} else {
			text[0] = "检测到订阅量较少的社区"
		}
		text[1] = depart.Name
		text[2] = "订阅人数：" + util.ToString(depart.SubScribeNum)
		text[3] = "社区电话：" + depart.Tel
		text[4] = "注意事项：" + depart.Prompt

		var buffer strings.Builder
		for _, s := range text {
			buffer.WriteString(fmt.Sprintf("- %v \n", s))
		}

		params := url.Values{}
		params.Add("vaccCode", depart.VaccineCode) // 8803 == 九价疫苗
		params.Add("depaCode", depart.Code)
		params.Add("vaccId", util.ToString(depart.DepaVaccId))
		params.Add("t", "1629365360744")

		rawQuery := params.Encode()
		url := fmt.Sprintf("%s?%s", config.DepartPageUrl, rawQuery)

		qrimg, err := qrcode.New(url, qrcode.Medium)
		if err != nil {
			zap.L().Error("failed to generate QR code", zap.Error(err))
			continue
		}
		qrPath := fmt.Sprintf("%s/%v-%v.png", imgPath, depart.RegionCode, depart.DepaVaccId)
		if err = qrimg.WriteFile(256, qrPath); err != nil {
			zap.L().Error("QR img write file path fail", zap.Error(err))
			continue
		}
		body["message"] = fmt.Sprintf(sendGroupMess, buffer.String(), qrPath)

		zap.L().Debug("qq-cq server request param", zap.Any("param", body))

		client := &http.Client{Timeout: time.Second * 5}
		bytesData, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", config.QQBotServ+"send_group_msg", bytes.NewReader(bytesData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			zap.L().Error("request qq-cq server fail", zap.Error(err))
			continue
		}
		respBytes, _ := ioutil.ReadAll(resp.Body)

		zap.L().Debug("qq-cq server resp body", zap.String("body", string(respBytes)))

		if err = TaskStorage.AddDepartData(depart); err != nil {
			zap.S().Error(err)
		}

		time.Sleep(time.Second * 5)
	}
}
