package aliyun

import (
	"errors"
	"fmt"

	"main/common"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

type Aliyun struct {
	Key common.AccessKey
}

func (ali Aliyun) newClient() (*alidns.Client, error) {
	return alidns.NewClientWithAccessKey("cn-shenzhen", ali.Key.Id, ali.Key.Secret)
}

func (ali Aliyun) getRecordList(info *common.RecordInfo) (*alidns.DescribeDomainRecordsResponse, error) {
	client, err := ali.newClient()

	if err != nil {
		return nil, err
	}

	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"

	request.DomainName = info.Domain
	request.RRKeyWord = info.Name

	response, err := client.DescribeDomainRecords(request)
	return response, err
}

func (ali Aliyun) getRecordId(info *common.RecordInfo) (string, error) {
	list, err := ali.getRecordList(info)
	if err != nil {
		return "", err
	}
	records := list.DomainRecords.Record
	for _, record := range records {
		if record.RR == info.Name {
			return record.RecordId, nil
		}
	}
	return "", errors.New("no found record id")
}

func (ali Aliyun) setRecordIp(info *common.RecordInfo, ip string) error {
	client, err := ali.newClient()

	if err != nil {
		return err
	}

	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.Type = "A"

	request.RecordId = info.Id
	request.RR = info.Name
	request.Value = ip

	_, err = client.UpdateDomainRecord(request)
	return err
}

func (ali Aliyun) SetRecordIp(info *common.RecordInfo, ip string) error {
	if info == nil {
		return errors.New("info is nil")
	}
	if info.Id == "" {
		id, err := ali.getRecordId(info)
		fmt.Println("record Id is: ", id)
		if err != nil {
			return err
		}
		info.Id = id
	}
	return ali.setRecordIp(info, ip)
}
