package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type NoticeS struct {
	Name    string `json:"name"`
	Level   int    `json:"level"` //0 1 2
	Time    string `json:"time"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type Level uint32

const (
	FatalLevel Level = iota
	WarnLevel
	MsgLevel
)

var errCodeMap = map[Level]string{
	FatalLevel: "程序错误，需立即解决！",
	WarnLevel:  "程序错误，稍后解决！",
	MsgLevel:   "提示信息",
}

func (e Level) String() string {
	if v, ok := errCodeMap[e]; ok {
		return v
	}
	return "未找到码，请检查错误码"
}

type Ding struct {
	MsgTypeN string            `json:"msgtype"`
	Text     map[string]string `json:"text"`
}



func notice(d Ding, token, sec string) {
	robotAddr, err := GetDingTalkURL(token, sec)
	if err != nil {
		log.Error(err, "url", robotAddr)
	}
	dJson, _ := json.Marshal(d)
	payload := strings.NewReader(string(dJson))
	req, _ := http.NewRequest("POST", robotAddr, payload)
	req.Header.Add("content-type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	log.Debug(string(body))
}

const dingTalkOAPI = "oapi.dingtalk.com"

var dingTalkURL url.URL = url.URL{
	Scheme: "https",
	Host:   dingTalkOAPI,
	Path:   "robot/send",
}

// GetDingTalkURL get DingTalk URL with accessToken & secret
// If no signature is set, the secret is set to ""
// 如果没有加签，secret 设置为 "" 即可
func GetDingTalkURL(accessToken string, secret string) (string, error) {
	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
	return GetDingTalkURLWithTimestamp(timestamp, accessToken, secret)
}

// GetDingTalkURLWithTimestamp get DingTalk URL with timestamp & accessToken & secret
func GetDingTalkURLWithTimestamp(timestamp string, accessToken string, secret string) (string, error) {
	dtu := dingTalkURL
	value := url.Values{}
	value.Set("access_token", accessToken)

	if secret == "" {
		dtu.RawQuery = value.Encode()
		return dtu.String(), nil
	}

	sign, err := sign(timestamp, secret)
	if err != nil {
		dtu.RawQuery = value.Encode()
		return dtu.String(), err
	}

	value.Set("timestamp", timestamp)
	value.Set("sign", sign)
	dtu.RawQuery = value.Encode()
	return dtu.String(), nil
}

func sign(timestamp string, secret string) (string, error) {
	stringToSign := fmt.Sprintf("%s\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	if _, err := io.WriteString(h, stringToSign); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func GetIntranetIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "获取ip错误 error" + err.Error()
	}
	for _, address := range addrs { // 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "获取ip错误 error"
}
