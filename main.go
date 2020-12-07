package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"main/qcloud"

	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"
)

type ddnsConfig struct {
	SecretId  string `yaml:"SecretId"`
	SecretKey string `yaml:"SecretKey"`
	Domain    string `yaml:"Domain"`
	Record    string `yaml:"Record"`
	Id        string `yaml:"Id"`
	Cron      string `yaml:"Cron"`
}

type qcloudStatus struct {
	Code     int    `json:"code"`
	CodeDesc string `json:"codeDesc"`
	Message  string `json:"message"`
}

func setRecordIp(config ddnsConfig, ip string) (qcloudStatus, error) {
	sKI := qcloud.SecretData{}
	sKI.SecretId = config.SecretId
	sKI.SecretKey = config.SecretKey
	record := qcloud.RecordData{}
	record.RecordId = config.Id
	record.SubDomain = config.Record
	record.RecordType = "A"
	record.Value = ip

	data, err := qcloud.ModifyRecord(sKI, config.Domain, record)

	if err != nil {
		return qcloudStatus{}, err
	}
	var status qcloudStatus
	err = json.Unmarshal(data, &status)
	if err != nil {
		return qcloudStatus{}, err
	}
	if status.Code == 0 {
		return status, nil
	}
	return status, errors.New(status.Message)
}

func getIp() (string, error) {
	resp, err := http.Get("http://ip.cip.cc")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ip := string(body)
	ip = strings.Replace(ip, "\r\n", "", -1)
	ip = strings.Replace(ip, "\n", "", -1)
	if ip == "" {
		return "", errors.New("ip is empty")
	}
	return ip, nil
}

func ipIsChanged(ip string, tempFile *os.File) (bool, error) {
	data, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		return false, err
	}
	oldIp := string(data)
	if oldIp == ip {
		return false, nil
	}
	tempFile.Write([]byte(ip))
	return true, nil
}

func main() {
	var configFilePath string
	var setting ddnsConfig
	var showDebugInfo bool

	flag.StringVar(&configFilePath, "c", "./config.yaml", "config file path")
	flag.BoolVar(&showDebugInfo, "d", false, "show debug information")
	flag.Parse()

	config, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		fmt.Print("read config file err: ", err)
	}

	yaml.Unmarshal(config, &setting)

	tempFile, err := ioutil.TempFile(os.TempDir(), "ddns-ip-*")

	if err != nil {
		fmt.Println("create temp file err: ", err)
		return
	}

	defer os.Remove(tempFile.Name())

	crontab := cron.New()
	task := func() {
		ip, err := getIp()
		if err != nil {
			fmt.Println("get ip err: ", err)
			return
		}
		if showDebugInfo {
			fmt.Println("now ip: ", ip)
		}
		ipChanged, err := ipIsChanged(ip, tempFile)
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
		fmt.Println("new ip: ", ip)
		result, err := setRecordIp(setting, ip)
		if err != nil {
			fmt.Println("set dns err: ", err)
			return
		}
		fmt.Println(result)
	}
	fmt.Println("start listen ip change")
	fmt.Println("cron run use: ", setting.Cron)
	crontab.AddFunc(setting.Cron, task)
	task()
	crontab.Start()

	select {}
}
