package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"notify/config"
	"notify/lib"
	"time"
)
var log *logrus.Entry
func main() {

	co := config.GetConfig()
	var eMsg string
	if co.DingToken == "" {
		eMsg = "配置文件钉钉token配置项不正确或配置文件缺失"
	}
	if co.RedisHost == "" {
		eMsg = "配置文件redis配置项不正确或缺失"
	}
	log = lib.GetLogInstance()
	if eMsg != "" {
		log.Fatal(eMsg)
	}
	go func() {
		log.Trace("程序启动")
	}()
	Notice(co.DingToken,co.DingSec,co.RedisHost,co.RedisPort)
}

func Notice(token, sec string, host, port string) {
	var addr string
	addr = host + ":" + port
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	key := "logQueue"
	content := "时间: " + time.Now().Format("2006-01-02 15:04:05") + "\n信息：程序启动通知"

	d := Ding{MsgTypeN: "text"}
	d.Text = make(map[string]string)
	d.Text["content"] = content
	notice(d, token, sec)

	pu := client.Subscribe(key)
	for msg := range pu.Channel() {
		ms := NoticeS{}
		e := json.Unmarshal([]byte(msg.Payload), &ms)
		log.Trace("channel：",msg.Channel)
		var content string
		if e != nil {
			content = "时间: " + time.Now().Format("2006-01-02 15:04:05") + "\n信息: json解析错误 " + GetIntranetIp() + "原始信息" + msg.Payload
		} else {
			if ms.Level < 0 {
				ms.Level = 0
			}
			mLevel := Level(ms.Level)
			content = "名称: " + ms.Name + "\n" + "时间: " + ms.Time + "\n级别: " + fmt.Sprintf("%s", mLevel) + "\n信息: " + ms.Message
		}
		content = content+"\nIP: "+GetIntranetIp()
		if ms.Type=="dingding" ||e !=nil ||ms.Type==""{
			d := Ding{MsgTypeN: "text"}
			d.Text = make(map[string]string)
			d.Text["content"] = content
			log.Error("content", content)
			notice(d, token, sec)
		}


	}

}