package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

type qcloudRecordListItem struct {
	Id   int    `yaml:"id"`
	Name string `yaml:"name"`
}
type qcloudRecordData struct {
	Records []qcloudRecordListItem `json:"records"`
}

type qcloudStatus struct {
	Code     int    `json:"code"`
	CodeDesc string `json:"codeDesc"`
	Message  string `json:"message"`
	// Data     qcloudStatusData `json:"data"`
}

type qcloudRecordInfo struct {
	qcloudStatus
	Data qcloudRecordData `json:"data"`
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

func getRecordList(config ddnsConfig) (qcloudRecordInfo, error) {
	sKI := qcloud.SecretData{}
	sKI.SecretId = config.SecretId
	sKI.SecretKey = config.SecretKey

	data, err := qcloud.GetRecordList(sKI, config.Domain)

	if err != nil {
		return qcloudRecordInfo{}, err
	}
	// fmt.Println(string(data))
	var status qcloudRecordInfo
	err = json.Unmarshal(data, &status)
	if err != nil {
		return qcloudRecordInfo{}, err
	}
	if status.Code == 0 {
		return status, nil
	}
	return status, errors.New(status.Message)
}

func getRecordId(config ddnsConfig, name string) (string, error) {
	recordState, err := getRecordList(config)
	if err != nil {
		return "", err
	}
	recordList := recordState.Data.Records
	for _, record := range recordList {
		if record.Name == name {
			return fmt.Sprint(record.Id), nil
		}
	}
	return "", errors.New("not found")
}

func updateRecordId(config *ddnsConfig) {
	if config.Id == "" {
		recordId, err := getRecordId(*config, config.Record)
		if err != nil {
			fmt.Println("get record list err: ", err)
			return
		}
		config.Id = recordId
		fmt.Println("get setting Id :", config.Id)
		fmt.Println("please write to config file")
	}
}

func getIp() (string, error) {
	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get("http://ip.cip.cc")
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
	tempFile.Truncate(0)
	tempFile.Seek(0, io.SeekStart)
	tempFile.Write([]byte(ip))
	tempFile.Sync()
	return true, nil
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

	if showRecordList {
		recordList, err := getRecordList(setting)
		if err == nil {
			fmt.Println("allRecord ", recordList)
		}
		return
	}

	updateRecordId(&setting)
	if setting.Id == "" {
		fmt.Println("get setting id failed")
		fmt.Println("please check Domain and Record in your config file")
	}

	tempFile, err := ioutil.TempFile(os.TempDir(), "ddns-ip-*")

	if err != nil {
		fmt.Println("create temp file err: ", err)
		return
	}

	defer os.Remove(tempFile.Name())

	crontab := cron.New()
	task := func() {
		updateRecordId(&setting)
		if setting.Id == "" {
			fmt.Println("get setting id failed")
			return
		}
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
		fmt.Println("== ", time.Now(), " ==")
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

	fmt.Println("waiting")

	select {
	case s := <-exitSignal:
		fmt.Println("exit with :", s)
		break
	}
}
