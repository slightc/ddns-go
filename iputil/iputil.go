package iputil

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

var tempFile *os.File

func Init() error {
	tFile, err := ioutil.TempFile(os.TempDir(), "ddns-ip-*")

	tempFile = tFile

	return err
}

func Deinit() {
	os.Remove(tempFile.Name())
}

func GetIp() (string, error) {
	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get("http://cip.cc")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ip := string(body)
	reg := regexp.MustCompile(`\d{1,3}\.\d{1,3}.\d{1,3}.\d{1,3}`)
	ip = reg.FindString(ip)
	if ip == "" {
		return "", errors.New("ip is empty")
	}
	return ip, nil
}

func IpIsChanged(ip string) (bool, error) {
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
