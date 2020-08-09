package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
	"runtime"
	"sync"
)

type Conf struct {
	RedisHost       string            `yaml:"redis_host"`
	RedisPort       string            `yaml:"redis_port"`
	Mail            map[string]string `yaml:"mail"`
	MailPort        int               `yaml:"mail_port"`
	DingToken       string            `yaml:"ding_token"`
	DingSec       string            `yaml:"ding_sec"`
}

var on sync.Once
var co *Conf

func GetConfig() *Conf  {
	if co == nil {
		config := Conf{}
		config.getConf()
		co = &config
	}
	return co
}
func  (c *Conf)getConf(){
		on.Do(func() {
			yamlFile, err := ioutil.ReadFile("./config/config.yaml")
			//开发模式使用
			if err != nil {
				yamlFile, err = ioutil.ReadFile(getCurrentPath() + "/config.yaml")
				if err != nil {
					log.Printf("yamlFile02.Get err   #%v ", err)
				}
			}
			err = yaml.Unmarshal(yamlFile, c)
			if err != nil {
				log.Fatalf("Unmarshal: %v", err)
			}
		})
	//编译前使用， 编译后的可执行文件getCurrentPath() 返回为开发模式下的目录

}

func getCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
