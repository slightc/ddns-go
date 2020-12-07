package qcloud

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

var qcloudApiHost = "cns.api.qcloud.com/v2/index.php"

type SecretData struct {
	SecretId  string
	SecretKey string
}

type RecordData struct {
	RecordId   string
	SubDomain  string
	RecordType string
	RecordLine string
	Value      string
	Ttl        int
	Mx         int
}

func addKv(kvList *[]string, key string, value string) {
	*kvList = append(*kvList, key+"="+value)
}

func createComKv(secretId string, addData []string) []string {
	kv := []string{}
	if addData != nil {
		kv = addData
	}
	nowTime := time.Now().Unix()
	rand.Seed(nowTime)
	addKv(&kv, "Nonce", strconv.FormatInt(int64(rand.Intn(10000)), 10))
	addKv(&kv, "Timestamp", strconv.FormatInt(nowTime, 10))
	addKv(&kv, "SecretId", secretId)
	return kv
}

func createParam(kvList []string) string {
	sort.Strings(kvList)

	out := strings.Join(kvList, "&")

	return out
}

func cryptoSignParam(data string, key string) string {
	k := []byte(key)
	mac := hmac.New(sha1.New, k)
	mac.Write([]byte(data))
	b64Data := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return url.QueryEscape(b64Data)
}

func signParam(method string, kvList []string, key string) string {
	params := createParam(kvList)
	out := method + qcloudApiHost + "?" + params
	return cryptoSignParam(out, key)
}

///////////////////////////////////////////////////////////////////////////

func qcloudGet(sKI SecretData, addParams []string) ([]byte, error) {
	params := createComKv(sKI.SecretId, addParams)
	addKv(&params, "Signature", signParam("GET", params, sKI.SecretKey))

	resp, err := http.Get("https://" + qcloudApiHost + "?" + createParam(params))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func fromString(data string, defaultData string) string {
	if data == "" {
		return defaultData
	}
	return data
}

func fromInt(data int, defaultData int) string {
	if data == 0 {
		return strconv.FormatInt(int64(defaultData), 10)
	}
	return strconv.FormatInt(int64(data), 10)
}

///////////////////////////////////////////////////////////////////////////

func GetRecordList(sKI SecretData, domain string) ([]byte, error) {
	params := []string{}
	addKv(&params, "Action", "RecordList")
	addKv(&params, "domain", domain)

	return qcloudGet(sKI, params)
}

func ModifyRecord(sKI SecretData, domain string, record RecordData) ([]byte, error) {

	if record.RecordId == "" || record.SubDomain == "" || record.RecordType == "" || record.Value == "" {
		return nil, errors.New("no must param")
	}

	params := []string{}
	addKv(&params, "Action", "RecordModify")
	addKv(&params, "domain", domain)

	addKv(&params, "recordId", record.RecordId)
	addKv(&params, "subDomain", record.SubDomain)
	addKv(&params, "recordType", record.RecordType)
	addKv(&params, "recordLine", (fromString(record.RecordLine, "默认")))
	addKv(&params, "value", record.Value)

	addKv(&params, "ttl", fromInt(record.Ttl, 600))
	addKv(&params, "mx", fromInt(record.Mx, 0))

	return qcloudGet(sKI, params)
}

///////////////////////////////////////////////////////////////////////////
