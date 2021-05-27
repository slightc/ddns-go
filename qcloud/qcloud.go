package qcloud

import (
	"errors"
	"fmt"
	"main/common"
	"strconv"

	qCommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	qProfile "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

type Qcloud struct {
	Key common.AccessKey
}

func (q Qcloud) getClient() (*dnspod.Client, error) {
	credential := qCommon.NewCredential(
		q.Key.Id,
		q.Key.Secret,
	)

	cpf := qProfile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "dnspod.tencentcloudapi.com"
	return dnspod.NewClient(credential, "", cpf)
}

func (q Qcloud) getRecordList(info *common.RecordInfo) (*dnspod.DescribeRecordListResponse, error) {
	client, err := q.getClient()
	if err != nil {
		return nil, err
	}

	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = qCommon.StringPtr(info.Domain)

	return client.DescribeRecordList(request)
}

func (q Qcloud) getRecordId(info *common.RecordInfo) (string, error) {
	list, err := q.getRecordList(info)
	if err != nil {
		return "", err
	}
	records := list.Response.RecordList
	for _, record := range records {
		if *record.Name == info.Name {
			recordId := strconv.FormatInt(int64(*record.RecordId), 10)
			return recordId, nil
		}
	}
	return "", errors.New("no found record id")
}

func (q Qcloud) setRecordIp(info *common.RecordInfo, ip string) error {
	client, err := q.getClient()
	if err != nil {
		return err
	}

	recordId, err := strconv.ParseInt(info.Id, 10, 64)
	if err != nil {
		return err
	}

	request := dnspod.NewModifyRecordRequest()
	request.RecordType = qCommon.StringPtr("A")
	request.RecordLine = qCommon.StringPtr("默认")

	request.Domain = qCommon.StringPtr(info.Domain)
	request.SubDomain = qCommon.StringPtr(info.Name)
	request.RecordId = qCommon.Uint64Ptr(uint64(recordId))
	request.Value = qCommon.StringPtr(ip)

	_, err = client.ModifyRecord(request)

	return err
}

func (q Qcloud) SetRecordIp(info *common.RecordInfo, ip string) error {
	if info == nil {
		return errors.New("info is nil")
	}
	if info.Id == "" {
		id, err := q.getRecordId(info)
		fmt.Println("record Id is: ", id)
		if err != nil {
			return err
		}
		info.Id = id
	}
	return q.setRecordIp(info, ip)
}
