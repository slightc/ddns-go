package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"main/aliyun"
	"main/common"
	"main/iputil"
	"main/qcloud"

	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"
)

type ddnsConfig struct {
	Type      string `yaml:"Type"`
	SecretId  string `yaml:"SecretId"`
	SecretKey string `yaml:"SecretKey"`
	Domain    string `yaml:"Domain"`
	Record    string `yaml:"Record"`
	RecordId  string `yaml:"RecordId"`
	Cron      string `yaml:"Cron"`
}

func main() {
	var configFilePath string
	var setting ddnsConfig
	var showDebugInfo bool
	var showRecordList bool
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	flag.StringVar(&configFilePath, "c", "/etc/ddns-go/config.yaml", "config file path")
	flag.BoolVar(&showDebugInfo, "d", false, "show debug information")
	flag.BoolVar(&showRecordList, "l", false, "show record list")
	flag.Parse()

	config, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		configFilePath = "./config.yaml"
		fmt.Println("read config file err: ", err)
		fmt.Println("try use ", configFilePath)
		config, err = ioutil.ReadFile(configFilePath)
	}

	if err != nil {
		fmt.Println("read config file err: ", err)
		fmt.Println("exit with no config file")
		return
	}

	yaml.Unmarshal(config, &setting)

	var recordInfo = common.RecordInfo{Domain: setting.Domain, Name: setting.Record, Id: setting.RecordId}
	var accessKey = common.AccessKey{Id: setting.SecretId, Secret: setting.SecretKey}
	var recordHandler common.RecordHandler = nil

	fmt.Printf("%#v\n", recordInfo)

	if setting.Type == "aliyun" {
		recordHandler = aliyun.Aliyun{Key: accessKey}
	}
	if setting.Type == "qcloud" {
		recordHandler = qcloud.Qcloud{Key: accessKey}
	}

	if recordHandler == nil {
		fmt.Println("recordHandler create failed")
		return
	}

	err = iputil.Init()
	if err != nil {
		fmt.Println("create temp file err: ", err)
		return
	}
	defer iputil.Deinit()

	crontab := cron.New()
	task := func() {
		ip, err := iputil.GetIp()
		if err != nil {
			fmt.Println("get ip err: ", err)
			return
		}
		if showDebugInfo {
			fmt.Println("now ip: ", ip)
		}
		ipChanged, err := iputil.IpIsChanged(ip)
		if err != nil {
			fmt.Println("check ip changed err: ", err)
			return
		}
		if !ipChanged {
			if showDebugInfo {
				fmt.Println("ip not changed")
			}
			return
		}
		fmt.Println("== ", time.Now(), " ==")
		fmt.Println("new ip: ", ip)
		err = recordHandler.SetRecordIp(&recordInfo, ip)
		if err != nil {
			fmt.Println("set dns err: ", err)
			return
		}
		fmt.Println("ip modify success")
	}
	fmt.Println("start listen ip change")
	fmt.Println("cron run use: ", setting.Cron)
	crontab.AddFunc(setting.Cron, task)
	task()
	crontab.Start()

	fmt.Println("waiting")

	select {
	case s := <-exitSignal:
		fmt.Println("exit with :", s)
		break
	}
}
